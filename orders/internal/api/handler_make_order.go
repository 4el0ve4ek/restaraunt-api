package api

import (
	stdhttp "net/http"
	"strconv"

	"github.com/pkg/errors"

	"github.com/4el0ve4ek/restaraunt-api/library/pkg/http"
	"orders/internal/context"
)

type inputMakeOrder = struct {
	DishIDToQuantity map[string]int `json:"dishes"`
	SpecialRequests  string         `json:"special_requests"`
}
type outputMakeOrder = struct {
	OrderID int `json:"order_id"`
}

type makeOrderManager interface {
	AddOrder(ctx context.Context, dishesToQuantity map[int]int, specialRequests string) (struct {
		OrderID           int
		Forbidden         bool
		NoSuchDishes      bool
		NotEnoughQuantity bool
	}, error)
}

func NewMakeOrderHandler(makeOrderManager makeOrderManager) *makeOrderHandler {
	return &makeOrderHandler{makeOrderManager: makeOrderManager}
}

type makeOrderHandler struct {
	makeOrderManager makeOrderManager
}

func (h *makeOrderHandler) ServeJSON(ctx context.Context, r *http.Request, input inputMakeOrder) (http.Response[outputMakeOrder], error) {
	var ret http.Response[outputMakeOrder]

	if !ctx.GetUser().IsPresent() {
		ret.StatusCode.Set(stdhttp.StatusUnauthorized)
		return ret, nil
	}

	dishIDToQuantity := make(map[int]int, len(input.DishIDToQuantity))
	for id, quantity := range input.DishIDToQuantity {
		dishID, err := strconv.Atoi(id)
		if err != nil {
			ret.StatusCode.Set(stdhttp.StatusBadRequest)
			return ret, errors.Wrap(err, "parse dishID")
		}
		dishIDToQuantity[dishID] = quantity
	}

	status, err := h.makeOrderManager.AddOrder(ctx, dishIDToQuantity, input.SpecialRequests)
	if err != nil {
		return ret, errors.Wrap(err, "add dish")
	}
	if status.Forbidden {
		ret.StatusCode.Set(stdhttp.StatusMethodNotAllowed)
		return ret, nil
	}
	if status.NotEnoughQuantity {
		ret.StatusCode.Set(stdhttp.StatusBadRequest)
		return ret, nil
	}
	if status.NoSuchDishes {
		ret.StatusCode.Set(stdhttp.StatusBadRequest)
		return ret, nil
	}
	ret.Content.Set(outputMakeOrder{
		OrderID: status.OrderID,
	})
	return ret, nil
}
