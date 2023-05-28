package api

import (
	"fmt"
	stdhttp "net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/4el0ve4ek/restaraunt-api/library/pkg/http"
	"github.com/4el0ve4ek/restaraunt-api/library/pkg/log"

	"auth/internal/manager/user"
)

func NewServant(
	cfg Config,
	logger log.Logger,
	userManager user.Manager,
) (*servant, error) {
	server := &stdhttp.Server{
		ReadHeaderTimeout: time.Minute,
	}
	server.Addr = fmt.Sprintf(":%d", cfg.Port)
	router := chi.NewRouter()

	router.Use(middleware.StripSlashes, middleware.Recoverer, middleware.RequestID, middleware.Logger)
	router.Handle("/debug/", middleware.Profiler())

	router.Method("POST", "/register", wrap[inputRegisterUser, outputRegisterUser](NewRegisterUserHandler(userManager), logger))
	router.Method("POST", "/login", wrap[inputLoginUser, outputLoginUser](NewLoginUserHandler(userManager), logger))
	router.Method("GET", "/user", wrap[inputGetUser, outputGetUser](NewGetUserHandler(userManager), logger))

	server.Handler = router
	return &servant{server: server}, nil
}

type servant struct {
	server *stdhttp.Server
}

func (s *servant) GetServer() *stdhttp.Server {
	return s.server
}

func wrap[I, O any](handler http.JSONHandler[I, O], logger log.Logger) stdhttp.Handler {
	return http.NewWrapper(
		http.NewJSONWrapper(handler),
		logger,
	)
}
