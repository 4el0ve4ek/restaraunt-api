package api

import (
	"fmt"
	stdhttp "net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/4el0ve4ek/restaraunt-api/library/pkg/http"
	"github.com/4el0ve4ek/restaraunt-api/library/pkg/log"

	"orders/internal/managers/dishes"
	"orders/internal/managers/order"
	"orders/internal/services/auth"
)

func NewServant(
	cfg Config,
	logger log.Logger,
	authService auth.Service,
	dishManager dishes.Manager,
	orderManager order.Manager,
) (*servant, error) {
	server := &stdhttp.Server{
		ReadHeaderTimeout: time.Minute,
	}
	server.Addr = fmt.Sprintf(":%d", cfg.Port)
	router := chi.NewRouter()

	router.Use(middleware.StripSlashes, middleware.RequestID, middleware.Logger, middleware.Recoverer)
	router.Handle("/debug/", middleware.Profiler())

	core := core{
		logger:      logger,
		authService: authService,
	}

	router.Method("GET", "/dishes", wrap[inputGetMenu, outputGetMenu](NewGetMenuHandler(dishManager), core))
	router.Method("POST", "/dishes", wrap[inputAddDish, outputAddDish](NewAddDishHandler(dishManager), core))
	router.Method("DELETE", "/dishes/{dishID}", wrap[inputDeleteDish, outputDeleteDish](NewDeleteDishHandler(dishManager), core))
	router.Method("PUT", "/dishes/{dishID}", wrap[inputModifyDish, outputModifyDish](NewModifyDishHandler(dishManager), core))
	// TODO getDishByID /dishes/{dishID}

	router.Method("POST", "/orders", wrap[inputMakeOrder, outputMakeOrder](NewMakeOrderHandler(orderManager), core))
	router.Method("GET", "/orders/{orderID}", wrap[inputGetOrder, outputGetOrder](NewGetOrderHandler(orderManager), core))

	server.Handler = router
	return &servant{server: server}, nil
}

type servant struct {
	server *stdhttp.Server
}

func (s *servant) GetServer() *stdhttp.Server {
	return s.server
}

type core struct {
	logger      log.Logger
	authService auth.Service
}

func wrap[I, O any](handler contextHandler[I, O], core core) stdhttp.Handler {
	return http.NewWrapper(
		http.NewJSONWrapper[I, O](
			newWrapper(handler, core.authService),
		),
		core.logger,
	)
}
