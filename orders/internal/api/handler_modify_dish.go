package api

import (
	stdhttp "net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"

	"github.com/4el0ve4ek/restaraunt-api/library/pkg/http"
	"orders/internal/context"
	"orders/internal/models"
)

type inputModifyDish = models.Dish
type outputModifyDish = http.None

type modifyDishManager interface {
	UpdateDish(ctx context.Context, dish models.Dish) (struct {
		Forbidden bool
		NoID      bool
	}, error)
}

func NewModifyDishHandler(modifyDishManager modifyDishManager) *modifyDishHandler {
	return &modifyDishHandler{modifyDishManager: modifyDishManager}
}

type modifyDishHandler struct {
	modifyDishManager modifyDishManager
}

func (h *modifyDishHandler) ServeJSON(ctx context.Context, r *http.Request, input inputModifyDish) (http.Response[outputModifyDish], error) {
	var ret http.Response[outputModifyDish]

	if !ctx.GetUser().IsPresent() {
		ret.StatusCode.Set(stdhttp.StatusUnauthorized)
		return ret, nil
	}

	dishIDS := chi.URLParam(r, "dishID")
	dishID, err := strconv.Atoi(dishIDS)
	if err != nil {
		ret.StatusCode.Set(stdhttp.StatusBadRequest)
		return ret, nil
	}
	input.ID = dishID

	status, err := h.modifyDishManager.UpdateDish(ctx, input)
	if err != nil {
		return ret, errors.Wrap(err, "modify dish")
	}
	if status.Forbidden {
		ret.StatusCode.Set(stdhttp.StatusMethodNotAllowed)
		return ret, nil
	}

	if status.NoID {
		ret.StatusCode.Set(stdhttp.StatusBadRequest)
		return ret, nil
	}

	return ret, nil
}
