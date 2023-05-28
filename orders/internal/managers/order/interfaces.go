package order

import (
	stdcontext "context"

	"orders/internal/models"
)

type orderRepository interface {
	AddOrder(ctx stdcontext.Context, order models.Order) (int, error)
	GetOrderByID(ctx stdcontext.Context, orderID int) (models.Order, error)
	GetWaitingOrdersID(ctx stdcontext.Context) ([]int, error)
	SetSuccessOrderByID(ctx stdcontext.Context, orderID int) error
	SetProcessingOrderByID(ctx stdcontext.Context, orderID int) error
	SetCancelOrderByID(ctx stdcontext.Context, orderID int) error
	ReduceDishQuantity(ctx stdcontext.Context, id int, quantity int) error
}

type dishesGetter interface {
	GetAllDishes(ctx stdcontext.Context) ([]models.Dish, error)
}
