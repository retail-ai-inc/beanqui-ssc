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
func (t *Log) List(beanContext *BeanContext) error {

	result, cancel := response.Get()
	defer cancel()
	r := beanContext.Request
	w := beanContext.Writer

	id := r.FormValue("id")
	msgType := r.FormValue("msgType")

	if id == "" || msgType == "" {
		// error
		result.Code = errorx.MissParameterCode
		result.Msg = errorx.MissParameterMsg
		return result.Json(w, http.StatusBadRequest)
	}
	data, err := detailHandler(r.Context(), id, msgType)
	if err != nil {
		result.Code = "1003"
		result.Msg = err.Error()
		return result.Json(w, http.StatusInternalServerError)
	}
	result.Data = data
	return result.Json(w, http.StatusOK)
}

func (t *Log) Retry(beanContext *BeanContext) error {

	result, cancel := response.Get()
	defer cancel()

	r := beanContext.Request
	w := beanContext.Writer

	id := r.PostFormValue("id")
	msgType := r.PostFormValue("msgType")
	if msgType == "" {
		msgType = "success"
	}
	if id == "" {
		result.Code = errorx.MissParameterCode
		result.Msg = errorx.MissParameterMsg
		return result.Json(w, http.StatusInternalServerError)
	}
	if err := retryHandler(r.Context(), id, msgType); err != nil {
		result.Code = errorx.InternalServerErrorCode
		result.Msg = err.Error()
		return result.Json(w, http.StatusInternalServerError)
	}
	return result.Json(w, http.StatusOK)
}

func (t *Log) Delete(beanContext *BeanContext) error {
	result, cancel := response.Get()
	defer cancel()

	w := beanContext.Writer
	r := beanContext.Request

	msgType := r.FormValue("msgType")
	score := r.FormValue("score")
	key := strings.Join([]string{redisx.BqConfig.Redis.Prefix, "logs", msgType}, ":")

	if err := redisx.ZRemRangeByScore(r.Context(), key, score, score); err != nil {
		result.Code = errorx.InternalServerErrorCode
		result.Msg = err.Error()
		return result.Json(w, http.StatusInternalServerError)
	}
	return result.Json(w, http.StatusOK)
}
func (t *Log) Add(w http.ResponseWriter, r *http.Request) {
}

// log detail
func detailHandler(ctx context.Context, id, msgType string) (map[string]any, error) {

	key := strings.Join([]string{redisx.BqConfig.Redis.Prefix, "logs", msgType}, ":")

	var build strings.Builder
	build.Grow(3)
	build.WriteString("*")
	build.WriteString(id)
	build.WriteString("*")

	vals, _, err := redisx.ZScan(ctx, key, 0, build.String(), 1)
	if err != nil {
		return nil, err
	}
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

	keys, _, err := redisx.ZScan(ctx, key, 0, build.String(), 1)
	if err != nil {
		return err
	}
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
