package metadata

import "context"

// requestIdKey 是承载请求 id 的 context 键类型。
//
// 使用私有空结构体作为键，与其它包的键互不冲突。
type requestIdKey struct{}

// WithRequestId 返回一个携带请求 id 的 Context。
//
// 原 Context 不会被修改；若 Context 中已有请求 id，本次注入的值会覆盖它。
func WithRequestId(ctx context.Context, rid string) context.Context {
	return context.WithValue(ctx, requestIdKey{}, rid)
}

// RequestIdFrom 从 Context 中取出请求 id。
//
// Context 中没有请求 id 时返回空字符串。
func RequestIdFrom(ctx context.Context) string {
	var rid, _ = ctx.Value(requestIdKey{}).(string)
	return rid
}
