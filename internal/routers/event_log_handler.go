package routers

import (
	"github.com/retail-ai-inc/beanqui/internal/mongox"
	"github.com/retail-ai-inc/beanqui/internal/routers/results"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"time"
)

type EventLog struct {
}

func NewEventLog() *EventLog {
	return &EventLog{}
}

func (t *EventLog) List(w http.ResponseWriter, r *http.Request) {

	//ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	result, cancelR := results.Get()

	defer func() {
		//cancel()
		cancelR()
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
	ticker := time.NewTicker(time.Second)
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
