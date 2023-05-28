package jwt

import (
	"context"
	"time"

	"github.com/4el0ve4ek/restaraunt-api/library/pkg/optional"
)

type sessionRepository interface {
	AddSession(ctx context.Context, userID int, sessionToken string, createdAt time.Time) error
	GetSession(ctx context.Context, sessionToken string) (optional.Optional[int], error)
}
