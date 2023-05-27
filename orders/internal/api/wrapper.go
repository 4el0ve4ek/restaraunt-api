package api

import (
	"github.com/pkg/errors"

	"github.com/4el0ve4ek/restaraunt-api/library/pkg/http"
	"orders/internal/context"
	"orders/internal/services/auth"
)

type contextHandler[I, O any] interface {
	ServeJSON(context context.Context, request *http.Request, input I) (http.Response[O], error)
}

func newWrapper[I, O any](h contextHandler[I, O], service auth.Service) *wrapper[I, O] {
	return &wrapper[I, O]{
		h:           h,
		authService: service,
	}
}

type wrapper[I, O any] struct {
	h           contextHandler[I, O]
	authService auth.Service
}

func (w *wrapper[I, O]) ServeJSON(request *http.Request, input I) (http.Response[O], error) {
	ctx, err := context.NewContext(request, w.authService)
	if err != nil {
		return http.Response[O]{}, errors.Wrap(err, "create ctx")
	}

	return w.h.ServeJSON(ctx, request, input)
}
