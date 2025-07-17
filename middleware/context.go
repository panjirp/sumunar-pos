package middleware

import (
	"context"
	"errors"
)

// contextKey untuk menghindari bentrok key antar package
type contextKey string

const (
	ContextKeyUserID   contextKey = "user_id"
	ContextKeyUsername contextKey = "username"
	ContextKeyEmail    contextKey = "email"
	ContextKeyRole     contextKey = "role"
)

// SetDataToContext menyisipkan data ke context
func SetDataToContext(ctx context.Context, key, value string) context.Context {
	return context.WithValue(ctx, key, value)
}

// GetUserIDFromContext mengembalikan UUID user dari context
func GetUserIDFromContext(ctx context.Context) (string, error) {
	raw := ctx.Value(ContextKeyUserID)
	id, ok := raw.(string)
	if !ok || id == "" {
		return "", errors.New("user_id not found in context")
	}
	return id, nil
}

func GetUsernameFromContext(ctx context.Context) (string, error) {
	raw := ctx.Value(ContextKeyUsername)
	username, ok := raw.(string)
	if !ok || username == "" {
		return "", errors.New("username not found in context")
	}
	return username, nil
}

func GetRoleFromContext(ctx context.Context) (string, error) {
	raw := ctx.Value(ContextKeyRole)
	role, ok := raw.(string)
	if !ok || role == "" {
		return "", errors.New("role not found in context")
	}
	return role, nil
}
