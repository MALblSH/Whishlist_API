package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/MALblSH/Wishlist_API/internal/application/service"
	"github.com/MALblSH/Wishlist_API/internal/infrastructure/dto"
	"github.com/MALblSH/Wishlist_API/internal/infrastructure/httppkg/httputils"
	"github.com/MALblSH/Wishlist_API/internal/infrastructure/logging"
)

type ItemHandler struct {
	svc service.ItemService
}

func NewItemHandler(svc service.ItemService) *ItemHandler {
	return &ItemHandler{
		svc: svc,
	}
}

func (h *ItemHandler) List(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	userID, ok := httputils.GetUserIDFromContext(r)
	if !ok {
		log.Warn("unauthorized access attempt to list items")
		httputils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var wishlistID uuid.UUID
	wishlistID, ok = httputils.ParseUUIDParam(w, r, "id")
	if !ok {
		log.Warn("invalid wishlist id",
			slog.String("user_id", userID.String()),
		)
		return
	}

	resp, err := h.svc.List(r.Context(), userID, wishlistID)
	if err != nil {
		log.Warn("failed to get list items",
			slog.String("user_id", userID.String()),
			slog.String("wishlist_id", wishlistID.String()),
			slog.Any("error", err),
		)
		httputils.MapError(w, err)
		return
	}

	log.Info("item list received successfully",
		slog.String("user_id", userID.String()),
		slog.String("wishlist_id", wishlistID.String()),
		slog.Int("item_count", len(resp.Items)),
	)

	httputils.WriteJSON(w, http.StatusOK, resp)
}

func (h *ItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	userID, ok := httputils.GetUserIDFromContext(r)
	if !ok {
		log.Warn("unauthorized access attempt to create item")
		httputils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var wishlistID uuid.UUID
	wishlistID, ok = httputils.ParseUUIDParam(w, r, "id")
	if !ok {
		log.Warn("invalid wishlist id",
			slog.String("user_id", userID.String()),
		)
		return
	}

	var req dto.WishListItemRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Warn("failed to decode item create request body",
			slog.String("error", err.Error()),
		)
		httputils.WriteJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	var resp dto.WishListItemResponse
	resp, err = h.svc.Create(r.Context(), userID, wishlistID, req)
	if err != nil {
		log.Warn("failed to create item",
			slog.String("user_id", userID.String()),
			slog.String("wishlist_id", wishlistID.String()),
			slog.Any("error", err),
		)
		httputils.MapError(w, err)
		return
	}

	log.Info("item created successfully",
		slog.String("user_id", userID.String()),
		slog.String("wishlist_id", wishlistID.String()),
		slog.String("item_id", resp.ID.String()),
	)

	httputils.WriteJSON(w, http.StatusCreated, resp)
}

func (h *ItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	userID, ok := httputils.GetUserIDFromContext(r)
	if !ok {
		log.Warn("unauthorized access attempt to delete item")
		httputils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var wishlistID uuid.UUID
	wishlistID, ok = httputils.ParseUUIDParam(w, r, "id")
	if !ok {
		log.Warn("invalid wishlist id",
			slog.String("user_id", userID.String()),
		)
		return
	}

	var itemID uuid.UUID
	itemID, ok = httputils.ParseUUIDParam(w, r, "itemId")
	if !ok {
		log.Warn("invalid item id",
			slog.String("user_id", userID.String()),
			slog.String("wishlist_id", wishlistID.String()),
		)
		return
	}

	err := h.svc.Delete(r.Context(), userID, wishlistID, itemID)
	if err != nil {
		log.Warn("failed to delete item",
			slog.String("user_id", userID.String()),
			slog.String("wishlist_id", wishlistID.String()),
			slog.String("item_id", itemID.String()),
			slog.Any("error", err),
		)
		httputils.MapError(w, err)
		return
	}

	log.Info("item deleted successfully",
		slog.String("user_id", userID.String()),
		slog.String("wishlist_id", wishlistID.String()),
		slog.String("item_id", itemID.String()),
	)

	httputils.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "item successfully deleted",
	})
}

func (h *ItemHandler) Update(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	userID, ok := httputils.GetUserIDFromContext(r)
	if !ok {
		log.Warn("unauthorized access attempt to update item")
		httputils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var wishlistID uuid.UUID
	wishlistID, ok = httputils.ParseUUIDParam(w, r, "id")
	if !ok {
		log.Warn("invalid wishlist id",
			slog.String("user_id", userID.String()),
		)
		return
	}

	var itemID uuid.UUID
	itemID, ok = httputils.ParseUUIDParam(w, r, "itemId")
	if !ok {
		log.Warn("invalid item id",
			slog.String("user_id", userID.String()),
			slog.String("wishlist_id", wishlistID.String()),
		)
		return
	}

	var req dto.WishListItemUpdateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Warn("failed to decode item update request body",
			slog.String("error", err.Error()),
		)
		httputils.WriteJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	var resp dto.WishListItemResponse
	resp, err = h.svc.Update(r.Context(), userID, wishlistID, itemID, req)
	if err != nil {
		log.Warn("failed to update item",
			slog.String("user_id", userID.String()),
			slog.String("wishlist_id", wishlistID.String()),
			slog.String("item_id", itemID.String()),
			slog.String("error", err.Error()),
		)
		httputils.MapError(w, err)
		return
	}

	log.Info("item updated successfully",
		slog.String("user_id", userID.String()),
		slog.String("wishlist_id", wishlistID.String()),
		slog.String("item_id", itemID.String()),
	)

	httputils.WriteJSON(w, http.StatusOK, resp)
}
