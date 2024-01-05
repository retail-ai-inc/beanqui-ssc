package routers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/retail-ai-inc/beanq/helper/json"
	"github.com/retail-ai-inc/beanq/helper/stringx"
	"github.com/retail-ai-inc/beanqui/internal/redisx"
	"github.com/retail-ai-inc/beanqui/internal/routers/results"
)

type RedisInfo struct {
}

func (t *RedisInfo) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	result, cancel := results.Get()
	defer cancel()

	client := redisx.Client()

	// rq := ctx.Request()
	flusher, ok := w.(http.Flusher)
	if !ok {

	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	var b []byte

	for {
		d, err := redisx.Info(r.Context(), client)

		if err != nil {
			result.Code = "1001"
			result.Msg = err.Error()
			b, err = json.Marshal(result)

			_, err = w.Write(stringx.StringToByte(fmt.Sprintf("id:%d\n", time.Now().Unix())))
			_, err = w.Write([]byte("event:redis_info\n"))
			_, err = w.Write([]byte(fmt.Sprintf("data:%s\n\n", string(b))))

		}
		if err == nil {
			result.Data = d
			b, err = json.Marshal(result)

			_, err = w.Write([]byte(fmt.Sprintf("id:%d\n", time.Now().Unix())))
			_, err = w.Write([]byte("event:redis_info\n"))
			_, err = w.Write([]byte(fmt.Sprintf("data:%s\n\n", string(b))))

		}
		flusher.Flush()
		time.Sleep(10 * time.Second)
	}
}
