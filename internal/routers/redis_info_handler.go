package routers

import (
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/retail-ai-inc/beanqui/internal/redisx"
	"github.com/retail-ai-inc/beanqui/internal/routers/results"
)

type RedisInfo struct {
	client redis.UniversalClient
}

func NewRedisInfo(client redis.UniversalClient) *RedisInfo {
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

	ctx := r.Context()
	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			d, err := redisx.Info(ctx, t.client)

			if err != nil {
				result.Code = "1001"
				result.Msg = err.Error()
			}

			if err == nil {
				result.Data = d
			}
			_ = result.EventMsg(w, "redis_info")
			flusher.Flush()
			ticker.Reset(10 * time.Second)

		}
	}
}
