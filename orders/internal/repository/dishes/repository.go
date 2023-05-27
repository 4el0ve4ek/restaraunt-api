package dishes

import (
	"context"
	stderrors "errors"

	"github.com/pkg/errors"

	"github.com/4el0ve4ek/restaraunt-api/library/pkg/database/postgres"
	"orders/internal/models"
)

var ErrNoSuchID error = stderrors.New("no dish with such id")

func NewRepository(db *postgres.DB) *repository {
	return &repository{db: db}
}

type repository struct {
	db *postgres.DB
}

func (r *repository) GetAllDishes(ctx context.Context) ([]models.Dish, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`
		SELECT id, name, description, price, quantity, is_available FROM dish
		`,
	)
	if err != nil {
		return []models.Dish{}, errors.Wrap(err, "do query")
	}
	if !rows.NextResultSet() {
		return []models.Dish{}, nil
	}
	ret := make([]models.Dish, 0)
	for rows.Next() {
		var dish models.Dish
		err := rows.Scan(&dish.ID, &dish.Description, &dish.Price, &dish.Quantity, &dish.Available)
		if err != nil {
			return []models.Dish{}, errors.Wrap(err, "scan row")
		}

		if dish.Quantity == 0 {
			dish.Available = false
		}

		ret = append(ret, dish)
	}
	return ret, nil
}
func (r *repository) GetAvailableDishes(ctx context.Context) ([]models.Dish, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`
		SELECT id, name, description, price, quantity, is_available 
		FROM dish
		WHERE is_available = true AND quantity != 0
		`,
	)
	if err != nil {
		return []models.Dish{}, errors.Wrap(err, "do query")
	}
	if !rows.NextResultSet() {
		return []models.Dish{}, nil
	}
	ret := make([]models.Dish, 0)
	for rows.Next() {
		var dish models.Dish
		err := rows.Scan(&dish.ID, &dish.Description, &dish.Price, &dish.Quantity, &dish.Available)
		if err != nil {
			return []models.Dish{}, errors.Wrap(err, "scan row")
		}
		ret = append(ret, dish)
	}
	return ret, nil
}

func (r *repository) UpdateDish(ctx context.Context, dish models.Dish) error {
	result, err := r.db.ExecContext(
		ctx,
		`
		UPDATE dish SET name = $2, description = $3, price = $4, quantity = $5, is_available = $6 
		WHERE id = $1
		`,
		dish.ID, dish.Name, dish.Description, dish.Price, dish.Quantity, dish.Available,
	)
	if err != nil {
		return errors.Wrap(err, "exec result")
	}

	if rows, _ := result.RowsAffected(); rows != 1 {
		return ErrNoSuchID
	}

	return nil
}

func (r *repository) DeleteDishByID(ctx context.Context, dishID int) error {
	result, err := r.db.ExecContext(
		ctx,
		`
		DELETE FROM dish 
		WHERE id = $1
		`,
		dishID,
	)
	if err != nil {
		return errors.Wrap(err, "exec result")
	}

	if rows, _ := result.RowsAffected(); rows != 1 {
		return ErrNoSuchID
	}

	return nil
}

func (r *repository) AddDish(ctx context.Context, dish models.Dish) error {
	_, err := r.db.ExecContext(
		ctx,
		`
		INSERT INTO dish (name, desicription, price, quantity, is_available)
		VALUES ($1, $2, $3, $4, $5)
		`,
		dish.Name, dish.Description, dish.Price, dish.Quantity, dish.Available,
	)
	if err != nil {
		return errors.Wrap(err, "exec result")
	}
	return nil
}
