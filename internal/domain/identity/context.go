package identity

import "context"

type userKey struct{}

func WithUserId(ctx context.Context, userId uint64) context.Context {
	return context.WithValue(ctx, userKey{}, userId)
}

func UserId(ctx context.Context) uint64 {
	var value, _ = ctx.Value(userKey{}).(uint64)
	return value
}
