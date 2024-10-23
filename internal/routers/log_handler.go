package routers

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/retail-ai-inc/beanq/v3"
	"github.com/retail-ai-inc/beanq/v3/helper/json"
	"github.com/retail-ai-inc/beanqui/internal/redisx"
	"github.com/retail-ai-inc/beanqui/internal/routers/errorx"
	"github.com/retail-ai-inc/beanqui/internal/routers/response"
)

type Log struct {
}

func NewLog() *Log {
	return &Log{}
}

// del ,retry,archive,detail
func (t *Log) List(w http.ResponseWriter, r *http.Request) {

	result, cancel := response.Get()
	defer cancel()

	id := r.FormValue("id")
	msgType := r.FormValue("msgType")

	if id == "" || msgType == "" {
		// error
		result.Code = errorx.MissParameterCode
		result.Msg = errorx.MissParameterMsg
		_ = result.Json(w, http.StatusBadRequest)
		return
	}
	data, err := detailHandler(r.Context(), id, msgType)
	if err != nil {
		result.Code = "1003"
		result.Msg = err.Error()
		_ = result.Json(w, http.StatusInternalServerError)
		return
	}
	result.Data = data
	_ = result.Json(w, http.StatusOK)
	return

}

func (t *Log) Retry(w http.ResponseWriter, r *http.Request) {

	result, cancel := response.Get()
	defer cancel()

	id := r.PostFormValue("id")
	msgType := r.PostFormValue("msgType")
	if msgType == "" {
		msgType = "success"
	}
	if id == "" {
		result.Code = errorx.MissParameterCode
		result.Msg = errorx.MissParameterMsg
		_ = result.Json(w, http.StatusInternalServerError)
		return
	}
	if err := retryHandler(r.Context(), id, msgType); err != nil {
		result.Code = errorx.InternalServerErrorCode
		result.Msg = err.Error()
		_ = result.Json(w, http.StatusInternalServerError)
		return
	}
	return
}

func (t *Log) Delete(w http.ResponseWriter, r *http.Request) {
	result, cancel := response.Get()
	defer cancel()

	msgType := r.FormValue("msgType")
	score := r.FormValue("score")
	client := redisx.Client()
	key := strings.Join([]string{redisx.BqConfig.Redis.Prefix, "logs", msgType}, ":")
	cmd := client.ZRemRangeByScore(r.Context(), key, score, score)
	if cmd.Err() != nil {
		result.Code = errorx.InternalServerErrorCode
		result.Msg = cmd.Err().Error()
		_ = result.Json(w, http.StatusInternalServerError)
		return
	}
	_ = result.Json(w, http.StatusOK)
	return
}
func (t *Log) Add(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut {

	}
}

// log detail
func detailHandler(ctx context.Context, id, msgType string) (map[string]any, error) {

	key := strings.Join([]string{redisx.BqConfig.Redis.Prefix, "logs", msgType}, ":")

	var build strings.Builder
	build.Grow(3)
	build.WriteString("*")
	build.WriteString(id)
	build.WriteString("*")
	client := redisx.Client()
	vals, _ := client.ZScan(ctx, key, 0, build.String(), 1).Val()
	if len(vals) <= 0 {
		return nil, errors.New("record is empty")
	}

	m := make(map[string]any)
	if err := json.Unmarshal([]byte(vals[0]), &m); err != nil {
		return nil, err
	}
	return m, nil
}

func retryHandler(ctx context.Context, id, msgType string) error {

	key := strings.Join([]string{redisx.BqConfig.Redis.Prefix, "logs", msgType}, ":")

	var build strings.Builder
	build.Grow(3)
	build.WriteString("*")
	build.WriteString(id)
	build.WriteString("*")
	client := redisx.Client()

	keys, _ := client.ZScan(ctx, key, 0, build.String(), 1).Val()
	if len(keys) <= 0 {
		return errors.New("record is empty")
	}

	valByte := []byte(keys[0])

	newJson := json.Json
	payload := newJson.Get(valByte, "Payload").ToString()
	executeTime := newJson.Get(valByte, "ExecuteTime").ToString()
	groupName := newJson.Get(valByte, "Channel").ToString()
	queue := newJson.Get(valByte, "Topic").ToString()

	dup, err := time.ParseInLocation(time.RFC3339, executeTime, time.Local)
	if err != nil {
		return err
	}

	bq := beanq.New(&redisx.BqConfig)
	if err := bq.BQ().PublishAtTime(groupName, queue, []byte(payload), dup); err != nil {
		return err
	}

	return nil
}
