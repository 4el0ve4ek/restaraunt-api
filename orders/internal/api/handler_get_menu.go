package api

import (
	"github.com/pkg/errors"

	"github.com/4el0ve4ek/restaraunt-api/library/pkg/http"
	"orders/internal/context"
	"orders/internal/models"
)

type inputGetMenu = http.None
type outputGetMenu = struct {
	Dishes []models.Dish `json:"dishes"`
}

type getMenuManager interface {
	GetAllDishes(ctx context.Context) ([]models.Dish, error)
}

func NewGetMenuHandler(getMenuManager getMenuManager) *getMenuHandler {
	return &getMenuHandler{getMenuManager: getMenuManager}
}

type getMenuHandler struct {
	getMenuManager getMenuManager
}

func (h *getMenuHandler) ServeJSON(ctx context.Context, r *http.Request, input inputGetMenu) (http.Response[outputGetMenu], error) {
	var ret http.Response[outputGetMenu]
	dishes, err := h.getMenuManager.GetAllDishes(ctx)
	if err != nil {
		return ret, errors.Wrap(err, "get all dishes")
	}
	ret.Content.Set(outputGetMenu{
		Dishes: dishes,
	})
	return ret, nil
}
