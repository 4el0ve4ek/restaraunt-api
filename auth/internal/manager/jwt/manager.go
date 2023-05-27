package jwt

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

const (
	Header = "Authorization"

	accessTokenTTL = time.Hour * 24 * 7 // 7 days
)

func NewManager(cfg Config) *manager {
	return &manager{
		cfg: cfg,
	}
}

type manager struct {
	cfg Config
}

func (m *manager) CreateToken(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ID:        strconv.Itoa(userID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenTTL)),
	})

	jwtSignedValue, err := token.SignedString([]byte(m.cfg.SecretKey))
	if err != nil {
		return "", errors.Wrap(err, "sign token")
	}
	return jwtSignedValue, nil
}

func (m *manager) ExtractToken(jwtToken string) (int, error) {
	token, err := jwt.ParseWithClaims(jwtToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenUnverifiable
		}
		return []byte(m.cfg.SecretKey), nil
	})
	if err != nil {
		return 0, errors.New("Invalid access token")
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return 0, errors.New("Invalid access token")
	}

	id, err := strconv.Atoi(claims.ID)
	if err != nil {
		return 0, err
	}
	return id, nil
}
