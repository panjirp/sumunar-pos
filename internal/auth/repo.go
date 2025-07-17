package auth

import (
	"context"

	"sumunar-pos-core/pkg/db"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *RefreshToken) error
	FindByToken(ctx context.Context, token string) (*RefreshToken, error)
	Revoke(ctx context.Context, tokenID string) error
}

type refreshTokenRepo struct {
	db db.DBTX
}

func NewRefreshTokenRepo(db db.DBTX) RefreshTokenRepository {
	return &refreshTokenRepo{db}
}

func (r *refreshTokenRepo) Create(ctx context.Context, t *RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (id, token, user_id, expires_at, revoked, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(ctx, query, t.ID, t.Token, t.UserID, t.ExpiresAt, t.Revoked, t.CreatedAt)
	return err
}

func (r *refreshTokenRepo) FindByToken(ctx context.Context, token string) (*RefreshToken, error) {
	var t RefreshToken
	query := `
		SELECT id, token, user_id, expires_at, revoked, created_at
		FROM refresh_tokens
		WHERE token = $1 AND revoked = false AND expires_at > NOW()`
	err := r.db.QueryRow(ctx, query, token).Scan(
		&t.ID, &t.Token, &t.UserID, &t.ExpiresAt, &t.Revoked, &t.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *refreshTokenRepo) Revoke(ctx context.Context, tokenID string) error {
	_, err := r.db.Exec(ctx, `UPDATE refresh_tokens SET revoked = true WHERE id = $1`, tokenID)
	return err
}

func (r *refreshTokenRepo) RevokeAllByUser(ctx context.Context, userID string) error {
	_, err := r.db.Exec(ctx, `UPDATE refresh_tokens SET revoked = true WHERE user_id = $1`, userID)
	return err
}
