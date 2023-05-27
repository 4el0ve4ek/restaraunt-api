package dishes

import (
	"context"

	"orders/internal/models"
)

type dishesRepository interface {
	GetAllDishes(context.Context) ([]models.Dish, error)
	GetAvailableDishes(context.Context) ([]models.Dish, error)
	UpdateDish(ctx context.Context, dish models.Dish) error
	DeleteDishByID(ctx context.Context, dishID int) error
	AddDish(ctx context.Context, dish models.Dish) error
}
