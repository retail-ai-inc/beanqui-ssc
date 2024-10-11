package routers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/retail-ai-inc/beanqui/internal/jwtx"
	"github.com/retail-ai-inc/beanqui/internal/routers/consts"
	"github.com/retail-ai-inc/beanqui/internal/routers/response"
)

func Auth(next HandleFunc) HandleFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		result, cancelr := response.Get()
		defer cancelr()

		accept := request.Header.Get("Accept")

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
				result.Code = consts.InternalServerErrorCode
				result.Msg = "missing parameter"
				_ = result.Json(writer, http.StatusInternalServerError)
				return
			}

			token, err = jwtx.ParseHsToken(strs[1])

			if err != nil {
				result.Code = consts.InternalServerErrorCode
				result.Msg = err.Error()
				_ = result.Json(writer, http.StatusUnauthorized)
				return
			}
		}
		if auth == "" {
			auth = request.FormValue("token")
			token, err = jwtx.ParseHsToken(auth)

			if err != nil {
				result.Code = consts.InternalServerErrorMsg
				result.Msg = err.Error()
				_ = result.Json(writer, http.StatusUnauthorized)
				return
			}
		}

		//
		_, err = token.GetExpirationTime()
		if err != nil {
			result.Code = consts.AuthExpireCode
			result.Msg = err.Error()
			_ = result.Json(writer, http.StatusUnauthorized)
			return
		}

		next(writer, request)
	}
}
