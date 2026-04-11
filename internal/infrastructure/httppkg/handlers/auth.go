package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/MALblSH/Wishlist_API/internal/application/service"
	"github.com/MALblSH/Wishlist_API/internal/infrastructure/dto"
	"github.com/MALblSH/Wishlist_API/internal/infrastructure/httppkg/httputils"
	"github.com/MALblSH/Wishlist_API/internal/infrastructure/logging"
)

type AuthHandler struct {
	svc service.AuthService
}

func NewAuthHandler(svc service.AuthService) *AuthHandler {
	return &AuthHandler{
		svc: svc,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	var req dto.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Warn("failed to decode register request body",
			slog.String("error", err.Error()),
		)
		httputils.WriteJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	err = h.svc.Register(r.Context(), req)
	if err != nil {
		log.Warn("failed to register user",
			slog.String("email", req.Email),
			slog.String("error", err.Error()),
		)
		httputils.MapError(w, err)
		return
	}

	log.Info("User registered successfully",
		slog.String("email", req.Email),
	)

	httputils.WriteJSON(w, http.StatusCreated, map[string]string{
		"message": "user registered successfully",
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	var req dto.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Warn("failed to decode login request body",
			slog.String("error", err.Error()),
		)
		httputils.WriteJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	var resp dto.LoginResponse
	resp, err = h.svc.Login(r.Context(), req)
	if err != nil {
		log.Warn("failed to login user",
			slog.String("email", req.Email),
			slog.String("error", err.Error()),
		)
		httputils.MapError(w, err)
		return
	}

	log.Info("User logged successfully",
		slog.String("email", req.Email),
	)

	httputils.WriteJSON(w, http.StatusOK, resp)
}
