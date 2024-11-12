package routers

import (
	"github.com/spf13/cast"
	"net/http"
	"runtime"
	"strings"

	"github.com/retail-ai-inc/beanqui/internal/redisx"
	"github.com/retail-ai-inc/beanqui/internal/routers/errorx"
	"github.com/retail-ai-inc/beanqui/internal/routers/response"
	"github.com/spf13/viper"
)

type Dashboard struct {
}

func NewDashboard() *Dashboard {
	return &Dashboard{}
}

func (t *Dashboard) Info(ctx *BeanContext) error {

	result, cancel := response.Get()
	defer cancel()

	w := ctx.Writer
	r := ctx.Request

	server, err := redisx.Server(r.Context())
	if err != nil {
		result.Code = errorx.InternalServerErrorCode
		result.Msg = err.Error()
		return result.Json(w, http.StatusOK)
	}
	persistence, err := redisx.Persistence(r.Context())
	if err != nil {
		result.Code = errorx.InternalServerErrorCode
		result.Msg = err.Error()
		return result.Json(w, http.StatusOK)
	}
	memory, err := redisx.Memory(r.Context())
	if err != nil {
		result.Code = errorx.InternalServerErrorCode
		result.Msg = err.Error()
		return result.Json(w, http.StatusOK)
	}

	command, err := redisx.CommandStats(r.Context())
	if err != nil {
		result.Code = errorx.InternalServerErrorCode
		result.Msg = err.Error()
		return result.Json(w, http.StatusOK)
	}

	clients, err := redisx.Clients(r.Context())
	if err != nil {
		result.Code = errorx.InternalServerErrorCode
		result.Msg = err.Error()
		return result.Json(w, http.StatusOK)
	}
	stats, err := redisx.Stats(r.Context())
	if err != nil {
		result.Code = errorx.InternalServerErrorCode
		result.Msg = err.Error()
		return result.Json(w, http.StatusOK)
	}

	keyspace, err := redisx.KeySpace(r.Context())
	if err != nil {
		result.Code = errorx.InternalServerErrorCode
		result.Msg = err.Error()
		return result.Json(w, http.StatusOK)
	}

	numCpu := runtime.NumCPU()

	// get queue total
	keys, err := redisx.Keys(r.Context(), strings.Join([]string{redisx.BqConfig.Redis.Prefix, "*", "stream"}, ":"))
	if err != nil {
		result.Code = errorx.InternalServerErrorCode
		result.Msg = err.Error()
		return result.Json(w, http.StatusInternalServerError)

	}
	keysLen := len(keys)

	// db size
	dbSize, err := redisx.DbSize(r.Context())
	if err != nil {
		result.Code = errorx.InternalServerErrorCode
		result.Msg = err.Error()
		return result.Json(w, http.StatusInternalServerError)
	}

	// Queue Past 10 Minutes
	prefix := viper.GetString("redis.prefix")
	failKey := strings.Join([]string{prefix, "logs", "fail"}, ":")
	failCount := redisx.ZCard(r.Context(), failKey)

	successKey := strings.Join([]string{prefix, "logs", "success"}, ":")
	successCount := redisx.ZCard(r.Context(), successKey)

	result.Data = map[string]any{
		"queue_total":   keysLen,
		"db_size":       dbSize,
		"num_cpu":       numCpu,
		"fail_count":    failCount,
		"success_count": successCount,
		"used_memory":   cast.ToInt(memory["used_memory_rss"]) / 1024 / 1024,
		"total_memory":  cast.ToInt(memory["total_system_memory"]) / 1024 / 1024,
		"commands":      command,
		"clients":       clients,
		"stats":         stats,
		"keyspace":      keyspace,
		"memory":        memory,
		"server":        server,
		"persistence":   persistence,
	}
	return result.Json(w, http.StatusOK)
}
