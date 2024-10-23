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
}

func NewEventLog() *EventLog {
	return &EventLog{}
}

func (t *EventLog) List(w http.ResponseWriter, r *http.Request) {

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
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	mog := mongox.NewMongo()
	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()

	datas := make(map[string]any, 3)
	ctx := r.Context()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:

			data, total, err := mog.EventLogs(ctx, filter, page, pageSize)
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

func (t *EventLog) Detail(w http.ResponseWriter, r *http.Request) {
	res, cancel := response.Get()
	defer cancel()

	id := r.URL.Query().Get("id")
	mog := mongox.NewMongo()
	data, err := mog.DetailEventLog(r.Context(), id)
	if err != nil {
		res.Msg = err.Error()
		res.Code = errorx.InternalServerErrorCode
		_ = res.Json(w, http.StatusInternalServerError)
		return
	}
	res.Data = data
	_ = res.Json(w, http.StatusOK)
	return
}

func (t *EventLog) Delete(w http.ResponseWriter, r *http.Request) {
	res, cancel := response.Get()
	defer cancel()

	id := r.URL.Query().Get("id")
	mog := mongox.NewMongo()
	count, err := mog.Delete(r.Context(), id)
	if err != nil {
		res.Msg = err.Error()
		res.Code = errorx.InternalServerErrorCode
		_ = res.Json(w, http.StatusInternalServerError)
		return
	}
	res.Data = count
	_ = res.Json(w, http.StatusOK)
	return
}

type editInfo struct {
	Id      string `json:"id"`
	Payload any    `json:"payload"`
}

func (t *EventLog) Edit(w http.ResponseWriter, r *http.Request) {
	res, cancel := response.Get()
	defer cancel()

	var info editInfo
	b, err := io.ReadAll(r.Body)
	if err != nil {
		res.Msg = err.Error()
		res.Code = errorx.InternalServerErrorCode
		_ = res.Json(w, http.StatusInternalServerError)
		return
	}
	if err := json.Unmarshal(b, &info); err != nil {
		res.Msg = err.Error()
		res.Code = errorx.InternalServerErrorCode
		_ = res.Json(w, http.StatusInternalServerError)
		return
	}

	mog := mongox.NewMongo()
	count, err := mog.Edit(r.Context(), info.Id, info.Payload)
	if err != nil {
		res.Msg = err.Error()
		res.Code = errorx.InternalServerErrorCode
		_ = res.Json(w, http.StatusInternalServerError)
		return
	}
	res.Data = count
	_ = res.Json(w, http.StatusOK)
	return
}

func (t *EventLog) Retry(w http.ResponseWriter, r *http.Request) {
	res, cancel := response.Get()
	defer cancel()
	m := make(map[string]any)
	id := r.FormValue("id")
	m["uniqueId"] = id
	ctx := r.Context()

	data := make(map[string]any)
	if err := json.Unmarshal([]byte(r.FormValue("data")), &data); err != nil {
		res.Msg = err.Error()
		res.Code = errorx.InternalServerErrorCode
		_ = res.Json(w, http.StatusInternalServerError)
		return
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
		return
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
			_ = res.Json(w, http.StatusOK)
			return
		}
		if err := bq.BQ().WithContext(ctx).PublishAtTime(channel, topic, []byte(payload), dup); err != nil {
			res.Msg = err.Error()
			res.Code = errorx.InternalServerErrorCode
			_ = res.Json(w, http.StatusOK)
			return
		}
		return
	}
	if err := bq.BQ().WithContext(ctx).Publish(channel, topic, []byte(payload)); err != nil {
		res.Msg = err.Error()
		res.Code = errorx.InternalServerErrorCode
		_ = res.Json(w, http.StatusOK)
		return
	}

	return
}
