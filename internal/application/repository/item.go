package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/MALblSH/Wishlist_API/internal/domain"
)

type ItemRepository interface {
	Create(ctx context.Context, wishlistID uuid.UUID, title string, description *string, productURL *string, priority int) (domain.WishlistItem, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.WishlistItem, error)
	ListByWishlistID(ctx context.Context, wishlistID uuid.UUID) ([]domain.WishlistItem, error)
	Update(ctx context.Context, item domain.WishlistItem) (domain.WishlistItem, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Reserve(ctx context.Context, id uuid.UUID) error
}
