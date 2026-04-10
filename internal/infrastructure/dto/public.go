package dto

type PublicWishListResponse struct {
	ID          string                 `json:"id"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	EventDate   string                 `json:"event_date"`
	Items       []WishListItemResponse `json:"items"`
}
