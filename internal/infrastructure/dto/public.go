package dto

import (
	"time"

	"github.com/google/uuid"
)

type PublicWishListResponse struct {
	ID          uuid.UUID              `json:"id"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	EventDate   time.Time              `json:"event_date"`
	Items       []WishListItemResponse `json:"items"`
}
