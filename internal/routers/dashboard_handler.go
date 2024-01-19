package routers

import (
	"net/http"
	"runtime"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/retail-ai-inc/beanqui/internal/redisx"
	"github.com/retail-ai-inc/beanqui/internal/routers/consts"
	"github.com/retail-ai-inc/beanqui/internal/routers/results"
)

type Dashboard struct {
	client *redis.Client
}

func NewDashboard(client *redis.Client) *Dashboard {
	return &Dashboard{client: client}
}

func (t *Dashboard) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	result, cancel := results.Get()
	defer cancel()

	numCpu := runtime.NumCPU()

	// get queue total
	keys, err := redisx.Keys(r.Context(), t.client, strings.Join([]string{redisx.BqConfig.Prefix, "*", "stream"}, ":"))
	if err != nil {
		result.Code = consts.InternalServerErrorCode
		result.Msg = err.Error()
		_ = result.Json(w, http.StatusInternalServerError)
		return
	}
	keysLen := len(keys)

	// db size
	db_size, err := t.client.DBSize(r.Context()).Result()
	if err != nil {

		result.Code = consts.InternalServerErrorCode
		result.Msg = err.Error()
		_ = result.Json(w, http.StatusInternalServerError)

		return
	}

	result.Data = map[string]any{
		"queue_total": keysLen,
		"db_size":     db_size,
		"num_cpu":     numCpu,
	}
	_ = result.Json(w, http.StatusOK)
	return
}
