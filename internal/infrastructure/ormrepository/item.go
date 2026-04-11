package ormrepository

import (
	"context"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/MALblSH/Wishlist_API/internal/application/repository"
	"github.com/MALblSH/Wishlist_API/internal/domain"
)

type itemRepository struct {
	pool *pgxpool.Pool
}

func NewItemRepository(pool *pgxpool.Pool) repository.ItemRepository {
	return &itemRepository{pool: pool}
}

func (i itemRepository) Create(ctx context.Context, wishlistID uuid.UUID, title string, description *string, productURL *string, priority int) (domain.WishlistItem, error) {
	id := uuid.New()
	now := time.Now()

	var err error
	var sql string
	var args []interface{}
	sql, args, err = psql.
		Insert("wishlist_items").
		Columns("id", "wishlist_id", "title", "description", "product_url", "priority", "is_reserved", "created_at", "updated_at").
		Values(id, wishlistID, title, description, productURL, priority, false, now, now).
		ToSql()
	if err != nil {
		return domain.WishlistItem{}, fmt.Errorf("build item.Create query: %w", err)
	}

	_, err = i.pool.Exec(ctx, sql, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.WishlistItem{}, domain.ErrAlreadyExists
		}
		return domain.WishlistItem{}, fmt.Errorf("execute item.Create query: %w", err)
	}

	return domain.WishlistItem{
		ID:          id,
		WishlistID:  wishlistID,
		Title:       title,
		Description: description,
		ProductURL:  productURL,
		Priority:    priority,
		IsReserved:  false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil

}

func (i itemRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.WishlistItem, error) {
	var err error
	var sql string
	var args []interface{}
	sql, args, err = psql.
		Select("id", "wishlist_id", "title", "description", "product_url", "priority", "is_reserved", "created_at", "updated_at").
		From("wishlist_items").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return domain.WishlistItem{}, fmt.Errorf("build item.GetByID query: %w", err)
	}

	row := i.pool.QueryRow(ctx, sql, args...)
	var item domain.WishlistItem
	err = row.Scan(
		&item.ID,
		&item.WishlistID,
		&item.Title,
		&item.Description,
		&item.ProductURL,
		&item.Priority,
		&item.IsReserved,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.WishlistItem{}, domain.ErrNotFound
		}
		return domain.WishlistItem{}, fmt.Errorf("execute item.GetByID query: %w", err)
	}

	return item, nil
}

func (i itemRepository) ListByWishlistID(ctx context.Context, wishlistID uuid.UUID) ([]domain.WishlistItem, error) {
	var err error
	var sql string
	var args []interface{}
	sql, args, err = psql.
		Select("id", "wishlist_id", "title", "description", "product_url", "priority", "is_reserved", "created_at", "updated_at").
		From("wishlist_items").
		Where(sq.Eq{"wishlist_id": wishlistID}).
		OrderBy("priority DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build item.ListByWishlistID query: %w", err)
	}

	var rows pgx.Rows
	rows, err = i.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("execute item.ListByWishlistID query: %w", err)
	}
	defer rows.Close()

	var items []domain.WishlistItem
	for rows.Next() {
		var item domain.WishlistItem
		err = rows.Scan(
			&item.ID,
			&item.WishlistID,
			&item.Title,
			&item.Description,
			&item.ProductURL,
			&item.Priority,
			&item.IsReserved,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan item.ListByWishlistID row: %w", err)
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate item.ListByWishlistID rows: %w", err)
	}
	return items, nil
}

func (i itemRepository) Update(ctx context.Context, item domain.WishlistItem) (domain.WishlistItem, error) {
	now := time.Now()

	var err error
	var sql string
	var args []interface{}
	sql, args, err = psql.
		Update("wishlist_items").
		Set("title", item.Title).
		Set("description", item.Description).
		Set("product_url", item.ProductURL).
		Set("priority", item.Priority).
		Set("updated_at", now).
		Where(sq.Eq{"id": item.ID}).
		ToSql()
	if err != nil {
		return domain.WishlistItem{}, fmt.Errorf("build item.Update query: %w", err)
	}

	var cmdTag pgconn.CommandTag
	cmdTag, err = i.pool.Exec(ctx, sql, args...)
	if err != nil {
		return domain.WishlistItem{}, fmt.Errorf("execute item.Update query: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return domain.WishlistItem{}, domain.ErrNotFound
	}

	item.UpdatedAt = now
	return item, nil
}

func (i itemRepository) Delete(ctx context.Context, id uuid.UUID) error {
	var err error
	var sql string
	var args []interface{}
	sql, args, err = psql.
		Delete("wishlist_items").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build item.Delete query: %w", err)
	}

	var cmdTag pgconn.CommandTag
	cmdTag, err = i.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("execute item.Delete query: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (i itemRepository) Reserve(ctx context.Context, id uuid.UUID) error {
	now := time.Now()

	var err error
	var sql string
	var args []interface{}
	sql, args, err = psql.
		Update("wishlist_items").
		Set("is_reserved", true).
		Set("updated_at", now).
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build item.Reserve query: %w", err)
	}

	var cmdTag pgconn.CommandTag
	cmdTag, err = i.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("execute item.Reserve query: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}
