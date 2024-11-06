package routers

import (
	"github.com/retail-ai-inc/beanqui/internal/redisx"
	"github.com/retail-ai-inc/beanqui/internal/routers/errorx"
	"github.com/retail-ai-inc/beanqui/internal/routers/response"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"net/http"
	"strings"
)

type Dlq struct {
}

func NewDlq() *Dlq {
	return &Dlq{}
}

func (t *Dlq) List(ctx *BeanContext) error {

	w := ctx.Writer
	r := ctx.Request

	res, cancel := response.Get()
	defer cancel()

	client := redisx.Client()
	stream := strings.Join([]string{viper.GetString("redis.prefix"), "beanq-logic-log"}, ":")
	msgs, err := client.XRevRange(r.Context(), stream, "+", "-").Result()
	if err != nil {
		res.Code = errorx.InternalServerErrorMsg
		res.Msg = err.Error()
		return res.Json(w, http.StatusOK)

	}
	data := make([]map[string]any, 0)
	for _, msg := range msgs {
		val := msg.Values
		if v, ok := val["pendingRetry"]; ok {
			if cast.ToInt(v) > 0 {
				data = append(data, val)
			}
		}
	}
	res.Data = data
	return res.Json(w, http.StatusOK)

}
