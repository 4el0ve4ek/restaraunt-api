package http

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/4el0ve4ek/restaraunt-api/library/pkg/log"
	"github.com/4el0ve4ek/restaraunt-api/library/pkg/optional"
)

func NewWrapper(handler Handler, logger log.Logger) *wrapper {
	return &wrapper{
		h:      handler,
		logger: logger,
	}
}

type wrapper struct {
	h      Handler
	logger log.Logger
}

func (w *wrapper) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	response, err := w.h.ServeHTTP(request)
	if err != nil {
		w.logger.Error(errors.Wrap(err, "serve http"))
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(response.StatusCode.GetOrDefault(http.StatusOK))
	_, err = writer.Write(response.Content.Get())
	w.logger.Error(errors.Wrap(err, "write response"))
}

type Handler interface {
	ServeHTTP(request *http.Request) (Response[[]byte], error)
}

type Response[T any] struct {
	StatusCode optional.Optional[int]
	Content    optional.Optional[T]
}
