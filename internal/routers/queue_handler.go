package routers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/retail-ai-inc/beanqui/internal/redisx"
	"github.com/retail-ai-inc/beanqui/internal/routers/consts"
	"github.com/retail-ai-inc/beanqui/internal/routers/results"
	"github.com/spf13/viper"
)

type Queue struct {
}

func NewQueue() *Queue {
	return &Queue{}
}

func (t *Queue) List(w http.ResponseWriter, r *http.Request) {
	result, cancel := results.Get()
	defer cancel()

	// url like: queue?list&page=1&pageSize=10
	query := r.URL.RawQuery
	querys := strings.Split(query, "&")
	if len(querys) < 1 {
		result.Code = "1004"
		result.Msg = "404"
		_ = result.Json(w, http.StatusNotFound)
		return
	}
	bt, err := redisx.QueueInfo(r.Context())
	if err != nil {
		result.Code = consts.InternalServerErrorCode
		result.Msg = err.Error()
		_ = result.Json(w, http.StatusInternalServerError)
		return
	}

	result.Data = bt
	_ = result.Json(w, http.StatusOK)
	return

}
func (t *Queue) Detail(w http.ResponseWriter, r *http.Request) {
	queueDetail(w, r)
	return
}

func queueDetail(w http.ResponseWriter, r *http.Request) {

	result, cancel := results.Get()
	defer cancel()

	flusher, ok := w.(http.Flusher)
	if !ok {

	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	id := r.FormValue("id")
	prefix := viper.GetString("redis.prefix")
	id = strings.Join([]string{prefix, id, "normal_stream", "stream"}, ":")
	fmt.Printf("id值:%+v \n", id)
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
