package routers

import (
	"net/http"
	"strings"

	"github.com/retail-ai-inc/beanqui/internal/jwtx"
	"github.com/retail-ai-inc/beanqui/internal/routers/consts"
	"github.com/retail-ai-inc/beanqui/internal/routers/results"
	"github.com/retail-ai-inc/beanqui/internal/simple_router"
)

func Auth(next simple_router.HandlerFunc) simple_router.HandlerFunc {

	return func(ctx *simple_router.Context) error {

		result, cancel := results.Get()
		defer cancel()
		req := ctx.Request()

		auth := req.Header.Get("Beanq-Authorization")

		strs := strings.Split(auth, " ")
		if len(strs) < 2 {
			// return data format err
			result.Code = consts.InternalServerErrorCode
			result.Msg = "missing parameter"
			return ctx.Json(http.StatusInternalServerError, result)
		}

		token, err := jwtx.ParseRsaToken(strs[1])
		if err != nil {
			result.Code = consts.InternalServerErrorCode
			result.Msg = err.Error()
			return ctx.Json(http.StatusUnauthorized, result)
		}
		//
		_, err = token.Claims.GetExpirationTime()
		if err != nil {
			result.Code = consts.InternalServerErrorCode
			result.Msg = err.Error()
			return ctx.Json(http.StatusUnauthorized, result)
		}

		return next(ctx)
	}
}
