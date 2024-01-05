package routers

import (
	"net/http"
	"strings"

	"github.com/retail-ai-inc/beanq/helper/json"
	"github.com/retail-ai-inc/beanqui/internal/redisx"
	"github.com/retail-ai-inc/beanqui/internal/routers/consts"
	"github.com/retail-ai-inc/beanqui/internal/routers/results"
	"github.com/spf13/cast"
)

type Logs struct {
}

func (t *Logs) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	resultRes, cancel := results.Get()
	defer cancel()

	client := redisx.Client()

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

	match := strings.Join([]string{matchStr, "*"}, ":")

	data := make(map[string]any)

	allKeys := client.Keys(r.Context(), match).Val()
	data["total"] = len(allKeys)

	keys, cursor, err := client.Scan(r.Context(), gCursor, match, 10).Result()
	if err != nil {
		resultRes.Code = "1005"
		resultRes.Msg = err.Error()
		_ = resultRes.Json(w, http.StatusInternalServerError)
		return
	}

	msgs := make([]*redisx.Msg, 0, 10)
	for _, key := range keys {

		str, err := client.Get(r.Context(), key).Result()
		if err != nil {
			continue
		}
		m := new(redisx.Msg)
		if err := json.Unmarshal([]byte(str), &m); err != nil {
			continue
		}
		msgs = append(msgs, m)
	}

	data["data"] = msgs
	data["cursor"] = cursor
	resultRes.Data = data

	_ = resultRes.Json(w, http.StatusOK)
	return
}
