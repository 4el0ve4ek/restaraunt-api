package order

import (
	"context"

	"github.com/pkg/errors"

	"github.com/4el0ve4ek/restaraunt-api/library/pkg/database/postgres"
	"orders/internal/models"
)

func NewRepository(db *postgres.DB) *repository {
	return &repository{db: db}
}

type repository struct {
	db *postgres.DB
}

func (r *repository) AddOrder(ctx context.Context, order models.Order) (int, error) {
	row := r.db.QueryRowContext(
		ctx,
		`
		INSERT INTO "order"(user_id, status, special_requests)
		VALUES ($1, $2, $3)
		RETURNING id
		`,
		order.UserID, order.Status, order.SpecialRequests,
	)
	if err := row.Err(); err != nil {
		return 0, errors.Wrap(err, "insert into order")
	}
	var orderID int
	if err := row.Scan(&orderID); err != nil {
		return 0, errors.Wrap(err, "scan order id")
	}
	for _, orderDish := range order.Dishes {
		_, err := r.db.ExecContext(
			ctx,
			`
			INSERT INTO "order_dish"(order_id, dish_id, quantity, price)
			VALUES ($1, $2, $3, $4)
			`,
			orderID, orderDish.DishID, orderDish.Quantity, orderDish.Price,
		)
		if err != nil {
			return 0, errors.Wrap(err, "insert into order_dish") // TODO: use tx with rollback
		}
	}
	return orderID, nil
}

func (r *repository) GetOrderByID(ctx context.Context, orderID int) (models.Order, error) {
	row := r.db.QueryRowContext(
		ctx,
		`
		SELECT id, user_id, status, special_requests FROM "order" WHERE id = $1
		`,
		orderID,
	)
	if err := row.Err(); err != nil {
		return models.Order{}, errors.Wrap(err, "select from order")
	}

	var order models.Order
	if err := row.Scan(&order.ID, &order.UserID, &order.Status, &order.SpecialRequests); err != nil {
		return models.Order{}, errors.Wrap(err, "scan row")
	}

	rows, err := r.db.QueryContext(
		ctx,
		`
		SELECT dish_id, quantity, price FROM "order_dish" WHERE order_id = $1
		`,
		orderID,
	)
	if err != nil {
		return models.Order{}, errors.Wrap(err, "scan rows")
	}
	for rows.Next() {
		var orderDish models.OrderDish
		err := rows.Scan(&orderDish.DishID, &orderDish.Quantity, &orderDish.Price)
		if err != nil {
			return models.Order{}, errors.Wrap(err, "scan row")
		}
		order.Dishes = append(order.Dishes, orderDish)
	}
	rows.Close()
	return order, nil
}

func (r *repository) GetWaitingOrdersID(ctx context.Context) ([]int, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`
		SELECT id FROM "order" WHERE status = 'waiting';
		`,
	)
	if err != nil {
		return []int{}, errors.Wrap(err, "select waiting")
	}
	ret := make([]int, 0)
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return []int{}, errors.Wrap(err, "scan order id")
		}
		ret = append(ret, id)
	}
	return ret, nil
}

func (r *repository) SetProcessingOrderByID(ctx context.Context, orderID int) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE "order" SET status = 'processing' WHERE id = $1`,
		orderID,
	)
	return errors.Wrap(err, "update status")
}

func (r *repository) SetSuccessOrderByID(ctx context.Context, orderID int) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE "order" SET status = 'success' WHERE id = $1`,
		orderID,
	)
	return errors.Wrap(err, "update status")
}

func (r *repository) SetCancelOrderByID(ctx context.Context, orderID int) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE "order" SET status = 'cancel' WHERE id = $1`,
		orderID,
	)
	return errors.Wrap(err, "update status")
}

func (r *repository) ReduceDishQuantity(ctx context.Context, dishID int, quantity int) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE "dish" SET quantity = quantity - $2 WHERE id = $1`,
		dishID, quantity,
	)
	return errors.Wrap(err, "update quantity")
}
