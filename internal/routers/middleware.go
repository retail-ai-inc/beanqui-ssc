package routers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/retail-ai-inc/beanqui/internal/jwtx"
	"github.com/retail-ai-inc/beanqui/internal/routers/errorx"
	"github.com/retail-ai-inc/beanqui/internal/routers/response"
)

func MigrateMiddleWare(next HandleFunc) HandleFunc {
	return HeaderRule(Auth(next))
}

func HeaderRule(next HandleFunc) HandleFunc {
	return func(ctx *BeanContext) error {
		ctx.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		ctx.Writer.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline';")
		ctx.Writer.Header().Set("X-Frame-Options", "SAMEORIGIN")
		ctx.Writer.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		return next(ctx)
	}
}

func Auth(next HandleFunc) HandleFunc {
	return func(ctx *BeanContext) error {

		result, cancelr := response.Get()
		defer cancelr()
		request := ctx.Request
		writer := ctx.Writer

		accept := request.Header.Get("Accept")
		//for SSE
		if !strings.EqualFold(accept, "text/event-stream") {
			ctx, cancel := context.WithTimeout(request.Context(), 20*time.Second)
			defer cancel()
			request = request.WithContext(ctx)
		}

		var (
			err   error
			token *jwtx.Claim
		)

		auth := request.Header.Get("Beanq-Authorization")
		if auth != "" {
			strs := strings.Split(auth, " ")
			if len(strs) < 2 {
				// return data format err
				result.Code = errorx.InternalServerErrorCode
				result.Msg = "missing parameter"
				return result.Json(writer, http.StatusInternalServerError)

			}

			token, err = jwtx.ParseHsToken(strs[1])

			if err != nil {
				result.Code = errorx.InternalServerErrorCode
				result.Msg = err.Error()
				return result.Json(writer, http.StatusUnauthorized)

			}
		}
		if auth == "" {
			auth = request.FormValue("token")
			token, err = jwtx.ParseHsToken(auth)
			if err != nil {
				result.Code = errorx.InternalServerErrorMsg
				result.Msg = err.Error()
				return result.Json(writer, http.StatusUnauthorized)

			}
		}

		//
		_, err = token.GetExpirationTime()
		if err != nil {
			result.Code = errorx.AuthExpireCode
			result.Msg = err.Error()
			return result.Json(writer, http.StatusUnauthorized)

		}
		return next(ctx)
	}
}
