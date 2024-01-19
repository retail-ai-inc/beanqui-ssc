package routers

import (
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/retail-ai-inc/beanqui/internal/redisx"
	"github.com/retail-ai-inc/beanqui/internal/routers/results"
)

type Client struct {
	client *redis.Client
}

func NewClient(client *redis.Client) *Client {
	return &Client{client: client}
}

func (t *Client) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	result, cancel := results.Get()
	defer cancel()

	data, err := redisx.ClientList(r.Context(), t.client)
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
