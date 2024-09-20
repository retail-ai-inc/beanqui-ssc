package routers

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/retail-ai-inc/beanq"
	"github.com/retail-ai-inc/beanq/helper/json"
	"github.com/retail-ai-inc/beanqui/internal/redisx"
	"github.com/retail-ai-inc/beanqui/internal/routers/consts"
	"github.com/retail-ai-inc/beanqui/internal/routers/results"
)

type Log struct {
	client redis.UniversalClient
}

func NewLog(client redis.UniversalClient) *Log {
	return &Log{client: client}
}

// del ,retry,archive,detail
func (t *Log) List(w http.ResponseWriter, r *http.Request) {

	result, cancel := results.Get()
	defer cancel()

	id := r.FormValue("id")
	msgType := r.FormValue("msgType")

	if id == "" || msgType == "" {
		// error
		result.Code = consts.MissParameterCode
		result.Msg = consts.MissParameterMsg
		_ = result.Json(w, http.StatusBadRequest)
		return
	}
	data, err := detailHandler(r.Context(), t.client, id, msgType)
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

	result, cancel := results.Get()
	defer cancel()

	id := r.PostFormValue("id")
	msgType := r.PostFormValue("msgType")
	if msgType == "" {
		msgType = "success"
	}
	if id == "" {
		result.Code = consts.MissParameterCode
		result.Msg = consts.MissParameterMsg
		_ = result.Json(w, http.StatusInternalServerError)
		return
	}
	if err := retryHandler(r.Context(), t.client, id, msgType); err != nil {
		result.Code = consts.InternalServerErrorCode
		result.Msg = err.Error()
		_ = result.Json(w, http.StatusInternalServerError)
		return
	}
	return
}

func (t *Log) Delete(w http.ResponseWriter, r *http.Request) {
	result, cancel := results.Get()
	defer cancel()

	msgType := r.FormValue("msgType")
	score := r.FormValue("score")

	key := strings.Join([]string{redisx.BqConfig.Redis.Prefix, "logs", msgType}, ":")
	cmd := t.client.ZRemRangeByScore(r.Context(), key, score, score)
	if cmd.Err() != nil {
		result.Code = consts.InternalServerErrorCode
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
func detailHandler(ctx context.Context, client redis.UniversalClient, id, msgType string) (map[string]any, error) {

	key := strings.Join([]string{redisx.BqConfig.Redis.Prefix, "logs", msgType}, ":")

	var build strings.Builder
	build.Grow(3)
	build.WriteString("*")
	build.WriteString(id)
	build.WriteString("*")

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

func retryHandler(ctx context.Context, client redis.UniversalClient, id, msgType string) error {

	key := strings.Join([]string{redisx.BqConfig.Redis.Prefix, "logs", msgType}, ":")

	var build strings.Builder
	build.Grow(3)
	build.WriteString("*")
	build.WriteString(id)
	build.WriteString("*")

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

	publish := beanq.NewPublisher(redisx.BqConfig)
	task := beanq.NewMessage([]byte(payload))
	if err := publish.PublishWithContext(ctx, task, beanq.ExecuteTime(dup), beanq.Channel(groupName), beanq.Topic(queue)); err != nil {
		return err
	}

	return nil
}
