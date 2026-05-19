package core_http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	core_logger "github.com/qwinkki/http-rest-taskmgr/internal/core/logger"
	"go.uber.org/zap"
)

type Server struct {
	mux    *http.ServeMux
	config Config
	log    *core_logger.Logger
}

func NewServer(config Config, log *core_logger.Logger) *Server {
	return &Server{
		mux:    http.NewServeMux(),
		config: config,
		log:    log,
	}
}

func (s *Server) RegisterRoutes(routers ...*APIVersionRouter) {
	for _, router := range routers {
		s.mux.Handle("/api/"+string(router.apiVersion)+"/", router.ServeMux)
	}
}

func (s *Server) Run(ctx context.Context) error {
	server := &http.Server{
		Addr:    s.config.Address,
		Handler: s.mux,
	}

	ch := make(chan error, 1)

	go func() {
		defer close(ch)
		s.log.Warn("start http server", zap.String("address", s.config.Address))

		err := server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		return fmt.Errorf("listen and server http: %w", err)
	case <-ctx.Done():
		s.log.Warn("shutting down http server")

		shutdownctx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownctx); err != nil {
			_ = server.Close()
			return fmt.Errorf("shutdown http server: %w", err)
		}
		s.log.Warn("http server stopped ")
		return nil
	}

}
