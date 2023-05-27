package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"orders/internal/models"
)

const (
	AuthorizationHeader = "Authorization"
)

type Service interface {
	Get(ctx context.Context, authorizationHeaderValue string) (models.User, error)
}

func NewService(cfg Config) *service {
	return &service{
		cfg:        cfg,
		httpclient: &http.Client{},
	}
}

type service struct {
	cfg        Config
	httpclient *http.Client
}

func (s *service) Get(ctx context.Context, authorizationHeaderValue string) (models.User, error) {
	authRequest, err := http.NewRequest(http.MethodGet, s.getURL("/user"), nil)
	if err != nil {
		return models.User{}, errors.Wrap(err, "create request")
	}

	authRequest.Header.Set(AuthorizationHeader, authorizationHeaderValue)
	authRequest = authRequest.WithContext(ctx)

	response, err := s.httpclient.Do(authRequest)
	if err != nil {
		return models.User{}, errors.Wrap(err, "do request")
	}
	defer response.Body.Close()

	var user models.User
	bodyDecoder := json.NewDecoder(response.Body)
	err = bodyDecoder.Decode(&user)
	if err != nil {
		return models.User{}, errors.Wrap(err, "decode response as user")
	}

	return user, nil
}

func (s *service) getURL(path string) string {
	return fmt.Sprintf("http://%s:%s") + path
}
