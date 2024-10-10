package routers

import (
	"net/http"

	"github.com/retail-ai-inc/beanqui/internal/redisx"
	"github.com/retail-ai-inc/beanqui/internal/routers/consts"
	"github.com/retail-ai-inc/beanqui/internal/routers/response"
)

type Schedule struct {
}

func NewSchedule() *Schedule {
	return &Schedule{}
}

func (t *Schedule) List(w http.ResponseWriter, r *http.Request) {
	result, cancel := response.Get()
	defer cancel()

	bt, err := redisx.QueueInfo(r.Context())

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
