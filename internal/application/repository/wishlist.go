package repository

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/MALblSH/Wishlist_API/internal/domain"
)

type WishlistRepository interface {
	Create(ctx context.Context, userID uuid.UUID, title string, description string, eventDate time.Time, shareToken string) (domain.Wishlist, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.Wishlist, error)
	GetByShareToken(ctx context.Context, token string) (domain.Wishlist, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Wishlist, error)
	Update(ctx context.Context, wishlist domain.Wishlist) (domain.Wishlist, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
