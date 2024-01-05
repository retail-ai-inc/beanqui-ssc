package routers

import (
	"net/http"

	"github.com/retail-ai-inc/beanqui/internal/redisx"
	"github.com/retail-ai-inc/beanqui/internal/routers/results"
)

type Client struct {
}

func (t *Client) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	result, cancel := results.Get()
	defer cancel()

	client := redisx.Client()
	data, err := redisx.ClientList(r.Context(), client)
	if err != nil {
		result.Code = "1001"
		result.Msg = err.Error()
		_ = result.Json(w, http.StatusInternalServerError)
		return
	}
	result.Data = data
	_ = result.Json(w, http.StatusOK)
	return

}
