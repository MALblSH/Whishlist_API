package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"

	"github.com/MALblSH/Wishlist_API/internal/application/generator"
	"github.com/MALblSH/Wishlist_API/internal/application/repository"
	"github.com/MALblSH/Wishlist_API/internal/domain"
	"github.com/MALblSH/Wishlist_API/internal/infrastructure/dto"
	"github.com/MALblSH/Wishlist_API/internal/infrastructure/logging"
)

type WishlistService interface {
	Create(ctx context.Context, userID uuid.UUID, req dto.WishListRequest) (dto.WishListResponse, error)
	List(ctx context.Context, userID uuid.UUID) (dto.WishListListResponse, error)
	Get(ctx context.Context, userID uuid.UUID, wishlistID uuid.UUID) (dto.WishListResponse, error)
	Update(ctx context.Context, userID uuid.UUID, wishlistID uuid.UUID, req dto.WishListUpdateRequest) (dto.WishListResponse, error)
	Delete(ctx context.Context, userID uuid.UUID, wishlistID uuid.UUID) error

	GetByToken(ctx context.Context, token string) (dto.PublicWishListResponse, error)
}

type wishlistService struct {
	wishlistRepo repository.WishlistRepository
	itemRepo     repository.ItemRepository
}

func NewWishlistService(wishlistRepo repository.WishlistRepository, itemRepo repository.ItemRepository) WishlistService {
	return &wishlistService{
		wishlistRepo: wishlistRepo,
		itemRepo:     itemRepo,
	}
}

func (s *wishlistService) Create(ctx context.Context, userID uuid.UUID, req dto.WishListRequest) (dto.WishListResponse, error) {
	log := logging.FromContext(ctx)

	if req.Title == "" {
		return dto.WishListResponse{}, domain.ErrInvalidInput
	}

	shareToken, err := generator.GenerateShareToken()
	if err != nil {
		log.Error("generate share token:",
			slog.String("error", err.Error()),
		)
		return dto.WishListResponse{}, fmt.Errorf("generate share token: %w", err)
	}

	var wishlist domain.Wishlist
	wishlist, err = s.wishlistRepo.Create(ctx, userID, req.Title, req.Description, req.EventDate, shareToken)
	if err != nil {
		if errors.Is(err, domain.ErrAlreadyExists) {
			log.Warn("wishlist already exists",
				slog.String("user_id", userID.String()),
			)
			return dto.WishListResponse{}, domain.ErrAlreadyExists
		}
		return dto.WishListResponse{}, fmt.Errorf("create wishlist: %w", err)
	}

	log.Info("wishlist created",
		slog.String("wishlist_id", wishlist.ID.String()),
		slog.String("user_id", userID.String()),
		slog.String("title", req.Title),
	)

	return dto.WishListResponse{
		ID:          wishlist.ID,
		Title:       wishlist.Title,
		Description: wishlist.Description,
		EventDate:   wishlist.EventDate,
		ShareToken:  wishlist.ShareToken,
		CreatedAt:   wishlist.CreatedAt,
		UpdatedAt:   wishlist.UpdatedAt,
	}, nil
}

func (s *wishlistService) List(ctx context.Context, userID uuid.UUID) (dto.WishListListResponse, error) {
	log := logging.FromContext(ctx)

	wishlists, err := s.wishlistRepo.ListByUserID(ctx, userID)
	if err != nil {
		return dto.WishListListResponse{}, fmt.Errorf("list wishlists: %w", err)
	}

	response := dto.WishListListResponse{
		WishLists: make([]dto.WishListResponse, 0, len(wishlists)),
	}

	for _, wishlist := range wishlists {
		response.WishLists = append(response.WishLists, dto.WishListResponse{
			ID:          wishlist.ID,
			Title:       wishlist.Title,
			Description: wishlist.Description,
			EventDate:   wishlist.EventDate,
			ShareToken:  wishlist.ShareToken,
			CreatedAt:   wishlist.CreatedAt,
			UpdatedAt:   wishlist.UpdatedAt,
		})
	}

	log.Info("successfully get wishlists",
		slog.Int("count", len(wishlists)),
	)
	return response, nil
}

func (s *wishlistService) Get(ctx context.Context, userID uuid.UUID, wishlistID uuid.UUID) (dto.WishListResponse, error) {
	log := logging.FromContext(ctx)

	wishlist, err := s.wishlistRepo.GetByID(ctx, wishlistID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			log.Warn("wishlist not found",
				slog.String("wishlist_id", wishlistID.String()),
			)
			return dto.WishListResponse{}, domain.ErrInvalidCredentials
		}
		return dto.WishListResponse{}, fmt.Errorf("get wishlist: %w", err)
	}

	if wishlist.UserID != userID {
		log.Warn("wishlist does not belong to user",
			slog.String("wishlist_id", wishlistID.String()),
			slog.String("user_id", userID.String()),
		)
		return dto.WishListResponse{}, domain.ErrInvalidCredentials
	}

	log.Info("successfully get wishlist",
		slog.String("wishlist_id", wishlistID.String()),
		slog.String("user_id", userID.String()),
	)
	return dto.WishListResponse{
		ID:          wishlist.ID,
		Title:       wishlist.Title,
		Description: wishlist.Description,
		EventDate:   wishlist.EventDate,
		ShareToken:  wishlist.ShareToken,
		CreatedAt:   wishlist.CreatedAt,
		UpdatedAt:   wishlist.UpdatedAt,
	}, nil
}

func (s *wishlistService) Update(ctx context.Context, userID uuid.UUID, wishlistID uuid.UUID, req dto.WishListUpdateRequest) (dto.WishListResponse, error) {
	log := logging.FromContext(ctx)

	wishlist, err := s.wishlistRepo.GetByID(ctx, wishlistID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			log.Warn("wishlist not found",
				slog.String("wishlist_id", wishlistID.String()),
			)
			return dto.WishListResponse{}, domain.ErrInvalidCredentials
		}
		return dto.WishListResponse{}, fmt.Errorf("get wishlist: %w", err)
	}

	if wishlist.UserID != userID {
		log.Warn("wishlist does not belong to user",
			slog.String("wishlist_id", wishlistID.String()),
			slog.String("user_id", userID.String()),
		)
		return dto.WishListResponse{}, domain.ErrInvalidCredentials
	}

	if req.Title != nil {
		wishlist.Title = *req.Title
	}
	if req.Description != nil {
		wishlist.Description = *req.Description
	}
	if req.EventDate != nil {
		wishlist.EventDate = *req.EventDate
	}

	var updated domain.Wishlist
	updated, err = s.wishlistRepo.Update(ctx, wishlist)
	if err != nil {
		return dto.WishListResponse{}, fmt.Errorf("update wishlist: %w", err)
	}

	log.Info("wishlist updated",
		slog.String("wishlist_id", wishlistID.String()),
		slog.String("user_id", userID.String()),
	)

	return dto.WishListResponse{
		ID:          updated.ID,
		Title:       updated.Title,
		Description: updated.Description,
		EventDate:   updated.EventDate,
		ShareToken:  updated.ShareToken,
		CreatedAt:   updated.CreatedAt,
		UpdatedAt:   updated.UpdatedAt,
	}, nil
}

func (s *wishlistService) Delete(ctx context.Context, userID uuid.UUID, wishlistID uuid.UUID) error {
	log := logging.FromContext(ctx)

	wishlist, err := s.wishlistRepo.GetByID(ctx, wishlistID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			log.Warn("wishlist not found",
				slog.String("wishlist_id", wishlistID.String()),
			)
			return domain.ErrInvalidCredentials
		}
		return fmt.Errorf("get wishlist: %w", err)
	}

	if wishlist.UserID != userID {
		log.Warn("wishlist does not belong to user",
			slog.String("wishlist_id", wishlistID.String()),
			slog.String("user_id", userID.String()),
		)
		return domain.ErrInvalidCredentials
	}

	err = s.wishlistRepo.Delete(ctx, wishlistID)
	if err != nil {
		return fmt.Errorf("delete wishlist: %w", err)
	}

	log.Info("wishlist deleted",
		slog.String("wishlist_id", wishlistID.String()),
		slog.String("user_id", userID.String()),
	)
	return nil
}

func (s *wishlistService) GetByToken(ctx context.Context, token string) (dto.PublicWishListResponse, error) {
	log := logging.FromContext(ctx)

	wishlist, err := s.wishlistRepo.GetByShareToken(ctx, token)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			log.Warn("wishlist not found",
				slog.String("share_token", token),
			)
			return dto.PublicWishListResponse{}, domain.ErrInvalidCredentials
		}
		return dto.PublicWishListResponse{}, fmt.Errorf("get wishlist by share token: %w", err)
	}

	var items []domain.WishlistItem
	items, err = s.itemRepo.ListByWishlistID(ctx, wishlist.ID)
	if err != nil {
		return dto.PublicWishListResponse{}, fmt.Errorf("list items by wishlist id: %w", err)
	}

	responseItems := make([]dto.WishListItemResponse, 0, len(items))
	for _, item := range items {
		responseItems = append(responseItems, dto.WishListItemResponse{
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

	return dto.PublicWishListResponse{
		ID:          wishlist.ID,
		Title:       wishlist.Title,
		Description: wishlist.Description,
		EventDate:   wishlist.EventDate,
		Items:       responseItems,
	}, nil
}
