package routers

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/retail-ai-inc/beanq"
	"github.com/retail-ai-inc/beanq/helper/json"
	"github.com/retail-ai-inc/beanqui/internal/redisx"
	"github.com/retail-ai-inc/beanqui/internal/routers/consts"
	"github.com/retail-ai-inc/beanqui/internal/routers/results"
	"github.com/spf13/cast"
)

type Log struct {
}

// del ,retry,archive,detail
func (t *Log) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	result, cancel := results.Get()
	defer cancel()

	// job detail
	if r.Method == http.MethodGet {

		id := r.FormValue("id")
		msgType := r.FormValue("msgType")

		if id == "" || msgType == "" {
			// error
			result.Code = consts.MissParameterCode
			result.Msg = consts.MissParameterMsg
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
	// retry
	if r.Method == http.MethodPost {

		id := r.PostFormValue("id")
		if id == "" {
			result.Code = consts.MissParameterCode
			result.Msg = consts.MissParameterMsg
			_ = result.Json(w, http.StatusInternalServerError)
			return
		}
		if err := retryHandler(r.Context(), id); err != nil {
			result.Code = consts.InternalServerErrorCode
			result.Msg = err.Error()
			_ = result.Json(w, http.StatusInternalServerError)
			return
		}
		return
	}
	// delete job
	if r.Method == http.MethodDelete {

		id := r.FormValue("id")
		msgType := r.FormValue("msgType")
		if id == "" || msgType == "" {
			// error
			return
		}

		if err := deleteHandler(r.Context(), id, msgType); err != nil {
			// error
			return
		}
		_ = result.Json(w, http.StatusOK)
		return
	}
	if r.Method == http.MethodPut {

	}
}

// log detail
func detailHandler(ctx context.Context, id, msgType string) (map[string]any, error) {

	client := redisx.Client()

	res, err := client.Get(ctx, strings.Join([]string{redisx.BqConfig.Redis.Prefix, "logs", msgType, id}, ":")).Result()
	if err != nil {
		return nil, err
	}
	m := make(map[string]any)
	if err := json.Unmarshal([]byte(res), &m); err != nil {
		return nil, err
	}
	return m, nil
}

func deleteHandler(ctx context.Context, id, msgType string) (err error) {

	client := redisx.Client()

	key := strings.Join([]string{redisx.BqConfig.Redis.Prefix, "logs", msgType, id}, ":")
	cmd := client.Del(ctx, key)

	if cmd.Err() != nil {
		return err
	}

	return nil
}

func retryHandler(ctx context.Context, id string) error {

	client := redisx.Client()

	nid := cast.ToInt64(id)

	cmd := client.ZRange(ctx, strings.Join([]string{redisx.BqConfig.Redis.Prefix, "logs", "success"}, ":"), nid, nid)
	if err := cmd.Err(); err != nil {
		return err
	}
	vals := cmd.Val()
	if len(vals) < 1 {
		return errors.New("record is empty")
	}
	valByte := []byte(vals[0])

	newJson := json.Json
	payload := newJson.Get(valByte, "Payload").ToString()
	executeTime := newJson.Get(valByte, "ExecuteTime").ToString()
	groupName := newJson.Get(valByte, "Group").ToString()
	queue := newJson.Get(valByte, "Queue").ToString()
	queues := strings.Split(queue, ":")
	if len(queues) < 4 {
		return errors.New("data error")
	}

	dup, err := time.ParseInLocation(time.RFC3339, executeTime, time.Local)
	if err != nil {
		return err
	}

	publish := beanq.NewPublisher(redisx.BqConfig)
	task := beanq.NewMessage([]byte(payload))
	if err := publish.PublishWithContext(ctx, task, beanq.ExecuteTime(dup), beanq.Channel(groupName), beanq.Topic(queues[2])); err != nil {
		return err
	}

	return nil
}
