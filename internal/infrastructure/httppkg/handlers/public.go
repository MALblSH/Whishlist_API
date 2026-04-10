package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type PublicHandler struct {
}

func NewPublicHandler() *PublicHandler {
	return &PublicHandler{}
}

func (p *PublicHandler) GetWishList(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	fmt.Printf("PublicHandler.GetWishlist token: %s\n", token)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Get public wishlist called",
		"token":   token,
	})
}

func (p *PublicHandler) ReserveItem(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	itemID := chi.URLParam(r, "itemId")

	fmt.Printf("PublicHandler.ReserveItem token: %s, itemID: %s\n", token, itemID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Reserve item called",
		"token":   token,
		"itemId":  itemID,
	})
}
