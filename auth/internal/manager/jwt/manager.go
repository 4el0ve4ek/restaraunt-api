package jwt

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/4el0ve4ek/restaraunt-api/library/pkg/optional"
)

const (
	Header = "Authorization"

	accessTokenTTL = time.Hour * 24 * 7 // 7 days
)

func NewManager(cfg Config, sessionRepository sessionRepository) *manager {
	return &manager{
		cfg:               cfg,
		sessionRepository: sessionRepository,
	}
}

type manager struct {
	cfg               Config
	sessionRepository sessionRepository
}

func (m *manager) CreateToken(ctx context.Context, userID int) (string, error) {
	sessionToken := m.generateToken(userID)
	expiresAt := time.Now().Add(accessTokenTTL)
	err := m.sessionRepository.AddSession(ctx, userID, sessionToken, expiresAt)
	if err != nil {
		return "", errors.Wrap(err, "add session to db")
	}

	return sessionToken, nil
}

func (m *manager) ExtractToken(ctx context.Context, jwtToken string) (optional.Optional[int], error) {
	userID, err := m.sessionRepository.GetSession(ctx, jwtToken)
	return userID, errors.Wrap(err, "get session from db")
}

func (m *manager) generateToken(salt int) string {
	alph := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	alphy := []rune(alph)

	salty := []rune(strconv.Itoa(salt) + " " + m.cfg.SecretKey + " ")
	hashedValue := make([]rune, len(salty), len(salty)+50)
	copy(salty, hashedValue)

	for ls := len(salty); ls < cap(salty); ls++ {
		salty = append(salty, alphy[rand.Intn(len(alphy))])
	}

	hashed := sha256.Sum256([]byte(string(salty)))
	return fmt.Sprintf("%x", hashed[:])
}
