package api

import (
	"context"

	"github.com/4el0ve4ek/restaraunt-api/library/pkg/http"
	"github.com/pkg/errors"
)

type inputLoginUser = struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type outputLoginUser = struct {
	JWTToken string `json:"jwt_token"`
}

type loginManager interface {
	LoginUser(ctx context.Context, email, userPassword string) (string, error)
}

func NewLoginUserHandler(loginManager loginManager) *loginUserHandler {
	return &loginUserHandler{loginManager: loginManager}
}

type loginUserHandler struct {
	loginManager loginManager
}

func (h *loginUserHandler) ServeJSON(r *http.Request, input inputLoginUser) (http.Response[outputLoginUser], error) {
	var ret http.Response[outputLoginUser]
	token, err := h.loginManager.LoginUser(r.Context(), input.Email, input.Password)
	if err != nil {
		return ret, errors.Wrap(err, "login user")
	}
	ret.Content.Set(outputLoginUser{
		JWTToken: token,
	})
	return ret, nil
}
