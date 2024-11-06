package routers

import (
	"fmt"
	"github.com/retail-ai-inc/beanqui/internal/googleAuth"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/retail-ai-inc/beanqui/internal/jwtx"
	"github.com/retail-ai-inc/beanqui/internal/routers/errorx"
	"github.com/retail-ai-inc/beanqui/internal/routers/response"
	"github.com/spf13/viper"
)

type Login struct {
}

func NewLogin() *Login {
	return &Login{}
}

func (t *Login) Login(ctx *BeanContext) error {

	r := ctx.Request
	w := ctx.Writer

	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	result, cancel := response.Get()
	defer cancel()

	m := viper.GetStringMap("ui")
	user, pwd := "", ""
	if u, ok := m["username"].(string); ok {
		user = u
	}
	if p, ok := m["password"].(string); ok {
		pwd = p
	}

	if username != user || password != pwd {
		result.Code = errorx.InternalServerErrorCode
		result.Msg = "username or password mismatch"
		return result.Json(w, http.StatusUnauthorized)
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
		result.Code = errorx.InternalServerErrorCode
		result.Msg = err.Error()
		return result.Json(w, http.StatusInternalServerError)
	}

	result.Data = map[string]any{"token": token}

	return result.Json(w, http.StatusOK)

}

func (t *Login) GoogleLogin(ctx *BeanContext) error {
	w := ctx.Writer

	gAuth := googleAuth.New()
	state := "test_self"
	if state == "" {
		//do something
	}
	url := gAuth.AuthCodeUrl(state)
	w.Header().Set("Content-Type", "text/html;charset=UTF-8")
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
	return nil
}

func (t *Login) GoogleCallBack(ctx *BeanContext) error {

	r := ctx.Request
	w := ctx.Writer

	res, cancel := response.Get()
	defer cancel()

	state := r.FormValue("state")
	if state != "test_self" {
		//return ctx.Redirect(http.StatusTemporaryRedirect, "/login")
	}

	code := r.FormValue("code")
	auth := googleAuth.New()

	token, err := auth.Exchange(r.Context(), code)

	if err != nil {
		res.Code = errorx.InternalServerErrorCode
		res.Msg = err.Error()
		return res.Json(w, http.StatusOK)
	}

	userInfo, err := auth.Response(token.AccessToken)
	if err != nil {
		res.Code = errorx.InternalServerErrorCode
		res.Msg = err.Error()
		return res.Json(w, http.StatusOK)
	}
	//admin, err := handler.authService.CheckEmailByGoogleAuth(userInfo)
	//if err != nil {
	//	return ctx.JSON(http.StatusInternalServerError, echo.Map{"code": 10500, "msg": err.Error()})
	//}
	fmt.Printf("------token:%+v \n", userInfo)

	claim := jwtx.Claim{
		UserName: "111",
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
	jwtToken, err := jwtx.MakeHsToken(claim)
	if err != nil {
		res.Code = errorx.InternalServerErrorCode
		res.Msg = err.Error()
		return res.Json(w, http.StatusOK)
	}
	w.Header().Set("Content-Type", "text/html;charset=UTF-8")
	w.Header().Set("Location", "http://localhost:9090/#/login?token="+jwtToken)
	w.WriteHeader(http.StatusFound)
	return nil
}
