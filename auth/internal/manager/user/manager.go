package user

import (
	"context"
	"regexp"

	"github.com/pkg/errors"

	"auth/models"
)

type Manager interface {
	RegisterUser(ctx context.Context, username, email, userPassword string) (string, error)
	LoginUser(ctx context.Context, email, userPassword string) (string, error)
	GetUserByToken(ctx context.Context, token string) (models.User, error)
}

func NewManager(userRepository userRepository, passwordManager passwordManager, jwtManager tokenManager) *manager {
	return &manager{
		passwordManager: passwordManager,
		jwtManager:      jwtManager,
		userRepository:  userRepository,
	}
}

type manager struct {
	passwordManager passwordManager
	jwtManager      tokenManager
	userRepository  userRepository
}

func (m *manager) RegisterUser(ctx context.Context, username, email, userPassword string) (string, error) {
	passwordEncrypted, err := m.passwordManager.Encrypt(userPassword)
	if err != nil {
		return "", errors.Wrap(err, "encrypt password")
	}

	if !m.validateEmail(email) {
		return "", errors.Wrap(err, "invalid email")
	}

	userID, err := m.userRepository.AddNewUser(ctx, username, email, passwordEncrypted)
	if err != nil {
		return "", errors.Wrap(err, "add new user")
	}
	token, err := m.jwtManager.CreateToken(userID)
	if err != nil {
		return "", errors.Wrap(err, "create token")
	}
	return token, nil
}

func (m *manager) LoginUser(ctx context.Context, email, userPassword string) (string, error) {
	user, err := m.userRepository.GetUserWithEmail(ctx, email)
	if err != nil {
		return "", errors.Wrap(err, "get user by email")
	}

	equal, err := m.passwordManager.CompareWithHashed(user.HashedPassword, userPassword)
	if err != nil {
		return "", errors.Wrap(err, "encrypt password")
	}
	if !equal {
		return "", errors.New("passwords not equal")
	}

	token, err := m.jwtManager.CreateToken(user.ID)
	if err != nil {
		return "", errors.Wrap(err, "create token")
	}
	return token, nil
}

func (m *manager) GetUserByToken(ctx context.Context, token string) (models.User, error) {
	userID, err := m.jwtManager.ExtractToken(token)
	if err != nil {
		return models.User{}, errors.Wrap(err, "extract token")
	}

	user, err := m.userRepository.GetUserWithID(ctx, userID)
	if err != nil {
		return models.User{}, errors.Wrap(err, "get user from db")
	}

	return user, nil
}

var emailPattern = regexp.MustCompile(`.+@.+\.`)

func (m *manager) validateEmail(email string) bool {
	return emailPattern.MatchString(email)
}
