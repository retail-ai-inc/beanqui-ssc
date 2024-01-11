package routers

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/retail-ai-inc/beanqui/internal/jwtx"
	"github.com/retail-ai-inc/beanqui/internal/routers/consts"
	"github.com/retail-ai-inc/beanqui/internal/routers/results"
	"github.com/spf13/viper"
)

type Login struct {
}

func NewLogin() *Login {
	return &Login{}
}

func (t *Login) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	result, cancel := results.Get()
	defer cancel()

	m := viper.GetStringMap("ui")
	user, pwd := "", ""
	if u, ok := m["username"].(string); ok {
		user = u
	}
	if p, ok := m["password"].(string); ok {
		pwd = p
	}

	if username != user && password != pwd {
		result.Code = consts.InternalServerErrorCode
		result.Msg = "username or password mismatch"
		_ = result.Json(w, http.StatusUnauthorized)
		return
	}

	claim := jwtx.Claim{
		UserName: username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    viper.GetString("issuer"),
			Subject:   viper.GetString("subject"),
			Audience:  nil,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(viper.GetDuration("expiresAt"))),
			NotBefore: nil,
			IssuedAt:  nil,
			ID:        "",
		},
	}

	token, err := jwtx.MakeHsToken(claim)
	if err != nil {
		result.Code = consts.InternalServerErrorCode
		result.Msg = err.Error()
		_ = result.Json(w, http.StatusInternalServerError)
		return
	}

	result.Data = map[string]any{"token": token}

	_ = result.Json(w, http.StatusOK)
	return

}
