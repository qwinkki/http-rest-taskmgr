package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/qwinkki/http-rest-taskmgr/internal/core/logger"
	core_http_server "github.com/qwinkki/http-rest-taskmgr/internal/core/transport/http/server"
	user_transport_http "github.com/qwinkki/http-rest-taskmgr/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	core_logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		panic(fmt.Errorf("failed to initialize logger: %w", err))
		os.Exit(1)
	}
	defer core_logger.Close()

	core_logger.Debug("starting application")

	userstransportHTTP := user_transport_http.NewUserHTTPHandler(nil)
	usersRoutes := userstransportHTTP.Routers()

	apiVersion := core_http_server.NewAPIVersionRouter(core_http_server.APIVersionV1)
	apiVersion.RegisterRoutes(usersRoutes...)

	httpServer := core_http_server.NewServer(core_http_server.NewConfigMust(), core_logger)
	httpServer.RegisterRoutes(apiVersion)

	if err := httpServer.Run(ctx); err != nil {
		core_logger.Error("failed to run http server", zap.Error(err))
		os.Exit(1)
	}
}
