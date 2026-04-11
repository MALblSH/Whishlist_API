package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/MALblSH/Wishlist_API/internal/application/service"
	"github.com/MALblSH/Wishlist_API/internal/infrastructure/dto"
	"github.com/MALblSH/Wishlist_API/internal/infrastructure/httppkg/httputils"
	"github.com/MALblSH/Wishlist_API/internal/infrastructure/logging"
)

type WishlistHandler struct {
	svc service.WishlistService
}

func NewWishlistHandler(svc service.WishlistService) *WishlistHandler {
	return &WishlistHandler{
		svc: svc,
	}
}

func (h *WishlistHandler) List(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	userID, ok := httputils.GetUserIDFromContext(r)
	if !ok {
		log.Warn("unauthorized access attempt to list wishlists")
		httputils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	resp, err := h.svc.List(r.Context(), userID)
	if err != nil {
		log.Warn("failed to get list wishlists",
			slog.String("user_id", userID.String()),
			slog.String("error", err.Error()),
		)
		httputils.MapError(w, err)
		return
	}

	log.Info("wishlists received successfully",
		slog.String("user_id", userID.String()),
		slog.Int("wishlist_count", len(resp.WishLists)),
	)

	httputils.WriteJSON(w, http.StatusOK, resp)
}

func (h *WishlistHandler) Create(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	userID, ok := httputils.GetUserIDFromContext(r)
	if !ok {
		log.Warn("unauthorized access attempt to create wishlists")
		httputils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req dto.WishListRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Warn("failed to decode wishlist create request body",
			slog.String("error", err.Error()),
		)
		httputils.WriteJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	var resp dto.WishListResponse
	resp, err = h.svc.Create(r.Context(), userID, req)
	if err != nil {
		log.Warn("failed to create wishlist",
			slog.String("user_id", userID.String()),
			slog.String("error", err.Error()),
		)
		httputils.MapError(w, err)
		return
	}

	log.Info("wishlist created successfully",
		slog.String("user_id", userID.String()),
		slog.String("wishlist_id", resp.ID.String()),
	)

	httputils.WriteJSON(w, http.StatusOK, resp)
}

func (h *WishlistHandler) Delete(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	userID, ok := httputils.GetUserIDFromContext(r)
	if !ok {
		log.Warn("unauthorized access attempt to delete wishlist")
		httputils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var wishlistID uuid.UUID
	wishlistID, ok = httputils.ParseUUIDParam(w, r, "id")
	if !ok {
		log.Warn("invalid wishlist id in delete request")
		return
	}
	err := h.svc.Delete(r.Context(), userID, wishlistID)
	if err != nil {
		log.Warn("failed to delete wishlist",
			slog.String("user_id", userID.String()),
			slog.String("wishlist_id", wishlistID.String()),
			slog.String("error", err.Error()),
		)
		httputils.MapError(w, err)
		return
	}

	log.Info("wishlist deleted successfully")

	httputils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "wishlist deleted successfully",
	})
}

func (h *WishlistHandler) Update(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	userID, ok := httputils.GetUserIDFromContext(r)
	if !ok {
		log.Warn("unauthorized access attempt to update wishlist")
		httputils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var wishlistID uuid.UUID
	wishlistID, ok = httputils.ParseUUIDParam(w, r, "id")
	if !ok {
		log.Warn("invalid wishlist id in update request")
		return
	}

	var req dto.WishListUpdateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Warn("failed to decode wishlist update request body",
			slog.String("error", err.Error()),
		)
		httputils.WriteJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	var resp dto.WishListResponse
	resp, err = h.svc.Update(r.Context(), userID, wishlistID, req)
	if err != nil {
		log.Warn("failed to update wishlist",
			slog.String("user_id", userID.String()),
			slog.String("wishlist_id", wishlistID.String()),
			slog.String("error", err.Error()),
		)
		httputils.MapError(w, err)
		return
	}

	log.Info("wishlist updated successfully",
		slog.String("user_id", userID.String()),
		slog.String("wishlist_id", wishlistID.String()),
	)

	httputils.WriteJSON(w, http.StatusOK, resp)
}

func (h *WishlistHandler) Get(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	userID, ok := httputils.GetUserIDFromContext(r)
	if !ok {
		log.Warn("unauthorized access attempt to get wishlist")
		httputils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	wishlistID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		log.Warn("invalid wishlist id in get request",
			slog.String("error", err.Error()),
		)
		httputils.WriteJSONError(w, http.StatusBadRequest, "invalid wishlist id")
		return
	}

	var resp dto.WishListResponse
	resp, err = h.svc.Get(r.Context(), userID, wishlistID)
	if err != nil {
		log.Warn("failed to get wishlist",
			slog.String("user_id", userID.String()),
			slog.String("wishlist_id", wishlistID.String()),
			slog.String("error", err.Error()),
		)
		httputils.MapError(w, err)
		return
	}

	log.Info("wishlist returned successfully",
		slog.String("user_id", userID.String()),
		slog.String("wishlist_id", wishlistID.String()),
	)

	httputils.WriteJSON(w, http.StatusOK, resp)
}
