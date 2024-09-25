package routers

import (
	"context"
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

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	result, cancelR := results.Get()

	defer func() {
		cancel()
		cancelR()
	}()
	query := r.URL.Query()
	page := query.Get("page")
	pageSize := query.Get("pageSize")
	id := query.Get("id")
	status := query.Get("status")

	filter := bson.M{}
	if id != "" {
		filter["id"] = id
	}
	if status != "" {
		filter["status"] = status
	}

	mog := mongox.NewMongo()
	data, total, err := mog.EventLogs(ctx, filter, cast.ToInt64(page), cast.ToInt64(pageSize))
	if err != nil {
		return
	}
	datas := make(map[string]any, 0)
	datas["data"] = data
	datas["total"] = total
	datas["cursor"] = cast.ToInt64(page)
	result.Data = datas

	_ = result.Json(w, http.StatusOK)
	return

}
