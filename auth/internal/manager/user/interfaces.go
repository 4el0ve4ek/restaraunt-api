package user

import (
	"context"

	"auth/models"
)

type passwordManager interface {
	Encrypt(string) (string, error)
	CompareWithHashed(string, string) (bool, error)
}

type tokenManager interface {
	CreateToken(int) (string, error)
	ExtractToken(string) (int, error)
}

type userRepository interface {
	AddNewUser(ctx context.Context, username, email, encryptedPassword string) (int, error)
	GetUserWithID(ctx context.Context, userID int) (models.User, error)
	GetUserWithEmail(ctx context.Context, email string) (models.User, error)
}
