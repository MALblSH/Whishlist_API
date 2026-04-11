package handlers

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/MALblSH/Wishlist_API/internal/application/service"
	"github.com/MALblSH/Wishlist_API/internal/infrastructure/httppkg/httputils"
	"github.com/MALblSH/Wishlist_API/internal/infrastructure/logging"
)

type PublicHandler struct {
	wishlistService service.WishlistService
	itemService     service.ItemService
}

func NewPublicHandler(wishlistService service.WishlistService, itemService service.ItemService) *PublicHandler {
	return &PublicHandler{
		wishlistService: wishlistService,
		itemService:     itemService,
	}
}

func (p *PublicHandler) GetWishList(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	token := chi.URLParam(r, "token")

	resp, err := p.wishlistService.GetByToken(r.Context(), token)
	if err != nil {
		log.Warn("failed to get wishlist by token",
			slog.String("token", token),
			slog.String("error", err.Error()),
		)
		httputils.MapError(w, err)
		return
	}

	log.Info("wishlist received successfully",
		slog.String("token", token),
		slog.Int("items_count", len(resp.Items)),
	)

	httputils.WriteJSON(w, http.StatusOK, resp)
}

func (p *PublicHandler) ReserveItem(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	token := chi.URLParam(r, "token")
	itemID, err := uuid.Parse(chi.URLParam(r, "itemId"))
	if err != nil {
		log.Warn("invalid item id",
			slog.String("itemId", chi.URLParam(r, "itemId")),
			slog.String("error", err.Error()),
		)
		httputils.WriteJSONError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	err = p.itemService.Reserve(r.Context(), token, itemID)
	if err != nil {
		log.Warn("failed to reserve item",
			slog.String("token", token),
			slog.String("itemId", itemID.String()),
			slog.String("error", err.Error()),
		)
		httputils.MapError(w, err)
		return
	}

	log.Info("item reserved successfully",
		slog.String("token", token),
		slog.String("itemId", itemID.String()),
	)

	httputils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "item reserved successfully",
	})
}
