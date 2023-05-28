package session

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/4el0ve4ek/restaraunt-api/library/pkg/database/postgres"
	"github.com/4el0ve4ek/restaraunt-api/library/pkg/optional"
)

func NewRepository(db *postgres.DB) *repository {
	return &repository{db: db}
}

type repository struct {
	db *postgres.DB
}

func (r *repository) AddSession(ctx context.Context, userID int, sessionToken string, createdAt time.Time) error {
	_, err := r.db.ExecContext(
		ctx,
		`
		INSERT INTO session(user_id, session_token, expires_at) VALUES($1, $2, $3);
		`,
		userID, sessionToken, createdAt,
	)
	return errors.Wrap(err, "insert into db")
}

func (r *repository) GetSession(ctx context.Context, sessionToken string) (optional.Optional[int], error) {
	rows, err := r.db.QueryContext(
		ctx,
		`
		SELECT user_id FROM session WHERE session_token = $1 and expires_at >= $2;
		`,
		sessionToken, time.Now(),
	)

	var ret optional.Optional[int]
	if err != nil {
		return ret, errors.Wrap(err, "select from db")
	}

	if !rows.Next() {
		return ret, nil
	}

	var id int
	if err := rows.Scan(&id); err != nil {
		return ret, errors.Wrap(err, "scan db response")
	}

	ret.Set(id)
	return ret, nil
}
