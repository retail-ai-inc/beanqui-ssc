package routers

import (
	"net/http"
	"runtime"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/retail-ai-inc/beanqui/internal/redisx"
	"github.com/retail-ai-inc/beanqui/internal/routers/consts"
	"github.com/retail-ai-inc/beanqui/internal/routers/results"
	"github.com/spf13/viper"
)

type Dashboard struct {
	client redis.UniversalClient
}

func NewDashboard(client redis.UniversalClient) *Dashboard {
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

	// Queue Past 10 Minutes
	prefix := viper.GetString("redis.prefix")
	failKey := strings.Join([]string{prefix, "logs", "fail"}, ":")
	failCount := t.client.ZCard(r.Context(), failKey).Val()

	successKey := strings.Join([]string{prefix, "logs", "success"}, ":")
	successCount := t.client.ZCard(r.Context(), successKey).Val()

	result.Data = map[string]any{
		"queue_total":   keysLen,
		"db_size":       db_size,
		"num_cpu":       numCpu,
		"fail_count":    failCount,
		"success_count": successCount,
	}
	_ = result.Json(w, http.StatusOK)
	return
}
