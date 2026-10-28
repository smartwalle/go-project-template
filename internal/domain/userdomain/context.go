package userdomain

import "context"

// userKey 是承载用户信息的 context 键类型。
//
// 使用私有空结构体作为键，与其它包的用户键互不冲突。
type userKey struct{}

// ContextWithUser 返回一个携带用户信息的 Context。
//
// 原 Context 不会被修改；若 Context 中已有用户信息，本次注入的值会覆盖它。
func ContextWithUser(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, userKey{}, user)
}

// UserFromContext 从 Context 中取出用户信息。
//
// Context 中没有用户信息时返回 User 的零值。
func UserFromContext(ctx context.Context) User {
	user, _ := ctx.Value(userKey{}).(User)
	return user
}
