package api

import (
	"context"
	stdhttp "net/http"

	"github.com/4el0ve4ek/restaraunt-api/library/pkg/optional"
	"github.com/pkg/errors"

	"github.com/4el0ve4ek/restaraunt-api/library/pkg/http"
)

type inputRegisterUser = struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type outputRegisterUser = struct {
	Message string `json:"message"`
}

type registerManager interface {
	RegisterUser(ctx context.Context, username, email, userPassword string, role optional.Optional[string]) (struct {
		FieldsCollide bool
		InvalidEmail  bool
	}, error)
}

func NewRegisterUserHandler(registerManager registerManager) *registerUserHandler {
	return &registerUserHandler{registerManager: registerManager}
}

type registerUserHandler struct {
	registerManager registerManager
}

func (h *registerUserHandler) ServeJSON(r *http.Request, input inputRegisterUser) (http.Response[outputRegisterUser], error) {
	var ret http.Response[outputRegisterUser]

	var queryRole optional.Optional[string]
	query := r.URL.Query()
	if query.Has("role") {
		queryRole.Set(query.Get("role"))
	}

	status, err := h.registerManager.RegisterUser(r.Context(), input.Username, input.Email, input.Password, queryRole)
	if err != nil {
		return ret, errors.Wrap(err, "register user")
	}
	if status.InvalidEmail {
		ret.StatusCode.Set(stdhttp.StatusBadRequest)
		ret.Content.Set(outputRegisterUser{
			Message: "invalid email",
		})
		return ret, nil
	}
	if status.FieldsCollide {
		ret.StatusCode.Set(stdhttp.StatusConflict)
		ret.Content.Set(outputRegisterUser{
			Message: "email or username collides",
		})
		return ret, nil
	}
	ret.Content.Set(outputRegisterUser{
		Message: "success",
	})
	return ret, nil
}
