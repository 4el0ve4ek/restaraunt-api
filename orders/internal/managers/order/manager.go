package order

import (
	"github.com/pkg/errors"

	"orders/internal/context"
	"orders/internal/models"
)

type Manager interface {
	AddOrder(ctx context.Context, dishesToQuantity map[int]int, specialRequests string) (struct {
		OrderID           int
		Forbidden         bool
		NoSuchDishes      bool
		NotEnoughQuantity bool
	}, error)
	GetOrderByID(ctx context.Context, orderID int) (models.Order, error)
}

func NewManager(orderRepository orderRepository, dishesGetter dishesGetter) *manager {
	return &manager{
		orderRepository: orderRepository,
		dishesGetter:    dishesGetter,
	}
}

type manager struct {
	orderRepository orderRepository
	dishesGetter    dishesGetter
}

func (m *manager) AddOrder(ctx context.Context, dishesToQuantity map[int]int, specialRequests string) (struct {
	OrderID           int
	Forbidden         bool
	NoSuchDishes      bool
	NotEnoughQuantity bool
}, error) {
	var ret struct {
		OrderID           int
		Forbidden         bool
		NoSuchDishes      bool
		NotEnoughQuantity bool
	}
	if !ctx.GetUser().IsPresent() {
		ret.Forbidden = true
		return ret, nil
	}

	order := models.Order{
		Status:          "waiting",
		UserID:          ctx.GetUser().Get().ID,
		SpecialRequests: specialRequests,
		Dishes:          make([]models.OrderDish, 0, len(dishesToQuantity)),
	}

	menuDishes, err := m.dishesGetter.GetAllDishes(ctx)
	if err != nil {
		return ret, errors.Wrap(err, "get all dishes")
	}

	for _, dish := range menuDishes {
		quantity, ok := dishesToQuantity[dish.ID]
		if !ok {
			continue
		}
		if quantity > dish.Quantity {
			ret.NotEnoughQuantity = true
			return ret, nil
		}
		order.Dishes = append(order.Dishes, models.OrderDish{
			DishID:   dish.ID,
			Quantity: quantity,
			Price:    dish.Price,
		})
	}
	if len(order.Dishes) != len(dishesToQuantity) {
		ret.NoSuchDishes = true
		return ret, nil
	}

	orderID, err := m.orderRepository.AddOrder(ctx, order)
	if err != nil {
		return ret, errors.Wrap(err, "add order to db")
	}
	ret.OrderID = orderID
	return ret, nil
}

func (m *manager) GetOrderByID(ctx context.Context, orderID int) (models.Order, error) {
	return m.orderRepository.GetOrderByID(ctx, orderID)
}
