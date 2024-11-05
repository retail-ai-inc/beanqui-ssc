package routers

import (
	"fmt"
	"github.com/retail-ai-inc/beanqui/internal/redisx"
	"github.com/retail-ai-inc/beanqui/internal/routers/errorx"
	"github.com/retail-ai-inc/beanqui/internal/routers/response"
	"github.com/spf13/viper"
	"net/http"
	"strings"
)

type User struct {
}

func NewUser() *User {
	return &User{}
}

func (t *User) List(ctx *BeanContext) {

	res, cancel := response.Get()
	defer cancel()

	r := ctx.Request
	w := ctx.Writer

	pattern := strings.Join([]string{viper.GetString("redis.prefix"), "users:*"}, ":")
	client := redisx.Client()

	keys, err := client.Keys(r.Context(), pattern).Result()
	if err != nil {
		res.Code = errorx.InternalServerErrorMsg
		res.Msg = err.Error()
		_ = res.Json(w, http.StatusOK)
		return
	}

	data := make([]any, 0)
	for _, key := range keys {

		r, err := client.HGetAll(r.Context(), key).Result()
		if err != nil {
			fmt.Printf("hget err:%+v \n", err)
			continue
		}

		data = append(data, r)
	}
	res.Data = data
	_ = res.Json(w, http.StatusOK)
	return

}
