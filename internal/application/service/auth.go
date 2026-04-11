package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"golang.org/x/crypto/bcrypt"

	"github.com/MALblSH/Wishlist_API/internal/application/repository"
	"github.com/MALblSH/Wishlist_API/internal/domain"
	"github.com/MALblSH/Wishlist_API/internal/infrastructure/dto"
	"github.com/MALblSH/Wishlist_API/internal/infrastructure/logging"
	"github.com/MALblSH/Wishlist_API/internal/infrastructure/tokens/access"
)

type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) error
	Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error)
}

type authService struct {
	userRepo   repository.UserRepository
	jwtManager *access.Manager
}

func NewAuthService(userRepo repository.UserRepository, jwtManager *access.Manager) AuthService {
	return &authService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) error {
	log := logging.FromContext(ctx)

	if req.Email == "" || req.Password == "" {
		log.Warn("registration failed: email or password is empty")
		return domain.ErrInvalidInput
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("hashing password failed",
			slog.String("error", err.Error()),
		)
		return fmt.Errorf("hash password: %w", err)
	}

	err = s.userRepo.Create(ctx, req.Email, string(passwordHash))
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	log.Info("user created",
		slog.String("email", req.Email),
	)

	return nil
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (dto.LoginResponse, error) {
	log := logging.FromContext(ctx)

	if req.Email == "" || req.Password == "" {
		log.Warn("login failed: email or password is empty")
		return dto.LoginResponse{}, domain.ErrInvalidInput
	}

	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			log.Warn("user not found",
				slog.String("email", req.Email),
			)
			return dto.LoginResponse{}, domain.ErrInvalidCredentials
		}
		return dto.LoginResponse{}, fmt.Errorf("get user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		log.Warn("invalid password",
			slog.String("email", req.Email),
		)
		return dto.LoginResponse{}, domain.ErrInvalidCredentials
	}

	var accessToken string
	accessToken, err = s.jwtManager.GenerateAccessToken(user.ID)
	if err != nil {
		log.Error("generate share token:",
			slog.String("error", err.Error()),
		)
		return dto.LoginResponse{}, fmt.Errorf("generate share token: %w", err)
	}

	log.Info("successfully logged in",
		slog.String("email", req.Email),
	)

	return dto.LoginResponse{
		AccessToken: accessToken,
		ExpiresIn:   int64(s.jwtManager.AccessTokenTTL.Seconds()),
	}, nil

}
