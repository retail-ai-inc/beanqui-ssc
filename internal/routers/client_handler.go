package routers

import (
	"net/http"

	"github.com/retail-ai-inc/beanqui/internal/redisx"
	"github.com/retail-ai-inc/beanqui/internal/routers/response"
)

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

func (t *Client) List(ctx *BeanContext) error {

	r := ctx.Request
	w := ctx.Writer

	result, cancel := response.Get()
	defer cancel()

	data, err := redisx.ClientList(r.Context())
	if err != nil {
		result.Code = "1001"
		result.Msg = err.Error()
		return result.Json(w, http.StatusInternalServerError)

	}
	result.Data = data
	return result.Json(w, http.StatusOK)

}
