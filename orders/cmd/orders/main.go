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
	"orders/internal/api"
	dishesmanager "orders/internal/managers/dishes"
	ordermanager "orders/internal/managers/order"
	"orders/internal/repository/dishes"
	"orders/internal/repository/order"
	"orders/internal/services/auth"
)

type Config struct {
	API      api.Config      `yaml:"api"`
	Postgres postgres.Config `yaml:"postgres"`
	Auth     auth.Config     `yaml:"auth"`
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

	dishRepository := dishes.NewRepository(db)
	dishManager := dishesmanager.NewManager(dishRepository)

	orderRepository := order.NewRepository(db)
	orderManager := ordermanager.NewManager(orderRepository, dishRepository)

	processor := ordermanager.NewProcessor(logger, orderRepository, dishRepository)
	defer processor.Close()
	go processor.Run()

	authService := auth.NewService(config.Auth)

	servant, err := api.NewServant(config.API, logger, authService, dishManager, orderManager)
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
