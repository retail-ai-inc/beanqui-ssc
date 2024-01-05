package routers

import (
	"net/http"

	"github.com/retail-ai-inc/beanqui/internal/redisx"
	"github.com/retail-ai-inc/beanqui/internal/routers/consts"
	"github.com/retail-ai-inc/beanqui/internal/routers/results"
)

type Schedule struct {
}

func (t *Schedule) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	result, cancel := results.Get()
	defer cancel()

	bt, err := redisx.QueueInfo(r.Context(), redisx.Client(), redisx.ScheduleQueueKey(redisx.BqConfig.Redis.Prefix))

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
