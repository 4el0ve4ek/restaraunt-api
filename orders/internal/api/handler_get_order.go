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

type inputGetOrder = http.None
type outputGetOrder = models.Order

type getOrderManager interface {
	GetOrderByID(ctx context.Context, orderID int) (models.Order, error)
}

func NewGetOrderHandler(getOrderManager getOrderManager) *getOrderHandler {
	return &getOrderHandler{getOrderManager: getOrderManager}
}

type getOrderHandler struct {
	getOrderManager getOrderManager
}

func (h *getOrderHandler) ServeJSON(ctx context.Context, r *http.Request, input inputGetOrder) (http.Response[outputGetOrder], error) {
	var ret http.Response[outputGetOrder]

	orderIDS := chi.URLParam(r, "orderID")
	orderID, err := strconv.Atoi(orderIDS)
	if err != nil {
		ret.StatusCode.Set(stdhttp.StatusBadRequest)
		return ret, nil
	}
	order, err := h.getOrderManager.GetOrderByID(ctx, orderID)
	if err != nil {
		return ret, errors.Wrap(err, "get order")
	}
	ret.Content.Set(order)
	return ret, nil
}
