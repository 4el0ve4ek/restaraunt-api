package context

import (
	stdcontext "context"
	"net/http"

	"github.com/pkg/errors"

	"github.com/4el0ve4ek/restaraunt-api/library/pkg/optional"
	"orders/internal/models"
	"orders/internal/services/auth"
)

type Context interface {
	stdcontext.Context
	GetUser() optional.Optional[models.User]
}

func NewContext(
	request *http.Request,
	authService auth.Service,
) (*context, error) {
	ret := &context{
		Context: request.Context(),
	}
	if header := request.Header.Get(auth.AuthorizationHeader); header != "" {
		user, err := authService.Get(ret, header)
		if err != nil {
			return ret, errors.Wrap(err, "auth error")
		}
		ret.user.Set(user)
	}

	return ret, nil
}

type context struct {
	stdcontext.Context

	user optional.Optional[models.User]
}

func (c *context) GetUser() optional.Optional[models.User] {
	return c.user
}
