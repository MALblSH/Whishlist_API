package domain

import (
	"time"

	"github.com/google/uuid"
)

type WishlistItem struct {
	ID          uuid.UUID
	WishlistID  uuid.UUID
	Title       string
	Description *string
	ProductURL  *string
	Priority    int
	IsReserved  bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
