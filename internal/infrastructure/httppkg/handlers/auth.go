package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MALblSH/Wishlist_API/internal/infrastructure/dto"
	"github.com/MALblSH/Wishlist_API/internal/infrastructure/httppkg/httputils"
)

type AuthHandler struct {
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		httputils.WriteJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	fmt.Printf("AuthHandler.Register email: %s, password: %s\n", req.Email, req.Password)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Register called",
		"email":   req.Email,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		httputils.WriteJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	fmt.Printf("AuthHandler.Login email: %s, password: %s\n", req.Email, req.Password)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login called",
		"email":   req.Email,
	})
}
