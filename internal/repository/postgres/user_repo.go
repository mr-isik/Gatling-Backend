package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mr-isik/gatling-backend/internal/domain"
	"github.com/mr-isik/gatling-backend/internal/repository"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, u *domain.User) (*domain.User, error) {
	query := `INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, created_at`
	err := r.db.QueryRow(ctx, query, u.Email, u.PasswordHash).Scan(&u.ID, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	query := `SELECT id, email, password_hash, created_at FROM users WHERE id = $1`
	var u domain.User
	err := r.db.QueryRow(ctx, query, id).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, email, password_hash, created_at FROM users WHERE email = $1`
	var u domain.User
	err := r.db.QueryRow(ctx, query, email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) CreateAPIKey(ctx context.Context, k *domain.APIKey) (*domain.APIKey, error) {
	query := `INSERT INTO api_keys (user_id, key_hash, name) VALUES ($1, $2, $3) RETURNING id, created_at`
	err := r.db.QueryRow(ctx, query, k.UserID, k.KeyHash, k.Name).Scan(&k.ID, &k.CreatedAt)
	if err != nil {
		return nil, err
	}
	return k, nil
}

func (r *userRepository) GetAPIKeyByHash(ctx context.Context, hash string) (*domain.APIKey, error) {
	query := `SELECT id, user_id, key_hash, name, created_at FROM api_keys WHERE key_hash = $1`
	var k domain.APIKey
	err := r.db.QueryRow(ctx, query, hash).Scan(&k.ID, &k.UserID, &k.KeyHash, &k.Name, &k.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &k, nil
}

func (r *userRepository) DeleteAPIKey(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, "DELETE FROM api_keys WHERE id = $1", id)
	return err
}
