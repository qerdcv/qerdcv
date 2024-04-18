package domain

import (
	"context"
	"time"
)

type userSessionCtxKeyType string

var userSessionCtxKey userSessionCtxKeyType = "user-session"

type UserSession struct {
	ID        int
	UserID    int
	CreatedAt time.Time
	ExpiresAt time.Time
}

func UserSessionFromContext(ctx context.Context) UserSession {
	return ctx.Value(userSessionCtxKey).(UserSession)
}

func ContextWithUserSession(ctx context.Context, session UserSession) context.Context {
	return context.WithValue(ctx, userSessionCtxKey, session)
}
