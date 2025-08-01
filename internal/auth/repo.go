package auth

import (
	"context"

	"sumunar-pos-core/pkg/db"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *RefreshToken) error
	FindByToken(ctx context.Context, token string) (*RefreshToken, error)
	Revoke(ctx context.Context, tokenID string) error
	RevokeAllByUser(ctx context.Context, userID string) error
}

type refreshTokenRepo struct {
	db db.DBTX
}

func NewRefreshTokenRepo(db db.DBTX) RefreshTokenRepository {
	return &refreshTokenRepo{db}
}

func (r *refreshTokenRepo) Create(ctx context.Context, t *RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (id, token, user_id, expires_at, is_active, created_by, created_at, updated_by, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	_, err := r.db.Exec(ctx, query, t.ID, t.Token, t.UserID, t.ExpiresAt, t.IsActive, t.CreatedBy, t.CreatedAt, t.UpdatedBy, t.UpdatedAt)
	return err
}

func (r *refreshTokenRepo) FindByToken(ctx context.Context, token string) (*RefreshToken, error) {
	var t RefreshToken
	query := `
		SELECT id, token, user_id, expires_at, is_active, created_at
		FROM refresh_tokens
		WHERE token = $1 AND is_active = true AND expires_at > NOW()`
	err := r.db.QueryRow(ctx, query, token).Scan(
		&t.ID, &t.Token, &t.UserID, &t.ExpiresAt, &t.IsActive, &t.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *refreshTokenRepo) Revoke(ctx context.Context, tokenID string) error {
	_, err := r.db.Exec(ctx, `UPDATE refresh_tokens SET is_active = false WHERE id = $1`, tokenID)
	return err
}

func (r *refreshTokenRepo) RevokeAllByUser(ctx context.Context, userID string) error {
	_, err := r.db.Exec(ctx, `UPDATE refresh_tokens SET is_active = false WHERE user_id = $1`, userID)
	return err
}
