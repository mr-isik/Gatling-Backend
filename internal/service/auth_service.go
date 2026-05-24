package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/mail"
	"strings"

	"github.com/mr-isik/gatling-backend/internal/domain"
	"github.com/mr-isik/gatling-backend/internal/infra"
	"github.com/mr-isik/gatling-backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  repository.UserRepository
	jwtSecret string
	jwtCfg    infra.JWTConfig
}

func NewAuthService(userRepo repository.UserRepository, cfg infra.JWTConfig) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: cfg.Secret,
		jwtCfg:    cfg,
	}
}

func (s *AuthService) Register(ctx context.Context, email, password string) (*TokenPair, error) {
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, errors.New("invalid email format")
	}
	if len(password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:        strings.ToLower(email),
		PasswordHash: string(hash),
	}

	createdUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return GenerateTokenPair(createdUser.ID, createdUser.Email, s.jwtSecret, s.jwtCfg.AccessExpiryMin, s.jwtCfg.RefreshExpiryDay)
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*TokenPair, error) {
	user, err := s.userRepo.GetByEmail(ctx, strings.ToLower(email))
	if err != nil {
		return nil, domain.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, domain.ErrUnauthorized
	}

	return GenerateTokenPair(user.ID, user.Email, s.jwtSecret, s.jwtCfg.AccessExpiryMin, s.jwtCfg.RefreshExpiryDay)
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshTokenStr string) (*TokenPair, error) {
	claims, err := ValidateToken(refreshTokenStr, s.jwtSecret)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}
	return GenerateTokenPair(claims.UserID, claims.Email, s.jwtSecret, s.jwtCfg.AccessExpiryMin, s.jwtCfg.RefreshExpiryDay)
}

func (s *AuthService) CreateAPIKey(ctx context.Context, userID, name string) (string, *domain.APIKey, error) {
	rawKeyBytes := make([]byte, 32)
	if _, err := rand.Read(rawKeyBytes); err != nil {
		return "", nil, err
	}
	rawKey := hex.EncodeToString(rawKeyBytes)

	hash := sha256.Sum256([]byte(rawKey))
	keyHash := hex.EncodeToString(hash[:])

	apiKey := &domain.APIKey{
		UserID:  userID,
		KeyHash: keyHash,
		Name:    name,
	}

	createdKey, err := s.userRepo.CreateAPIKey(ctx, apiKey)
	if err != nil {
		return "", nil, err
	}
	return rawKey, createdKey, nil
}

func (s *AuthService) DeleteAPIKey(ctx context.Context, keyID string) error {
	return s.userRepo.DeleteAPIKey(ctx, keyID)
}
