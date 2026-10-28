package authdomain

import "context"

// authUserKey 是承载当前登录用户信息的 context 键类型。
//
// 用途：HTTP 中间件鉴权通过后，把当前登录用户写进 Context，
// 后续业务代码用 AuthUserFromContext 读取，
// 不必把用户信息一层层往下传。
//
// 使用私有空结构体作为键，与其它包的键互不冲突。
type authUserKey struct{}

// ContextWithAuthUser 返回一个携带当前登录用户信息的 Context。
//
// 由鉴权中间件在鉴权通过后调用；原 Context 不会被修改，
// 若 Context 中已有用户信息，本次注入的值会覆盖它。
func ContextWithAuthUser(ctx context.Context, user AuthUser) context.Context {
	return context.WithValue(ctx, authUserKey{}, user)
}

// AuthUserFromContext 从 Context 中取出当前登录用户信息。
//
// 供业务代码读取登录用户；Context 中没有用户信息
// （例如接口没有经过鉴权中间件）时返回 AuthUser 的零值。
func AuthUserFromContext(ctx context.Context) AuthUser {
	user, _ := ctx.Value(authUserKey{}).(AuthUser)
	return user
}
