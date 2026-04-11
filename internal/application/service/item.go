package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"github.com/MALblSH/Wishlist_API/internal/application/repository"
	"github.com/MALblSH/Wishlist_API/internal/domain"
	"github.com/MALblSH/Wishlist_API/internal/infrastructure/dto"
	"github.com/MALblSH/Wishlist_API/internal/infrastructure/logging"
)

type ItemService interface {
	Create(ctx context.Context, userID uuid.UUID, wishlistID uuid.UUID, req dto.WishListItemRequest) (dto.WishListItemResponse, error)
	List(ctx context.Context, userID uuid.UUID, wishlistID uuid.UUID) (dto.WishListItemListResponse, error)
	Update(ctx context.Context, userID uuid.UUID, wishlistID uuid.UUID, itemID uuid.UUID, req dto.WishListItemUpdateRequest) (dto.WishListItemResponse, error)
	Delete(ctx context.Context, userID uuid.UUID, wishlistID uuid.UUID, itemID uuid.UUID) error

	Reserve(ctx context.Context, token string, itemID uuid.UUID) error
}

type itemService struct {
	itemRepo     repository.ItemRepository
	wishlistRepo repository.WishlistRepository
}

func NewItemService(itemRepo repository.ItemRepository, wishlistRepo repository.WishlistRepository) ItemService {
	return &itemService{
		itemRepo:     itemRepo,
		wishlistRepo: wishlistRepo,
	}
}

func (s *itemService) Create(ctx context.Context, userID uuid.UUID, wishlistID uuid.UUID, req dto.WishListItemRequest) (dto.WishListItemResponse, error) {
	log := logging.FromContext(ctx)
	if req.Title == "" {
		return dto.WishListItemResponse{}, domain.ErrInvalidInput
	}

	wishlist, err := s.wishlistRepo.GetByID(ctx, wishlistID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			log.Warn("Wishlist not found",
				slog.String("wishlist_id", wishlistID.String()),
			)
			return dto.WishListItemResponse{}, domain.ErrInvalidCredentials
		}
		return dto.WishListItemResponse{}, fmt.Errorf("get wishlist: %w", err)
	}

	if wishlist.UserID != userID {
		log.Warn("Wishlist user ID mismatch",
			slog.String("wishlist_id", wishlistID.String()),
			slog.String("expected_user_id", wishlist.UserID.String()),
			slog.String("actual_user_id", userID.String()),
		)
		return dto.WishListItemResponse{}, domain.ErrForbidden
	}

	var item domain.WishlistItem
	item, err = s.itemRepo.Create(ctx, wishlistID, req.Title, req.Description, req.ProductURL, req.Priority)
	if err != nil {
		return dto.WishListItemResponse{}, fmt.Errorf("create item: %w", err)
	}

	log.Info("Item created",
		slog.String("item_id", item.ID.String()),
		slog.String("wishlist_id", wishlistID.String()),
		slog.String("title", item.Title),
	)

	return dto.WishListItemResponse{
		ID:          item.ID,
		Title:       item.Title,
		Description: item.Description,
		ProductURL:  item.ProductURL,
		Priority:    item.Priority,
		IsReserved:  item.IsReserved,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}, nil
}

func (s *itemService) List(ctx context.Context, userID uuid.UUID, wishlistID uuid.UUID) (dto.WishListItemListResponse, error) {
	log := logging.FromContext(ctx)
	wishlist, err := s.wishlistRepo.GetByID(ctx, wishlistID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			log.Warn("Wishlist not found",
				slog.String("wishlist_id", wishlistID.String()),
			)
			return dto.WishListItemListResponse{}, domain.ErrInvalidCredentials
		}
		return dto.WishListItemListResponse{}, fmt.Errorf("get wishlist: %w", err)
	}

	if wishlist.UserID != userID {
		log.Warn("Wishlist user ID mismatch",
			slog.String("wishlist_id", wishlistID.String()),
			slog.String("expected_user_id", wishlist.UserID.String()),
			slog.String("actual_user_id", userID.String()),
		)
		return dto.WishListItemListResponse{}, domain.ErrForbidden
	}

	var items []domain.WishlistItem
	items, err = s.itemRepo.ListByWishlistID(ctx, wishlistID)
	if err != nil {
		return dto.WishListItemListResponse{}, fmt.Errorf("list items: %w", err)
	}

	response := dto.WishListItemListResponse{
		Items: make([]dto.WishListItemResponse, 0, len(items)),
	}
	for _, item := range items {
		response.Items = append(response.Items, dto.WishListItemResponse{
			ID:          item.ID,
			Title:       item.Title,
			Description: item.Description,
			ProductURL:  item.ProductURL,
			Priority:    item.Priority,
			IsReserved:  item.IsReserved,
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
		})
	}
	log.Info("successfully get items",
		slog.Int("item_count", len(items)),
		slog.String("wishlist_id", wishlistID.String()),
	)
	return response, nil
}

func (s *itemService) Update(ctx context.Context, userID uuid.UUID, wishlistID uuid.UUID, itemID uuid.UUID, req dto.WishListItemUpdateRequest) (dto.WishListItemResponse, error) {
	log := logging.FromContext(ctx)
	wishlist, err := s.wishlistRepo.GetByID(ctx, wishlistID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			log.Warn("Wishlist not found",
				slog.String("wishlist_id", wishlistID.String()),
			)
			return dto.WishListItemResponse{}, domain.ErrInvalidCredentials
		}
		return dto.WishListItemResponse{}, fmt.Errorf("get wishlist: %w", err)
	}

	if wishlist.UserID != userID {
		log.Warn("Wishlist user ID mismatch",
			slog.String("wishlist_id", wishlistID.String()),
			slog.String("expected_user_id", wishlist.UserID.String()),
			slog.String("actual_user_id", userID.String()),
		)
		return dto.WishListItemResponse{}, domain.ErrForbidden
	}

	var item domain.WishlistItem
	item, err = s.itemRepo.GetByID(ctx, itemID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			log.Warn("Item not found",
				slog.String("item_id", itemID.String()),
			)
			return dto.WishListItemResponse{}, domain.ErrNotFound
		}
		return dto.WishListItemResponse{}, fmt.Errorf("get item: %w", err)
	}
	if req.Title != nil {
		item.Title = *req.Title
	}
	if req.Description != nil {
		item.Description = req.Description
	}
	if req.ProductURL != nil {
		item.ProductURL = req.ProductURL
	}
	if req.Priority != nil {
		item.Priority = *req.Priority
	}

	var updatedItem domain.WishlistItem
	updatedItem, err = s.itemRepo.Update(ctx, item)
	if err != nil {
		return dto.WishListItemResponse{}, fmt.Errorf("update item: %w", err)
	}

	log.Info("Item updated",
		slog.String("item_id", itemID.String()),
		slog.String("wishlist_id", wishlistID.String()),
		slog.String("title", updatedItem.Title),
	)

	return dto.WishListItemResponse{
		ID:          updatedItem.ID,
		Title:       updatedItem.Title,
		Description: updatedItem.Description,
		ProductURL:  updatedItem.ProductURL,
		Priority:    updatedItem.Priority,
		IsReserved:  updatedItem.IsReserved,
		CreatedAt:   updatedItem.CreatedAt,
		UpdatedAt:   updatedItem.UpdatedAt,
	}, nil
}

func (s *itemService) Delete(ctx context.Context, userID uuid.UUID, wishlistID uuid.UUID, itemID uuid.UUID) error {
	log := logging.FromContext(ctx)
	wishlist, err := s.wishlistRepo.GetByID(ctx, wishlistID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			log.Warn("Wishlist not found",
				slog.String("wishlist_id", wishlistID.String()),
			)
			return domain.ErrInvalidCredentials
		}
		return fmt.Errorf("get wishlist: %w", err)
	}

	if wishlist.UserID != userID {
		log.Warn("Wishlist user ID mismatch",
			slog.String("wishlist_id", wishlistID.String()),
			slog.String("expected_user_id", wishlist.UserID.String()),
			slog.String("actual_user_id", userID.String()),
		)
		return domain.ErrForbidden
	}

	err = s.itemRepo.Delete(ctx, itemID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			log.Warn("Item not found",
				slog.String("item_id", itemID.String()),
			)
			return domain.ErrNotFound
		}
		return fmt.Errorf("delete item: %w", err)
	}
	log.Info("Item deleted",
		slog.String("item_id", itemID.String()),
		slog.String("wishlist_id", wishlistID.String()),
	)
	return nil
}

func (s *itemService) Reserve(ctx context.Context, token string, itemID uuid.UUID) error {
	log := logging.FromContext(ctx)

	wishlist, err := s.wishlistRepo.GetByShareToken(ctx, token)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			log.Warn("Wishlist not found for token",
				slog.String("token", token),
			)
			return domain.ErrInvalidCredentials
		}
		return fmt.Errorf("get wishlist by token: %w", err)
	}

	var item domain.WishlistItem
	item, err = s.itemRepo.GetByID(ctx, itemID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			log.Warn("Item not found",
				slog.String("item_id", itemID.String()),
			)
			return domain.ErrNotFound
		}
		return fmt.Errorf("get item: %w", err)
	}

	if item.WishlistID != wishlist.ID {
		log.Warn("Item does not belong to the wishlist",
			slog.String("item_id", itemID.String()),
			slog.String("wishlist_id", wishlist.ID.String()),
		)
		return domain.ErrForbidden
	}

	if item.IsReserved {
		log.Warn("Item already reserved",
			slog.String("item_id", itemID.String()),
		)
		return domain.ErrAlreadyReserved
	}

	err = s.itemRepo.Reserve(ctx, itemID)
	if err != nil {
		return fmt.Errorf("reserve item: %w", err)
	}

	log.Info("Item reserved",
		slog.String("item_id", itemID.String()),
		slog.String("wishlist_id", item.ID.String()),
		slog.String("title", item.Title),
	)

	return nil
}
