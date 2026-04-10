package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/MALblSH/Wishlist_API/internal/infrastructure/dto"
	"github.com/MALblSH/Wishlist_API/internal/infrastructure/httppkg/httputils"
)

type ItemHandler struct {
}

func NewItemHandler() *ItemHandler {
	return &ItemHandler{}
}

func (h *ItemHandler) List(w http.ResponseWriter, r *http.Request) {
	claims, ok := httputils.GetUserIDFromContext(r)
	if !ok {
		httputils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	wishListID := chi.URLParam(r, "id")

	fmt.Printf("ItemHandler.List claims: %s, wishlist_id: %s\n", claims, wishListID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message":     "List items called",
		"claims":      claims,
		"wishlist_id": wishListID,
	})
}

func (h *ItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims, ok := httputils.GetUserIDFromContext(r)
	if !ok {
		httputils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req dto.WishListItemRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		httputils.WriteJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	fmt.Printf("ItemHandler.Create claims: %s, request: %+v\n", claims, req)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Create item called",
		"claims":  claims,
		"request": req,
	})
}

func (h *ItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	claims, ok := httputils.GetUserIDFromContext(r)
	if !ok {
		httputils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	itemID := chi.URLParam(r, "itemId")
	wishListID := chi.URLParam(r, "id")

	fmt.Printf("ItemHandler.Delete claims: %s, item_id: %s, wishlist_id: %s\n", claims, itemID, wishListID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Delete item called",
		"claims":  claims,
		"item_id": itemID,
	})
}

func (h *ItemHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims, ok := httputils.GetUserIDFromContext(r)
	if !ok {
		httputils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	itemID := chi.URLParam(r, "itemId")
	wishListID := chi.URLParam(r, "id")

	var req dto.WishListItemUpdateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		httputils.WriteJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	fmt.Printf("ItemHandler.Update claims: %s, item_id: %s, wishlist_id: %s\n", claims, itemID, wishListID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Update item called",
		"claims":  claims,
		"item_id": itemID,
		"request": req,
	})
}
