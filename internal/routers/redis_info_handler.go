package routers

import (
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/retail-ai-inc/beanqui/internal/redisx"
	"github.com/retail-ai-inc/beanqui/internal/routers/results"
)

type RedisInfo struct {
	client *redis.Client
}

func NewRedisInfo(client *redis.Client) *RedisInfo {
	return &RedisInfo{client: client}
}
func (t *RedisInfo) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	result, cancel := results.Get()
	defer cancel()

	flusher, ok := w.(http.Flusher)
	if !ok {

	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for {
		d, err := redisx.Info(r.Context(), t.client)

		if err != nil {
			result.Code = "1001"
			result.Msg = err.Error()
		}

		if err == nil {
			result.Data = d
		}
		_ = result.EventMsg(w, "redis_info")
		flusher.Flush()

		time.Sleep(10 * time.Second)
	}
}
