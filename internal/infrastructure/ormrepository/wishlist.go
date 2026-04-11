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

type wishlistRepository struct {
	pool *pgxpool.Pool
}

func NewWishlistRepository(pool *pgxpool.Pool) repository.WishlistRepository {
	return &wishlistRepository{pool: pool}
}

func (w wishlistRepository) Create(ctx context.Context, userID uuid.UUID, title string, description string, eventDate time.Time, shareToken string) (domain.Wishlist, error) {
	id := uuid.New()
	now := time.Now()

	var err error
	var sql string
	var args []interface{}
	sql, args, err = psql.
		Insert("wishlists").
		Columns("id", "user_id", "title", "description", "event_date", "share_token", "created_at", "updated_at").
		Values(id, userID, title, description, eventDate, shareToken, now, now).
		ToSql()
	if err != nil {
		return domain.Wishlist{}, fmt.Errorf("build wishlist.Create query: %w", err)
	}

	_, err = w.pool.Exec(ctx, sql, args...)
	if err != nil {
		return domain.Wishlist{}, fmt.Errorf("execute wishlist.Create query: %w", err)
	}

	return domain.Wishlist{
		ID:          id,
		UserID:      userID,
		Title:       title,
		Description: description,
		EventDate:   eventDate,
		ShareToken:  shareToken,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (w wishlistRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.Wishlist, error) {
	var err error
	var sql string
	var args []interface{}
	sql, args, err = psql.
		Select("id", "user_id", "title", "description", "event_date", "share_token", "created_at", "updated_at").
		From("wishlists").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return domain.Wishlist{}, fmt.Errorf("build wishlist.GetByID query: %w", err)
	}

	return w.scanWishlist(w.pool.QueryRow(ctx, sql, args...))
}

func (w wishlistRepository) GetByShareToken(ctx context.Context, token string) (domain.Wishlist, error) {
	var err error
	var sql string
	var args []interface{}
	sql, args, err = psql.
		Select("id", "user_id", "title", "description", "event_date", "share_token", "created_at", "updated_at").
		From("wishlists").
		Where(sq.Eq{"share_token": token}).
		ToSql()
	if err != nil {
		return domain.Wishlist{}, fmt.Errorf("build wishlist.GetByShareToken query: %w", err)
	}

	return w.scanWishlist(w.pool.QueryRow(ctx, sql, args...))
}

func (w wishlistRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Wishlist, error) {
	var err error
	var sql string
	var args []interface{}
	sql, args, err = psql.
		Select("id", "user_id", "title", "description", "event_date", "share_token", "created_at", "updated_at").
		From("wishlists").
		Where(sq.Eq{"user_id": userID}).
		OrderBy("created_at DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build wishlist.ListByUserID query: %w", err)
	}

	var rows pgx.Rows
	rows, err = w.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("execute wishlist.ListByUserID query: %w", err)
	}
	defer rows.Close()

	var wishlists []domain.Wishlist
	for rows.Next() {
		var wishlist domain.Wishlist
		err = rows.Scan(
			&wishlist.ID,
			&wishlist.UserID,
			&wishlist.Title,
			&wishlist.Description,
			&wishlist.EventDate,
			&wishlist.ShareToken,
			&wishlist.CreatedAt,
			&wishlist.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan wishlist.ListByUserID row: %w", err)
		}
		wishlists = append(wishlists, wishlist)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate wishlist.ListByUserID rows: %w", err)
	}

	return wishlists, nil
}

func (w wishlistRepository) Update(ctx context.Context, wishlist domain.Wishlist) (domain.Wishlist, error) {
	now := time.Now()

	var err error
	var sql string
	var args []interface{}
	sql, args, err = psql.
		Update("wishlists").
		Set("title", wishlist.Title).
		Set("description", wishlist.Description).
		Set("event_date", wishlist.EventDate).
		Set("updated_at", now).
		Where(sq.Eq{"id": wishlist.ID}).
		ToSql()
	if err != nil {
		return domain.Wishlist{}, fmt.Errorf("build wishlist.Update query: %w", err)
	}

	var cmdTag pgconn.CommandTag
	cmdTag, err = w.pool.Exec(ctx, sql, args...)
	if err != nil {
		return domain.Wishlist{}, fmt.Errorf("execute wishlist.Update query: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return domain.Wishlist{}, domain.ErrNotFound
	}

	wishlist.UpdatedAt = now
	return wishlist, nil
}

func (w wishlistRepository) Delete(ctx context.Context, id uuid.UUID) error {
	var err error
	var sql string
	var args []interface{}
	sql, args, err = psql.
		Delete("wishlists").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("build wishlist.Delete query: %w", err)
	}

	var cmdTag pgconn.CommandTag
	cmdTag, err = w.pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("execute wishlist.Delete query: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (w wishlistRepository) scanWishlist(row pgx.Row) (domain.Wishlist, error) {
	var wishlist domain.Wishlist
	err := row.Scan(
		&wishlist.ID,
		&wishlist.UserID,
		&wishlist.Title,
		&wishlist.Description,
		&wishlist.EventDate,
		&wishlist.ShareToken,
		&wishlist.CreatedAt,
		&wishlist.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Wishlist{}, domain.ErrNotFound
		}
		return domain.Wishlist{}, fmt.Errorf("scan wishlist row: %w", err)
	}
	return wishlist, nil
}
