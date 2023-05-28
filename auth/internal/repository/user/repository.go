package user

import (
	"context"

	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/4el0ve4ek/restaraunt-api/library/pkg/database/postgres"

	"auth/internal/models"
)

func NewRepository(db *postgres.DB) *repository {
	return &repository{db: db}
}

type repository struct {
	db *postgres.DB
}

func (r *repository) AddNewUser(ctx context.Context, username, email, encryptedPassword string) (bool, error) {
	_, err := r.db.ExecContext(
		ctx,
		`
		INSERT INTO "user"(username, email, password_hash, role)
 		VALUES($1, $2, $3, 'customer') 
 		`,
		username, email, encryptedPassword,
	)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Code == ("23505") {
			return true, nil
		}
		return false, errors.Wrap(err, "unable to add")
	}
	return false, nil
}

func (r *repository) GetUserWithID(ctx context.Context, userID int) (models.User, error) {
	row := r.db.QueryRowContext(
		ctx,
		`
		SELECT id, username, email, role, password_hash FROM "user" WHERE id = $1
 		`,
		userID,
	)
	if err := row.Err(); err != nil {
		return models.User{}, errors.Wrap(err, "unable to add")
	}
	var user models.User
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.HashedPassword); err != nil {
		return models.User{}, errors.Wrap(err, "unable to scan user id")
	}
	return user, nil
}
func (r *repository) GetUserWithEmail(ctx context.Context, email string) (models.User, error) {
	row := r.db.QueryRowContext(
		ctx,
		`
		SELECT id, username, email, role, password_hash FROM "user" WHERE email = $1
 		`,
		email,
	)
	if err := row.Err(); err != nil {
		return models.User{}, errors.Wrap(err, "unable to add")
	}
	var user models.User
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.HashedPassword); err != nil {
		return models.User{}, errors.Wrap(err, "unable to scan user id")
	}
	return user, nil
}
