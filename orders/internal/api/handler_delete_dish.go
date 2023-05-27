package api

import (
	stdhttp "net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"

	"github.com/4el0ve4ek/restaraunt-api/library/pkg/http"
	"orders/internal/context"
)

type inputDeleteDish = http.None
type outputDeleteDish = http.None

type deleteDishManager interface {
	DeleteDishByID(ctx context.Context, dishID int) (struct {
		Forbidden bool
		NoID      bool
	}, error)
}

func NewDeleteDishHandler(deleteDishManager deleteDishManager) *deleteDishHandler {
	return &deleteDishHandler{deleteDishManager: deleteDishManager}
}

type deleteDishHandler struct {
	deleteDishManager deleteDishManager
}

func (h *deleteDishHandler) ServeJSON(ctx context.Context, r *http.Request, input inputDeleteDish) (http.Response[outputDeleteDish], error) {
	var ret http.Response[outputDeleteDish]

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

	status, err := h.deleteDishManager.DeleteDishByID(ctx, dishID)
	if err != nil {
		return ret, errors.Wrap(err, "delete dish")
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
