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

func (t *User) List(ctx *BeanContext) error {

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
		return res.Json(w, http.StatusOK)
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
	return res.Json(w, http.StatusOK)
}

func (t *User) Add(ctx *BeanContext) error {
	res, cancal := response.Get()
	defer cancal()

	r := ctx.Request
	w := ctx.Writer

	account := r.PostFormValue("account")
	password := r.PostFormValue("password")
	typ := r.PostFormValue("type")
	active := r.PostFormValue("active")
	detail := r.PostFormValue("detail")

	if account == "" {
		res.Code = errorx.MissParameterCode
		res.Msg = "missing account"
		return res.Json(w, http.StatusOK)

	}

	client := redisx.Client()
	key := strings.Join([]string{viper.GetString("redis.prefix"), "users", account}, ":")
	data := make(map[string]any, 0)
	data["account"] = account
	data["password"] = password
	data["type"] = typ
	data["active"] = active
	data["detail"] = detail

	if err := client.HSet(r.Context(), key, data).Err(); err != nil {
		res.Code = errorx.InternalServerErrorCode
		res.Msg = err.Error()
		return res.Json(w, http.StatusOK)

	}
	return res.Json(w, http.StatusOK)
}

func (t *User) Delete(ctx *BeanContext) error {

	res, cancel := response.Get()
	defer cancel()
	account := ctx.Request.FormValue("account")
	if account == "" {
		res.Code = errorx.MissParameterMsg
		res.Msg = "missing account field"
		return res.Json(ctx.Writer, http.StatusOK)

	}

	client := redisx.Client()
	key := strings.Join([]string{viper.GetString("redis.prefix"), "users", account}, ":")
	if err := client.Del(ctx.Request.Context(), key).Err(); err != nil {
		res.Code = errorx.InternalServerErrorCode
		res.Msg = err.Error()
		return res.Json(ctx.Writer, http.StatusOK)

	}
	return res.Json(ctx.Writer, http.StatusOK)

}

func (t *User) Edit(ctx *BeanContext) error {

	res, cancel := response.Get()
	defer cancel()

	r := ctx.Request
	w := ctx.Writer

	account := r.FormValue("account")
	password := r.FormValue("password")
	active := r.FormValue("active")
	typ := r.FormValue("type")
	detail := r.FormValue("detail")

	key := strings.Join([]string{viper.GetString("redis.prefix"), "users", account}, ":")
	var data = map[string]any{
		"account":  account,
		"password": password,
		"active":   active,
		"detail":   detail,
		"type":     typ,
	}
	client := redisx.Client()
	if err := client.HSet(r.Context(), key, data).Err(); err != nil {
		res.Code = errorx.InternalServerErrorMsg
		res.Msg = err.Error()
		return res.Json(w, http.StatusOK)

	}
	return res.Json(w, http.StatusOK)

}
