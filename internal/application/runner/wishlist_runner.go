package runner

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/MALblSH/Wishlist_API/internal/infrastructure/httppkg/router"
	"github.com/MALblSH/Wishlist_API/internal/infrastructure/logging"
)

const (
	shutdownTimeout = 5 * time.Second
)

func RunWishlist() error {
	log := logging.Init()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	ctx = logging.With(ctx, log)

	rtr := router.RegisterRoutes()

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: rtr,
	}

	go func() {
		log.Info("HTTP server starting", "addr", httpServer.Addr)
		srvErr := httpServer.ListenAndServe()
		if srvErr != nil && !errors.Is(srvErr, http.ErrServerClosed) {
			log.Error(
				"HTTP server failed",
				slog.String("error", srvErr.Error()),
			)
		}
	}()

	<-ctx.Done()
	log.Info("shutting down HTTP server",
		slog.String("addr", httpServer.Addr),
	)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	err := httpServer.Shutdown(shutdownCtx)
	if err != nil {
		log.Error(
			"failed to shutdown HTTP server",
			slog.String("error", err.Error()),
			slog.String("addr", httpServer.Addr),
		)
		return fmt.Errorf("failed to shutdown HTTP server: %w", err)
	}
	log.Info("HTTP server stopped")

	return nil
}
