package routers

import (
	"context"
	"github.com/retail-ai-inc/beanqui/internal/mongox"
	"github.com/retail-ai-inc/beanqui/internal/routers/results"
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

	mog := mongox.NewMongo()
	data, err := mog.EventLogs(ctx, nil, 0, 10)
	if err != nil {
		return
	}
	result.Data = data
	_ = result.Json(w, http.StatusOK)
	return

}
