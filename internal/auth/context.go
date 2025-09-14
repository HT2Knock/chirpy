package auth

import (
	"context"

	"github.com/google/uuid"
)

type contextKey string

const userIDKey = contextKey("userID")

func WithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func UserIDFromContext(ctx context.Context) uuid.UUID {
	userID := ctx.Value(userIDKey).(uuid.UUID)

	return userID
}
