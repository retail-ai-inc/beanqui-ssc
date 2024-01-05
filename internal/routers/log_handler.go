package routers

import (
	"context"
	"net/http"
	"strings"

	"github.com/retail-ai-inc/beanq/helper/json"
	"github.com/retail-ai-inc/beanqui/internal/redisx"
	"github.com/retail-ai-inc/beanqui/internal/routers/results"
)

type Log struct {
}

// del ,retry,archive,detail
func (t *Log) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	result, cancel := results.Get()
	defer cancel()

	// job detail
	if r.Method == http.MethodGet {

		id := r.FormValue("id")
		msgType := r.FormValue("msgType")

		if id == "" || msgType == "" {
			// error
		}
		data, err := detailHandler(r.Context(), id, msgType)
		if err != nil {
			result.Code = "1003"
			result.Msg = err.Error()
			_ = result.Json(w, http.StatusInternalServerError)
			return
		}
		result.Data = data
		_ = result.Json(w, http.StatusOK)
		return
	}
	if r.Method == http.MethodPost {
		return
	}
	// delete job
	if r.Method == http.MethodDelete {

		id := r.FormValue("id")
		msgType := r.FormValue("msgType")
		if id == "" || msgType == "" {
			// error
			return
		}

		if err := deleteHandler(r.Context(), id, msgType); err != nil {
			// error
			return
		}
		_ = result.Json(w, http.StatusOK)
		return
	}
	if r.Method == http.MethodPut {

	}
}

// log detail
func detailHandler(ctx context.Context, id, msgType string) (map[string]any, error) {

	client := redisx.Client()

	res, err := client.Get(ctx, strings.Join([]string{redisx.BqConfig.Redis.Prefix, "logs", msgType, id}, ":")).Result()
	if err != nil {
		return nil, err
	}
	m := make(map[string]any)
	if err := json.Unmarshal([]byte(res), &m); err != nil {
		return nil, err
	}
	return m, nil
}

func deleteHandler(ctx context.Context, id, msgType string) (err error) {

	client := redisx.Client()

	key := strings.Join([]string{redisx.BqConfig.Redis.Prefix, "logs", msgType, id}, ":")
	cmd := client.Del(ctx, key)

	if cmd.Err() != nil {
		return err
	}

	return nil
}
