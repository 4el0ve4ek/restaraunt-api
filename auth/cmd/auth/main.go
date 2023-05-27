package main

import (
	stdhttp "net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"

	"github.com/4el0ve4ek/restaraunt-api/library/pkg/config"
	"github.com/4el0ve4ek/restaraunt-api/library/pkg/database/postgres"
	"github.com/4el0ve4ek/restaraunt-api/library/pkg/log"

	"auth/internal/api"
	"auth/internal/manager/jwt"
	"auth/internal/manager/password"
	usermanager "auth/internal/manager/user"
	"auth/internal/repository/user"
)

type Config struct {
	JWT      jwt.Config      `yaml:"jwt"`
	API      api.Config      `yaml:"api"`
	Postgres postgres.Config `yaml:"postgres"`
}

func main() {
	logger := log.NewLogger()
	config, err := config.ReadYML[Config]("./config.yaml")
	if err != nil {
		logger.Fatal(err)
	}

	db, err := postgres.NewDBConn(config.Postgres)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	jwtManager := jwt.NewManager(config.JWT)
	passwordManager := password.NewManager()
	userRepository := user.NewRepository(db)
	userManager := usermanager.NewManager(userRepository, passwordManager, jwtManager)
	servant, err := api.NewServant(config.API, logger, userManager)
	if err != nil {
		logger.Fatal(errors.Wrap(err, "create servant"))
	}

	server := servant.GetServer()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, stdhttp.ErrServerClosed) {
			logger.Error(errors.Wrap(err, "http server failure"))
			sigChan <- syscall.SIGINT
		}
	}()

	<-sigChan
	logger.Info(errors.New("shutting down"))
}
