package user

import (
	"context"

	"github.com/4el0ve4ek/restaraunt-api/library/pkg/optional"

	"auth/internal/models"
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
	AddNewUser(ctx context.Context, username, email, encryptedPassword string, role models.Role) (bool, error)
	GetUserWithID(ctx context.Context, userID int) (models.User, error)
	GetUserWithEmail(ctx context.Context, email string) (models.User, error)
}
