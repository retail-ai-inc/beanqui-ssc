package routers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/retail-ai-inc/beanqui/internal/redisx"
	"github.com/retail-ai-inc/beanqui/internal/routers/response"
)

type RedisInfo struct {
}

func NewRedisInfo() *RedisInfo {
	return &RedisInfo{}
}
func (t *RedisInfo) Info(ctx *BeanContext) error {

	result, cancel := response.Get()
	defer cancel()
	w := ctx.Writer
	r := ctx.Request

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "server error", http.StatusInternalServerError)
		return nil
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	nctx := r.Context()
	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-nctx.Done():
			return nctx.Err()
		case <-ticker.C:
			d, err := redisx.Info(nctx)

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

func (t *RedisInfo) Monitor(ctx *BeanContext) error {

	res, cancel := response.Get()
	defer cancel()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	w := ctx.Writer
	r := ctx.Request

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "server error", http.StatusInternalServerError)
		return nil
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for {
		select {
		case <-r.Context().Done():
			return nil
		case <-ticker.C:

			str := redisx.Monitor(r.Context())
			str = strings.ReplaceAll(str, "MONITOR:", "")
			if strings.Contains(str, "OK") {
				continue
			}
			res.Data = fmt.Sprintf("Time:%s,Command:%s", time.Now(), str)
			_ = res.EventMsg(w, "redis_monitor")
			flusher.Flush()
			ticker.Reset(time.Second)
		}
	}
}
