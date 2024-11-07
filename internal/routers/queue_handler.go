package routers

import (
	"net/http"
	"strings"
	"time"

	"github.com/retail-ai-inc/beanqui/internal/redisx"
	"github.com/retail-ai-inc/beanqui/internal/routers/errorx"
	"github.com/retail-ai-inc/beanqui/internal/routers/response"
	"github.com/spf13/viper"
)

type Queue struct {
}

func NewQueue() *Queue {
	return &Queue{}
}

func (t *Queue) List(ctx *BeanContext) error {
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
func (t *Queue) Detail(ctx *BeanContext) error {
	queueDetail(ctx.Writer, ctx.Request)
	return nil
}

func queueDetail(w http.ResponseWriter, r *http.Request) {

	result, cancel := response.Get()
	defer cancel()

	flusher, ok := w.(http.Flusher)
	if !ok {
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	id := r.FormValue("id")
	prefix := viper.GetString("redis.prefix")
	id = strings.Join([]string{prefix, id, "normal_stream", "stream"}, ":")

	ctx := r.Context()
	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()
	client := redisx.Client()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:

			cmd := client.XRangeN(ctx, id, "-", "+", 50)
			stream, err := cmd.Result()

			if err != nil {
				result.Code = "1004"
				result.Msg = err.Error()
			}

			if err == nil {
				result.Data = stream
			}
			_ = result.EventMsg(w, "queue_detail")
			flusher.Flush()
			ticker.Reset(10 * time.Second)
		}
	}
}
