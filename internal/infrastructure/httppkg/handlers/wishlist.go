package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/MALblSH/Wishlist_API/internal/infrastructure/dto"
	"github.com/MALblSH/Wishlist_API/internal/infrastructure/httppkg/httputils"
)

type WishlistHandler struct {
}

func NewWishlistHandler() *WishlistHandler {
	return &WishlistHandler{}
}

func (h *WishlistHandler) List(w http.ResponseWriter, r *http.Request) {
	claims, ok := httputils.GetUserIDFromContext(r)
	if !ok {
		httputils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	fmt.Printf("WishlistHandler.List claims: %s\n", claims)
}

func (h *WishlistHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims, ok := httputils.GetUserIDFromContext(r)
	if !ok {
		httputils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req dto.WishListRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		httputils.WriteJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	fmt.Printf("WishlistHandler.Create claims: %s, title: %s, description: %s, event_date: %s\n",
		claims, req.Title, req.Description, req.EventDate)

	json.NewEncoder(w).Encode(map[string]string{
		"message":    "Create wishlist called",
		"claims":     claims,
		"title":      req.Title,
		"event_date": req.EventDate,
	})
}

func (h *WishlistHandler) Delete(w http.ResponseWriter, r *http.Request) {
	claims, ok := httputils.GetUserIDFromContext(r)
	if !ok {
		httputils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	id := chi.URLParam(r, "id")

	fmt.Printf("WishlistHandler.Delete claims: %s, wishlist_id: %s\n", claims, id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message":     "Delete wishlist called",
		"claims":      claims,
		"wishlist_id": id,
	})
}

func (h *WishlistHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims, ok := httputils.GetUserIDFromContext(r)
	if !ok {
		httputils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	id := chi.URLParam(r, "id")
	var req dto.WishListUpdateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		httputils.WriteJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	fmt.Printf("WishlistHandler.Update claims: %s, wishlist_id: %s, title: %s, description: %s, event_date: %s\n",
		claims, id, req.Title, req.Description, req.EventDate)

	json.NewEncoder(w).Encode(map[string]string{
		"message":     "Update wishlist called",
		"claims":      claims,
		"wishlist_id": id,
	})
}

func (h *WishlistHandler) Get(w http.ResponseWriter, r *http.Request) {
	claims, ok := httputils.GetUserIDFromContext(r)
	if !ok {
		httputils.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	id := chi.URLParam(r, "id")

	fmt.Printf("WishlistHandler.Get claims: %s, wishlist_id: %s\n", claims, id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message":     "Get wishlist called",
		"claims":      claims,
		"wishlist_id": id,
	})
}
