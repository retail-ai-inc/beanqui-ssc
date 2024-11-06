package routers

import (
	"github.com/retail-ai-inc/beanqui/internal/redisx"
	"github.com/retail-ai-inc/beanqui/internal/routers/errorx"
	"github.com/retail-ai-inc/beanqui/internal/routers/response"
	"net/http"
)

type Schedule struct {
}

func NewSchedule() *Schedule {
	return &Schedule{}
}

func (t *Schedule) List(ctx *BeanContext) error {

	result, cancel := response.Get()
	defer cancel()

	bt, err := redisx.QueueInfo(ctx.Request.Context())

	if err != nil {
		result.Code = errorx.InternalServerErrorCode
		result.Msg = err.Error()

		return result.Json(ctx.Writer, http.StatusInternalServerError)
	}
	result.Data = bt
	return result.Json(ctx.Writer, http.StatusOK)
}
