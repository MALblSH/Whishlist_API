package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/MALblSH/Wishlist_API/internal/infrastructure/httppkg/handlers"
	authmiddleware "github.com/MALblSH/Wishlist_API/internal/infrastructure/httppkg/middleware"
)

func RegisterRoutes() *chi.Mux {
	authHandler := handlers.NewAuthHandler()
	wishlistHandler := handlers.NewWishlistHandler()
	itemHandler := handlers.NewItemHandler()
	publicHandler := handlers.NewPublicHandler()
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})

	r.Route("/wishlist", func(r chi.Router) {
		r.Use(authmiddleware.AuthMiddleWare)
		r.Get("/", wishlistHandler.List)
		r.Post("/", wishlistHandler.Create)
		r.Put("/{id}", wishlistHandler.Update)
		r.Delete("/{id}", wishlistHandler.Delete)
		r.Get("/{id}", wishlistHandler.Get)

		r.Route("/{id}/items", func(r chi.Router) {
			r.Get("/", itemHandler.List)
			r.Post("/", itemHandler.Create)
			r.Put("/{itemId}", itemHandler.Update)
			r.Delete("/{itemId}", itemHandler.Delete)
		})
	})

	r.Route("/public", func(r chi.Router) {
		r.Get("/{token}", publicHandler.GetWishList)
		r.Post("/{token}/items/{itemId}/reserve", publicHandler.ReserveItem)
	})
	return r
}
