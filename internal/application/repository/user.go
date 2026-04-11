package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/MALblSH/Wishlist_API/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, email string, passwordHash string) error
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.User, error)
}
