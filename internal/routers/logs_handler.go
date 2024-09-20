package routers

import (
	"net/http"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/retail-ai-inc/beanq/helper/json"
	"github.com/retail-ai-inc/beanqui/internal/redisx"
	"github.com/retail-ai-inc/beanqui/internal/routers/consts"
	"github.com/retail-ai-inc/beanqui/internal/routers/results"
	"github.com/spf13/cast"
)

type Logs struct {
	client redis.UniversalClient
}

func NewLogs(client redis.UniversalClient) *Logs {
	return &Logs{client: client}
}

func (t *Logs) List(w http.ResponseWriter, r *http.Request) {

	resultRes, cancel := results.Get()
	defer cancel()

	var (
		dataType string = "success"
		matchStr string = strings.Join([]string{redisx.BqConfig.Redis.Prefix, "logs", "success"}, ":")
	)

	dataType = r.FormValue("type")
	gCursor := cast.ToUint64(r.FormValue("cursor"))

	if dataType != "success" && dataType != "error" {
		resultRes.Code = consts.TypeErrorCode
		resultRes.Msg = consts.TypeErrorMsg

		_ = resultRes.Json(w, http.StatusInternalServerError)
		return

	}

	if dataType == "error" {
		matchStr = strings.Join([]string{redisx.BqConfig.Redis.Prefix, "logs", "fail"}, ":")
	}
	data := make(map[string]any)

	count := t.client.ZCard(r.Context(), matchStr).Val()
	data["total"] = count

	cmd := t.client.ZScan(r.Context(), matchStr, gCursor, "", 10)

	keys, cursor, err := cmd.Result()

	if err != nil {
		resultRes.Code = "1005"
		resultRes.Msg = err.Error()
		_ = resultRes.Json(w, http.StatusInternalServerError)
		return
	}

	msgs := make([]*redisx.Msg, 0, 10)
	m := new(redisx.Msg)

	for _, key := range keys {

		if err := json.Unmarshal([]byte(key), &m); err != nil {
			m.Score = key
			msgs = append(msgs, m)
			m = nil
		}

	}

	data["data"] = msgs
	data["cursor"] = cursor
	resultRes.Data = data

	_ = resultRes.Json(w, http.StatusOK)
	return
}
