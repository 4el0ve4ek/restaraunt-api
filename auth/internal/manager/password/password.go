package password

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func NewManager() *manager {
	return &manager{}
}

type manager struct{}

func (m *manager) Encrypt(password string) (string, error) {
	encryptedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.Wrap(err, "generate password")
	}
	return string(encryptedPass), nil
}

func (m *manager) CompareWithHashed(hashedPassword, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	if err == nil {
		return true, nil
	}
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, nil
	}
	return false, err
}
