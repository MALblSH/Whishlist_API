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

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/MALblSH/Wishlist_API/internal/application/service"
	"github.com/MALblSH/Wishlist_API/internal/infrastructure/config"
	"github.com/MALblSH/Wishlist_API/internal/infrastructure/httppkg/handlers"
	"github.com/MALblSH/Wishlist_API/internal/infrastructure/httppkg/router"
	"github.com/MALblSH/Wishlist_API/internal/infrastructure/logging"
	"github.com/MALblSH/Wishlist_API/internal/infrastructure/ormrepository"
	"github.com/MALblSH/Wishlist_API/internal/infrastructure/postgres"
	"github.com/MALblSH/Wishlist_API/internal/infrastructure/tokens/access"
)

const (
	AccessTokenTTL  = time.Hour * 24
	shutdownTimeout = 5 * time.Second
)

func RunWishlist() error {
	log := logging.Init()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	ctx = logging.With(ctx, log)

	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	var pool *pgxpool.Pool
	pool, err = postgres.NewPool(ctx, cfg)
	if err != nil {
		return fmt.Errorf("create pool: %w", err)
	}

	defer pool.Close()

	itemRepo := ormrepository.NewItemRepository(pool)
	wishlistRepo := ormrepository.NewWishlistRepository(pool)
	userRepo := ormrepository.NewUserRepository(pool)

	jwtManager := access.NewManager(cfg.AccessTokenSecret, AccessTokenTTL)

	itemService := service.NewItemService(itemRepo, wishlistRepo)
	wishlistService := service.NewWishlistService(wishlistRepo, itemRepo)
	authService := service.NewAuthService(userRepo, jwtManager)

	itemHandler := handlers.NewItemHandler(itemService)
	wishlistHandler := handlers.NewWishlistHandler(wishlistService)
	authHandler := handlers.NewAuthHandler(authService)
	publicHandler := handlers.NewPublicHandler(wishlistService, itemService)

	rtr := router.RegisterRoutes(authHandler, wishlistHandler, itemHandler, publicHandler, jwtManager)

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
	err = httpServer.Shutdown(shutdownCtx)
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
