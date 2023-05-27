package http

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"

	"github.com/4el0ve4ek/restaraunt-api/library/pkg/optional"
)

type JSONHandler[I, O any] interface {
	ServeJSON(request *http.Request, input I) (Response[O], error)
}

func NewJSONWrapper[I, O any](handler JSONHandler[I, O]) *jsonWrapper[I, O] {
	return &jsonWrapper[I, O]{
		h: handler,
	}
}

type jsonWrapper[I, O any] struct {
	h JSONHandler[I, O]
}

func (w *jsonWrapper[I, O]) ServeHTTP(request *http.Request) (Response[[]byte], error) {
	var input I
	switch any(input).(type) {
	case struct{}:
	default:
		bodyDecoder := json.NewDecoder(request.Body)
		err := bodyDecoder.Decode(&input)
		if err != nil && errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			err = errors.Wrap(err, "decode request body as json")
			return Response[[]byte]{
				Content:    optional.New([]byte(err.Error())),
				StatusCode: optional.New(http.StatusBadRequest),
			}, err
		}
	}
	response, err := w.h.ServeJSON(request, input)
	if err != nil {
		return Response[[]byte]{}, errors.Wrap(err, "serve json")
	}

	var ret Response[[]byte]
	if response.StatusCode.IsPresent() {
		ret.StatusCode.Set(response.StatusCode.Get())
	}
	if response.Content.IsPresent() {
		body, err := json.Marshal(response.Content.Get())
		if err != nil {
			return Response[[]byte]{}, errors.Wrap(err, "marshal response")
		}
		ret.Content.Set(body)
	}
	return ret, nil
}
