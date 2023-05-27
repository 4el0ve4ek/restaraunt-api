package order

import (
	stdcontext "context"

	"orders/internal/context"
	"orders/internal/models"
)

type orderRepository interface {
	AddOrder(ctx stdcontext.Context, order models.Order) (int, error)
	GetOrderByID(ctx stdcontext.Context, orderID int) (models.Order, error)
}

type dishesGetter interface {
	GetAllDishes(ctx context.Context) ([]models.Dish, error)
}
