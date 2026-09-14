package authctx

import (
	"context"


)

type ctxKey struct{}

func WithUserId(ctx context.Context, UserID string) context.Context {
	return context.WithValue(ctx, ctxKey{}, UserID)
}

func FromUserId(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(ctxKey{}).(string)
	return userID, ok

}
