package dishes

import (
	"github.com/pkg/errors"

	"orders/internal/context"
	"orders/internal/models"
	"orders/internal/repository/dishes"
)

type Manager interface {
	GetAllDishes(ctx context.Context) ([]models.Dish, error)
	UpdateDish(ctx context.Context, dish models.Dish) (struct {
		Forbidden bool
		NoID      bool
	}, error)
	DeleteDishByID(ctx context.Context, dishID int) (struct {
		Forbidden bool
		NoID      bool
	}, error)
	AddDish(ctx context.Context, dish models.Dish) (struct {
		Forbidden bool
	}, error)
}

func NewManager(dishesRepository dishesRepository) *manager {
	return &manager{
		dishesRepository: dishesRepository,
	}
}

type manager struct {
	dishesRepository dishesRepository
}

func (m *manager) GetAllDishes(ctx context.Context) ([]models.Dish, error) {
	user := ctx.GetUser()
	if !user.IsPresent() || user.Get().Role == models.Customer {
		return m.dishesRepository.GetAvailableDishes(ctx)
	}
	return m.dishesRepository.GetAllDishes(ctx)
}

func (m *manager) UpdateDish(ctx context.Context, dish models.Dish) (struct {
	Forbidden bool
	NoID      bool
}, error) {
	var ret struct {
		Forbidden bool
		NoID      bool
	}

	if ctx.GetUser().Get().Role != models.Manager {
		ret.Forbidden = true
		return ret, nil
	}

	err := m.dishesRepository.UpdateDish(ctx, dish)
	if errors.Is(err, dishes.ErrNoSuchID) {
		ret.NoID = true
		return ret, nil
	} else if err != nil {
		return ret, errors.Wrap(err, "modify dish")
	}
	return ret, nil
}

func (m *manager) DeleteDishByID(ctx context.Context, dishID int) (struct {
	Forbidden bool
	NoID      bool
}, error) {
	var ret struct {
		Forbidden bool
		NoID      bool
	}
	if ctx.GetUser().Get().Role != models.Manager {
		ret.Forbidden = true
		return ret, nil
	}

	err := m.dishesRepository.DeleteDishByID(ctx, dishID)
	if errors.Is(err, dishes.ErrNoSuchID) {
		ret.NoID = true
		return ret, nil
	} else if err != nil {
		return ret, errors.Wrap(err, "delete dish")
	}
	return ret, nil
}

func (m *manager) AddDish(ctx context.Context, dish models.Dish) (struct {
	Forbidden bool
}, error) {
	var ret struct {
		Forbidden bool
	}

	if ctx.GetUser().Get().Role != models.Manager {
		ret.Forbidden = true
		return ret, nil
	}
	return ret, m.dishesRepository.AddDish(ctx, dish)
}
