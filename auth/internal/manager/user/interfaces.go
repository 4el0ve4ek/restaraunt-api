package user

import (
	"context"

	"auth/internal/models"
	"github.com/4el0ve4ek/restaraunt-api/library/pkg/optional"
)

type passwordManager interface {
	Encrypt(string) (string, error)
	CompareWithHashed(string, string) (bool, error)
}

type tokenManager interface {
	CreateToken(context.Context, int) (string, error)
	ExtractToken(context.Context, string) (optional.Optional[int], error)
}

type userRepository interface {
	AddNewUser(ctx context.Context, username, email, encryptedPassword string) (int, error)
	GetUserWithID(ctx context.Context, userID int) (models.User, error)
	GetUserWithEmail(ctx context.Context, email string) (models.User, error)
}
