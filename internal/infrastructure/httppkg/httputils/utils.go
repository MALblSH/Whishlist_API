package httputils

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/MALblSH/Wishlist_API/internal/domain"
	"github.com/MALblSH/Wishlist_API/internal/infrastructure/httppkg/middleware"
)

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func WriteJSONError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, map[string]string{"error": message})
}

func GetUserIDFromContext(r *http.Request) (uuid.UUID, bool) {
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		return uuid.Nil, false
	}
	return claims.UserID, true
}

func ParseUUIDParam(w http.ResponseWriter, r *http.Request, param string) (uuid.UUID, bool) {
	id, err := uuid.Parse(chi.URLParam(r, param))
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, "invalid "+param)
		return uuid.Nil, false
	}
	return id, true
}

func MapError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		WriteJSONError(w, http.StatusBadRequest, "invalid input")
	case errors.Is(err, domain.ErrInvalidCredentials):
		WriteJSONError(w, http.StatusUnauthorized, "invalid credentials")
	case errors.Is(err, domain.ErrNotFound):
		WriteJSONError(w, http.StatusNotFound, "resource not found")
	case errors.Is(err, domain.ErrAlreadyReserved):
		WriteJSONError(w, http.StatusConflict, "item already reserved")
	case errors.Is(err, domain.ErrForbidden):
		WriteJSONError(w, http.StatusForbidden, "does not have access to this resource")
	case errors.Is(err, domain.ErrAlreadyExists):
		WriteJSONError(w, http.StatusConflict, "resource already exists")

	default:
		WriteJSONError(w, http.StatusInternalServerError, "internal server error")
	}
}
