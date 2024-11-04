package api

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"github.com/sheginabo/go-quick-api/internal/presentation/handlers"
	"github.com/sheginabo/go-quick-api/internal/presentation/middlewares"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
	"net/http"
	"time"
)

type Module struct {
	Mux        *http.ServeMux
	Server     *http.Server
	Stop       context.CancelFunc
	Middleware http.Handler
}

type AllHandlers struct {
	InternalHandler *handlers.InternalHandler
}

func NewModule(stop context.CancelFunc) *Module {
	mux := http.NewServeMux()

	module := &Module{
		Mux:  mux,
		Stop: stop,
	}

	module.SetupRoute(module.NewHandlers(stop))

	return module
}

func (module *Module) NewHandlers(stop context.CancelFunc) AllHandlers {
	return AllHandlers{
		InternalHandler: handlers.NewInternalHandler(),
	}
}

func (module *Module) SetupRoute(allHandlers AllHandlers) {
	// Basic check
	module.Mux.HandleFunc("/health", handlers.HealthCheck)
	module.Mux.HandleFunc("/test/ip", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := `{"X-Forwarded-For": "ip:" + r.Header.Get("X-Forwarded-For") + ", "c_ClientIP": "ip:" + r.RemoteAddr}`
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	})
	// Basic handler
	module.Mux.HandleFunc("/hello", allHandlers.InternalHandler.PostHello)

	// set middleware
	module.Middleware = middlewares.ChainMiddleware(module.Mux, middlewares.Logger, middlewares.CustomRecovery)
}

// Run api module
func (module *Module) Run(ctx context.Context, waitGroup *errgroup.Group) {
	module.Server = &http.Server{
		Addr:    viper.GetString("SERVER_ADDRESS"),
		Handler: module.Middleware,
	}

	waitGroup.Go(func() error {
		log.Info().Msgf("Starting HTTP server on %s\n", viper.GetString("SERVER_ADDRESS"))
		err := module.Server.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			log.Error().Err(err).Msg("HTTP server failed to serve")
			return err
		}
		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("graceful shutdown HTTP(api) server")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		shutdownErr := module.Server.Shutdown(shutdownCtx)
		if shutdownErr != nil {
			log.Error().Err(shutdownErr).Msg("failed to shutdown HTTP(api) server")
			return shutdownErr
		}

		log.Info().Msg("HTTP(api) server is stopped")
		return nil
	})
}
