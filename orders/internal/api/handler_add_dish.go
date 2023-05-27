package api

import (
	stdhttp "net/http"

	"github.com/pkg/errors"

	"github.com/4el0ve4ek/restaraunt-api/library/pkg/http"
	"orders/internal/context"
	"orders/internal/models"
)

type inputAddDish = models.Dish
type outputAddDish = http.None

type addDishManager interface {
	AddDish(ctx context.Context, dish models.Dish) (struct {
		Forbidden bool
	}, error)
}

func NewAddDishHandler(addDishManager addDishManager) *addDishHandler {
	return &addDishHandler{addDishManager: addDishManager}
}

type addDishHandler struct {
	addDishManager addDishManager
}

func (h *addDishHandler) ServeJSON(ctx context.Context, r *http.Request, input inputAddDish) (http.Response[outputAddDish], error) {
	var ret http.Response[outputAddDish]

	if !ctx.GetUser().IsPresent() {
		ret.StatusCode.Set(stdhttp.StatusUnauthorized)
		return ret, nil
	}

	status, err := h.addDishManager.AddDish(ctx, input)
	if err != nil {
		return ret, errors.Wrap(err, "add dish")
	}
	if status.Forbidden {
		ret.StatusCode.Set(stdhttp.StatusMethodNotAllowed)
	}
	return ret, nil
}
