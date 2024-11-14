package routers

import (
	"encoding/json"
	"github.com/retail-ai-inc/beanq/v3"
	"github.com/retail-ai-inc/beanqui/internal/mongox"
	"github.com/retail-ai-inc/beanqui/internal/redisx"
	"github.com/retail-ai-inc/beanqui/internal/routers/errorx"
	"github.com/retail-ai-inc/beanqui/internal/routers/response"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson"
	"io"
	"net/http"
	"time"
)

type EventLog struct {
	Id string `json:"id"`
}

func NewEventLog() *EventLog {
	return &EventLog{}
}

func (t *EventLog) List(ctx *BeanContext) error {

	r := ctx.Request
	w := ctx.Writer

	result, cancel := response.Get()
	defer func() {
		cancel()
	}()
	query := r.URL.Query()
	page := cast.ToInt64(query.Get("page"))
	pageSize := cast.ToInt64(query.Get("pageSize"))
	id := query.Get("id")
	status := query.Get("status")

	filter := bson.M{}
	if id != "" {
		filter["id"] = id
	}
	if status != "" {
		filter["status"] = status
	}
	if page <= 0 {
		page = 0
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	flush, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "server err", http.StatusInternalServerError)
		return nil
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	mog := mongox.NewMongo()
	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()

	datas := make(map[string]any, 3)
	nctx := r.Context()

	for {
		select {
		case <-nctx.Done():
			return nctx.Err()
		case <-ticker.C:

			data, total, err := mog.EventLogs(nctx, filter, page, pageSize)
			if err != nil {
				result.Code = "1001"
				result.Msg = err.Error()
			}
			if err == nil {
				datas["data"] = data
				datas["total"] = total
				datas["cursor"] = page
				result.Data = datas
			}

			_ = result.EventMsg(w, "event_log")
			flush.Flush()
			ticker.Reset(5 * time.Second)
		}
	}
}

func (t *EventLog) Detail(ctx *BeanContext) error {

	res, cancel := response.Get()
	defer cancel()

	r := ctx.Request
	w := ctx.Writer

	id := r.URL.Query().Get("id")
	mog := mongox.NewMongo()
	data, err := mog.DetailEventLog(r.Context(), id)
	if err != nil {
		res.Msg = err.Error()
		res.Code = errorx.InternalServerErrorCode
		return res.Json(w, http.StatusInternalServerError)

	}
	res.Data = data
	return res.Json(w, http.StatusOK)

}

func (t *EventLog) Delete(ctx *BeanContext) error {

	res, cancel := response.Get()
	defer cancel()

	w := ctx.Writer
	r := ctx.Request

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		res.Code = errorx.TypeErrorCode
		res.Msg = err.Error()
		return res.Json(ctx.Writer, http.StatusOK)
	}
	if err := json.Unmarshal(body, &t); err != nil {
		res.Code = errorx.TypeErrorCode
		res.Msg = err.Error()
		return res.Json(ctx.Writer, http.StatusOK)
	}

	id := t.Id
	mog := mongox.NewMongo()
	count, err := mog.Delete(r.Context(), id)
	if err != nil {
		res.Msg = err.Error()
		res.Code = errorx.InternalServerErrorCode
		return res.Json(w, http.StatusInternalServerError)
	}
	res.Data = count
	return res.Json(w, http.StatusOK)
}

type editInfo struct {
	Payload any    `json:"payload"`
	Id      string `json:"id"`
}

func (t *EventLog) Edit(ctx *BeanContext) error {
	res, cancel := response.Get()
	defer cancel()

	r := ctx.Request
	w := ctx.Writer

	var info editInfo
	b, err := io.ReadAll(r.Body)
	if err != nil {
		res.Msg = err.Error()
		res.Code = errorx.InternalServerErrorCode
		return res.Json(w, http.StatusInternalServerError)

	}
	if err := json.Unmarshal(b, &info); err != nil {
		res.Msg = err.Error()
		res.Code = errorx.InternalServerErrorCode
		return res.Json(w, http.StatusInternalServerError)

	}

	mog := mongox.NewMongo()
	count, err := mog.Edit(r.Context(), info.Id, info.Payload)
	if err != nil {
		res.Msg = err.Error()
		res.Code = errorx.InternalServerErrorCode
		return res.Json(w, http.StatusInternalServerError)

	}
	res.Data = count
	return res.Json(w, http.StatusOK)

}

func (t *EventLog) Retry(ctx *BeanContext) error {
	res, cancel := response.Get()
	defer cancel()

	w := ctx.Writer
	r := ctx.Request

	m := make(map[string]any)
	id := r.FormValue("id")
	m["uniqueId"] = id
	nctx := r.Context()

	data := make(map[string]any)
	if err := json.Unmarshal([]byte(r.FormValue("data")), &data); err != nil {
		res.Msg = err.Error()
		res.Code = errorx.InternalServerErrorCode
		return res.Json(w, http.StatusInternalServerError)

	}

	moodType := ""
	if v, ok := data["moodType"]; ok {
		moodType = v.(string)
	}
	payload := ""
	if v, ok := data["payload"]; ok {
		payload = v.(string)
	}
	channel := ""
	if v, ok := data["channel"]; ok {
		channel = v.(string)
	}
	topic := ""
	if v, ok := data["topic"]; ok {
		topic = v.(string)
	}

	bq := beanq.New(&redisx.BqConfig)

	if moodType == string(beanq.SEQUENTIAL) {
		return res.Json(w, http.StatusOK)
	}
	if moodType == string(beanq.DELAY) {
		executeTime := ""
		if v, ok := data["executeTime"]; ok {
			executeTime = v.(string)
		}
		dup, err := time.ParseInLocation(time.RFC3339, executeTime, time.Local)
		if err != nil {
			res.Msg = err.Error()
			res.Code = errorx.InternalServerErrorCode
			return res.Json(w, http.StatusOK)

		}
		if err := bq.BQ().WithContext(nctx).PublishAtTime(channel, topic, []byte(payload), dup); err != nil {
			res.Msg = err.Error()
			res.Code = errorx.InternalServerErrorCode
			return res.Json(w, http.StatusOK)

		}
		return res.Json(w, http.StatusOK)
	}
	if err := bq.BQ().WithContext(nctx).Publish(channel, topic, []byte(payload)); err != nil {
		res.Msg = err.Error()
		res.Code = errorx.InternalServerErrorCode
		return res.Json(w, http.StatusOK)

	}

	return res.Json(w, http.StatusOK)
}
