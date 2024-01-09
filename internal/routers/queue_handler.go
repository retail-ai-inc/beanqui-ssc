package routers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/retail-ai-inc/beanq/helper/json"
	"github.com/retail-ai-inc/beanqui/internal/redisx"
	"github.com/retail-ai-inc/beanqui/internal/routers/consts"
	"github.com/retail-ai-inc/beanqui/internal/routers/results"
	"github.com/spf13/viper"
)

type Queue struct {
}

func (t *Queue) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	client := redisx.Client()
	action := querys[0]
	if r.Method == http.MethodGet {
		// queue list
		if action == "list" {
			bt, err := redisx.QueueInfo(r.Context(), client, redisx.QueueKey(redisx.BqConfig.Redis.Prefix))
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
		// queue detail
		if action == "detail" {
			queueDetail(w, r, client)
		}
	}
}

func queueDetail(w http.ResponseWriter, r *http.Request, client *redis.Client) {

	result, cancel := results.Get()
	defer cancel()

	flusher, ok := w.(http.Flusher)
	if !ok {

	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	var b []byte

	id := r.FormValue("id")
	prefix := viper.GetString("redis.prefix")
	id = strings.Join([]string{prefix, id, "stream"}, ":")

	for {
		ctx := context.TODO()
		ctx, _ = context.WithTimeout(ctx, 10*time.Second)

		cmd, err := client.XInfoStreamFull(ctx, id, 10).Result()
		if err != nil {
			fmt.Println(err)
		}
		if err == nil {
			result.Data = cmd.Entries
			b, err = json.Marshal(result)

			_, err = w.Write([]byte(fmt.Sprintf("id:%d\n", time.Now().Unix())))
			_, err = w.Write([]byte("event:queue_detail\n"))
			_, err = w.Write([]byte(fmt.Sprintf("data:%s\n\n", string(b))))

		}
		flusher.Flush()
		time.Sleep(10 * time.Second)
	}
}
