package routers

import (
	"context"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/retail-ai-inc/beanq"
	"github.com/retail-ai-inc/beanq/helper/json"
	"github.com/retail-ai-inc/beanqui/internal/jwtx"
	"github.com/retail-ai-inc/beanqui/internal/redisx"
	"github.com/retail-ai-inc/beanqui/internal/routers/consts"
	"github.com/retail-ai-inc/beanqui/internal/routers/results"
	"github.com/retail-ai-inc/beanqui/internal/simple_router"

	"github.com/spf13/cast"
)

func IndexHandler(ctx *simple_router.Context) error {

	url := ctx.Request().RequestURI
	if strings.HasSuffix(url, ".vue") {
		ctx.Response().Header().Set("Content-Type", "application/octet-stream")
	}
	var dir string = "./"
	_, f, _, ok := runtime.Caller(0)
	if ok {
		dir = filepath.Dir(f)
	}

	hdl := http.FileServer(http.Dir(path.Join(dir, "../../ui/")))
	hdl.ServeHTTP(ctx.Response(), ctx.Request())
	return nil
}
func DashboardHandler(ctx *simple_router.Context) error {

	result, cancel := results.Get()
	defer cancel()

	client := redisx.Client()
	// get queue total
	keys, err := redisx.Keys(ctx.Context(), client, strings.Join([]string{redisx.BqConfig.Prefix, "*", "stream"}, ":"))
	if err != nil {
		result.Code = consts.InternalServerErrorCode
		result.Msg = err.Error()
		return ctx.Json(http.StatusInternalServerError, result)
	}
	keysLen := len(keys)

	// db size
	db_size, err := client.DBSize(ctx.Context()).Result()
	if err != nil {
		result.Msg = err.Error()
		return ctx.Json(http.StatusInternalServerError, result)
	}

	result.Data = map[string]any{
		"queue_total": keysLen,
		"db_size":     db_size,
	}
	return ctx.Json(http.StatusOK, result)
}
func LoginHandler(ctx *simple_router.Context) error {

	request := ctx.Request()
	username := request.PostFormValue("username")
	password := request.PostFormValue("password")

	result, cancel := results.Get()
	defer cancel()

	if username != "aa" && password != "bb" {
		result.Code = consts.InternalServerErrorCode
		result.Msg = "username or password mismatch"
		return ctx.Json(http.StatusUnauthorized, result)
	}
	claim := jwtx.Claim{
		UserName: username,
		Claims: jwt.RegisteredClaims{
			Issuer:    "Trial China",
			Subject:   "beanq monitor ui",
			Audience:  nil,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7200 * time.Second)),
			NotBefore: nil,
			IssuedAt:  nil,
			ID:        "",
		},
	}

	token, err := jwtx.MakeRsaToken(claim)
	if err != nil {
		result.Code = consts.InternalServerErrorCode
		result.Msg = err.Error()
		return ctx.Json(http.StatusInternalServerError, result)
	}

	result.Data = map[string]any{"token": token}
	return ctx.Json(http.StatusOK, result)

}

func ScheduleHandler(ctx *simple_router.Context) error {

	result, cancel := results.Get()
	defer cancel()

	bt, err := queueInfo(ctx.Context(), redisx.ScheduleQueueKey(redisx.BqConfig.Redis.Prefix))

	if err != nil {
		result.Code = consts.InternalServerErrorCode
		result.Msg = err.Error()
		return ctx.Json(http.StatusInternalServerError, result)
	}
	result.Data = bt
	return ctx.Json(http.StatusOK, result)
}

func QueueHandler(ctx *simple_router.Context) error {

	result, cancel := results.Get()
	defer cancel()
	nctx := ctx.Context()

	bt, err := queueInfo(nctx, redisx.QueueKey(redisx.BqConfig.Redis.Prefix))
	if err != nil {
		result.Code = consts.InternalServerErrorCode
		result.Msg = err.Error()
		return ctx.Json(http.StatusInternalServerError, result)
	}

	result.Data = bt

	return ctx.Json(http.StatusOK, result)
}

func LogArchiveHandler(ctx *simple_router.Context) error {

	result, cancel := results.Get()
	defer cancel()

	return ctx.Json(http.StatusOK, result)
}

func LogRetryHandler(ctx *simple_router.Context) error {

	result, cancel := results.Get()
	defer cancel()

	req := ctx.Request()
	id := req.PostFormValue("id")
	if id == "" {
		result.Code = consts.MissParameterCode
		result.Msg = consts.MissParameterMsg
		return ctx.Json(http.StatusInternalServerError, result)
	}
	client := redisx.Client()

	nid := cast.ToInt64(id)

	cmd := client.ZRange(ctx.Context(), strings.Join([]string{redisx.BqConfig.Redis.Prefix, "logs", "success"}, ":"), nid, nid)
	if err := cmd.Err(); err != nil {
		result.Code = consts.InternalServerErrorCode
		result.Msg = err.Error()
		return ctx.Json(http.StatusInternalServerError, result)
	}
	vals := cmd.Val()
	if len(vals) < 1 {
		result.Code = consts.InternalServerErrorCode
		result.Msg = "record is empty"
		return ctx.Json(http.StatusInternalServerError, result)
	}
	valByte := []byte(vals[0])

	newJson := json.Json
	payload := newJson.Get(valByte, "Payload").ToString()
	executeTime := newJson.Get(valByte, "ExecuteTime").ToString()
	groupName := newJson.Get(valByte, "Group").ToString()
	queue := newJson.Get(valByte, "Queue").ToString()
	queues := strings.Split(queue, ":")
	if len(queues) < 4 {
		result.Code = consts.InternalServerErrorCode
		result.Msg = "data error"
		return ctx.Json(http.StatusInternalServerError, result)
	}

	dup, err := time.ParseInLocation(time.RFC3339, executeTime, time.Local)
	if err != nil {
		result.Code = consts.InternalServerErrorCode
		result.Msg = err.Error()
		return ctx.Json(http.StatusInternalServerError, result)
	}

	publish := beanq.NewPublisher(redisx.BqConfig)
	task := beanq.NewTask([]byte(payload))
	if err := publish.PublishWithContext(ctx.Context(), task, beanq.ExecuteTime(dup), beanq.Group(groupName), beanq.Queue(queues[2])); err != nil {
		result.Code = consts.InternalServerErrorCode
		result.Msg = err.Error()
		return ctx.Json(http.StatusInternalServerError, result)
	}

	return ctx.Json(http.StatusOK, result)
}

func LogDelHandler(ctx *simple_router.Context) error {

	result, cancel := results.Get()
	defer cancel()
	req := ctx.Request()
	id := req.FormValue("id")
	if id == "" {
		result.Code = consts.MissParameterCode
		result.Msg = consts.MissParameterMsg
		return ctx.Json(http.StatusInternalServerError, result)
	}

	client := redisx.Client()

	nid := cast.ToInt64(id)
	var start int64
	start = nid - 1
	if start <= 0 {
		start = 0
	}

	cmd := client.ZRemRangeByRank(ctx.Context(), strings.Join([]string{redisx.BqConfig.Redis.Prefix, "logs", "success"}, ":"), start, nid)

	if cmd.Err() != nil {
		result.Code = consts.InternalServerErrorCode
		result.Msg = cmd.Err().Error()
		return ctx.Json(http.StatusInternalServerError, result)
	}

	return ctx.Json(http.StatusOK, result)
}

func LogHandler(ctx *simple_router.Context) error {

	resultRes, cancel := results.Get()
	defer cancel()

	client := redisx.Client()

	var (
		page, pageSize int64
		dataType       string = "success"
		matchStr       string = strings.Join([]string{redisx.BqConfig.Redis.Prefix, "logs", "success"}, ":")
	)

	req := ctx.Request()
	page = cast.ToInt64(req.FormValue("page"))
	pageSize = cast.ToInt64(req.FormValue("pageSize"))
	dataType = req.FormValue("type")

	if dataType != "success" && dataType != "error" {
		resultRes.Code = consts.TypeErrorCode
		resultRes.Msg = consts.TypeErrorMsg

		return ctx.Json(http.StatusInternalServerError, resultRes)

	}

	nowPage := (page - 1) * pageSize
	if nowPage <= 0 {
		nowPage = 0
	}
	nowPageSize := page * pageSize
	if nowPageSize <= 0 {
		nowPageSize = 9
	}

	if dataType == "error" {
		matchStr = strings.Join([]string{redisx.BqConfig.Redis.Prefix, "logs", "error"}, ":")
	}

	data, err := redisx.ZRange(ctx.Context(), client, matchStr, nowPage, nowPageSize)
	if err != nil {
		resultRes.Code = "1001"
		resultRes.Msg = err.Error()

		return ctx.Json(http.StatusInternalServerError, resultRes)
	}
	resultRes.Data = data

	return ctx.Json(http.StatusOK, resultRes)
}

func RedisHandler(ctx *simple_router.Context) error {

	result, cancel := results.Get()
	defer cancel()

	client := redisx.Client()

	d, err := redisx.Info(ctx.Context(), client)
	if err != nil {
		result.Code = "1001"
		result.Msg = err.Error()
		return ctx.Json(http.StatusInternalServerError, result)
	}

	result.Data = d

	return ctx.Json(http.StatusOK, result)
}
func ClientListHandler(ctx *simple_router.Context) error {

	result, cancel := results.Get()
	defer cancel()

	client := redisx.Client()
	data, err := redisx.ClientList(ctx.Context(), client)
	if err != nil {
		result.Code = "1001"
		result.Msg = err.Error()
		return ctx.Json(http.StatusInternalServerError, result)
	}
	result.Data = data
	return ctx.Json(http.StatusOK, result)
}
func queueInfo(ctx context.Context, queueKey string) (any, error) {

	client := redisx.Client()

	// get queues
	queues, err := redisx.Keys(ctx, client, queueKey)
	if err != nil {
		return nil, err
	}

	d := make([]map[string]any, 0, len(queues))
	for _, queue := range queues {

		queueArr := strings.Split(queue, ":")
		if len(queueArr) < 4 {
			continue
		}
		objStr := redisx.Object(ctx, client, queue)
		// get memory
		r, err := client.MemoryUsage(ctx, queue).Result()
		if err != nil {
			log.Println(err)
			continue
		}
		d = append(d, map[string]any{"group": queueArr[1], "queue": queueArr[2], "state": "Run", "size": objStr.SerizlizedLength, "memory": r, "process": objStr.LruSecondsIdle, "fail": 0})
	}

	return d, nil
}
