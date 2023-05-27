package api

import (
	"context"

	"github.com/pkg/errors"

	"github.com/4el0ve4ek/restaraunt-api/library/pkg/http"
)

type inputRegisterUser = struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type outputRegisterUser = struct {
	JWTToken string `json:"jwt_token"`
}

type registerManager interface {
	RegisterUser(ctx context.Context, username, email, userPassword string) (string, error)
}

func NewRegisterUserHandler(registerManager registerManager) *registerUserHandler {
	return &registerUserHandler{registerManager: registerManager}
}

type registerUserHandler struct {
	registerManager registerManager
}

func (h *registerUserHandler) ServeJSON(r *http.Request, input inputRegisterUser) (http.Response[outputRegisterUser], error) {
	var ret http.Response[outputRegisterUser]
	token, err := h.registerManager.RegisterUser(r.Context(), input.Username, input.Email, input.Password)
	if err != nil {
		return ret, errors.Wrap(err, "register user")
	}
	ret.Content.Set(outputRegisterUser{
		JWTToken: token,
	})
	return ret, nil
}
