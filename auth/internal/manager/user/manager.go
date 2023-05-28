package user

import (
	"context"
	"regexp"

	"github.com/pkg/errors"

	"auth/internal/models"
	"github.com/4el0ve4ek/restaraunt-api/library/pkg/optional"
)

type Manager interface {
	RegisterUser(ctx context.Context, username, email, userPassword string) (struct {
		FieldsCollide bool
		InvalidEmail  bool
	}, error)
	LoginUser(ctx context.Context, email, userPassword string) (string, error)
	GetUserByToken(ctx context.Context, token string) (optional.Optional[models.User], error)
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

func (m *manager) RegisterUser(ctx context.Context, username, email, userPassword string) (struct {
	FieldsCollide bool
	InvalidEmail  bool
}, error) {
	var ret struct {
		FieldsCollide bool
		InvalidEmail  bool
	}

	passwordEncrypted, err := m.passwordManager.Encrypt(userPassword)
	if err != nil {
		return ret, errors.Wrap(err, "encrypt password")
	}

	if !m.validateEmail(email) {
		ret.InvalidEmail = true
		return ret, nil
	}

	namesCollide, err := m.userRepository.AddNewUser(ctx, username, email, passwordEncrypted)
	if err != nil {
		return ret, errors.Wrap(err, "add new user")
	}

	if namesCollide {
		ret.FieldsCollide = true
	}
	return ret, nil
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

	token, err := m.jwtManager.CreateToken(ctx, user.ID)
	if err != nil {
		return "", errors.Wrap(err, "create token")
	}
	return token, nil
}

func (m *manager) GetUserByToken(ctx context.Context, token string) (optional.Optional[models.User], error) {
	ret := optional.NewEmpty[models.User]()

	userID, err := m.jwtManager.ExtractToken(ctx, token)
	if err != nil {
		return ret, errors.Wrap(err, "extract token")
	}

	if !userID.IsPresent() {
		return ret, nil
	}

	user, err := m.userRepository.GetUserWithID(ctx, userID.Get())
	if err != nil {
		return ret, errors.Wrap(err, "get user from db")
	}
	ret.Set(user)
	return ret, nil
}

var emailPattern = regexp.MustCompile(`.+@.+\.`)

func (m *manager) validateEmail(email string) bool {
	return emailPattern.MatchString(email)
}
