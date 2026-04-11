package ormrepository

import (
	"context"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/MALblSH/Wishlist_API/internal/application/repository"
	"github.com/MALblSH/Wishlist_API/internal/domain"
)

var (
	psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
)

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) repository.UserRepository {
	return &userRepository{pool: pool}
}

func (u userRepository) Create(ctx context.Context, email string, passwordHash string) error {
	var err error
	var sql string
	var args []interface{}
	sql, args, err = psql.
		Insert("users").
		Columns("id", "email", "password_hash").
		Values(uuid.New(), email, passwordHash).
		ToSql()
	if err != nil {
		return fmt.Errorf("build user.Create  query: %w", err)
	}

	_, err = u.pool.Exec(ctx, sql, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrAlreadyExists
		}
		return fmt.Errorf("execute user.Create query: %w", err)
	}
	return nil
}

func (u userRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	var err error
	var sql string
	var args []interface{}
	sql, args, err = psql.
		Select("id", "email", "password_hash", "created_at").
		From("users").
		Where(sq.Eq{"email": email}).
		ToSql()
	if err != nil {
		return domain.User{}, fmt.Errorf("build user.GetByEmail query: %w", err)
	}

	row := u.pool.QueryRow(ctx, sql, args...)

	var user domain.User
	err = row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrNotFound
		}
		return domain.User{}, fmt.Errorf("execute user.GetByEmail query: %w", err)
	}
	return user, nil
}

func (u userRepository) GetByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	var err error
	var sql string
	var args []interface{}
	sql, args, err = psql.
		Select("id", "email", "password_hash", "created_at").
		From("users").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return domain.User{}, fmt.Errorf("build user.GetByID query: %w", err)
	}

	row := u.pool.QueryRow(ctx, sql, args...)

	var user domain.User
	err = row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrNotFound
		}
		return domain.User{}, fmt.Errorf("execute user.GetByID query: %w", err)
	}
	return user, nil
}
