package routers

import (
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/retail-ai-inc/beanqui/internal/redisx"
	"github.com/retail-ai-inc/beanqui/internal/routers/consts"
	"github.com/retail-ai-inc/beanqui/internal/routers/results"
)

type Schedule struct {
	client *redis.Client
}

func NewSchedule(client *redis.Client) *Schedule {
	return &Schedule{client: client}
}

func (t *Schedule) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	result, cancel := results.Get()
	defer cancel()

	bt, err := redisx.QueueInfo(r.Context(), t.client, redisx.ScheduleQueueKey(redisx.BqConfig.Redis.Prefix))

	if err != nil {
		result.Code = consts.InternalServerErrorCode
		result.Msg = err.Error()

		_ = result.Json(w, http.StatusInternalServerError)

		return
	}
	result.Data = bt
	_ = result.Json(w, http.StatusOK)

	return
}
