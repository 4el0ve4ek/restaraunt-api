package api

import (
	"context"
	stdhttp "net/http"
	"strings"

	"github.com/pkg/errors"

	"github.com/4el0ve4ek/restaraunt-api/library/pkg/http"
	"github.com/4el0ve4ek/restaraunt-api/library/pkg/optional"

	"auth/internal/models"

	"auth/internal/manager/jwt"
)

type inputGetUser = http.None
type outputGetUser = models.User

type userGetManager interface {
	GetUserByToken(ctx context.Context, token string) (optional.Optional[models.User], error)
}

func NewGetUserHandler(userGetManager userGetManager) *getUserHandler {
	return &getUserHandler{userGetManager: userGetManager}
}

type getUserHandler struct {
	userGetManager userGetManager
}

func (h *getUserHandler) ServeJSON(r *http.Request, input inputGetUser) (http.Response[outputGetUser], error) {
	var ret http.Response[outputGetUser]
	header := r.Header.Get(jwt.Header)
	header = strings.TrimPrefix(header, "Bearer")
	token := strings.TrimSpace(header)

	user, err := h.userGetManager.GetUserByToken(r.Context(), token)
	if err != nil {
		return ret, errors.Wrap(err, "get user by token")
	}
	if !user.IsPresent() {
		ret.StatusCode.Set(stdhttp.StatusUnauthorized)
		return ret, nil
	}
	ret.Content.Set(user.Get())
	return ret, nil
}
