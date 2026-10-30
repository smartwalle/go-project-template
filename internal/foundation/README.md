# internal/foundation

Foundation 提供跨业务复用的通用技术能力，不依赖具体业务领域，也不包含业务专属模型、规则和实现。

Foundation 中的能力由需要它们的层按需依赖，不要求所有层都依赖 Foundation。

| 包 | 职责 |
| --- | --- |
| `errors/` | 通用错误类型和跨业务通用错误码 |
| `database/` | 数据库连接、连接池参数、关闭；唯一允许出现驱动和 DSN 的地方 |
| `logging/` | 全局 slog logger 的初始化、SQL 日志接入 |
| `metadata/` | 通用的上下文数据管理（请求 ID 等） |
| `pagination/` | 分页参数与分页计算 |

## 错误

这里定义跨业务可复用的通用错误语义，例如参数错误、资源不存在、资源已存在等。

```go
// 用例自己产生的通用错误。
return nil, errors.NotFound("user not found")
```

约定：

- `code` 表达错误语义，`message` 描述具体细节，调用方按 `code` 判断。
- 通用错误语义使用这里提供的构造函数。
- 业务专属的错误语义在对应领域包中定义，使用带业务前缀的专属 code。
- 判断错误语义按 code 比较，不按错误值比较：message 可以随业务变化。
- 不要把具体业务领域的错误定义放到 Foundation。

## 分页

分页参数和计算逻辑由 `pagination` 提供。

```go
type ListRequest struct {
	pagination.Page

	// Name 是用户名称的匹配条件，为空表示不过滤。
	Name string `json:"name"`
}

func (r *ListRequest) Validate() error {
	if r == nil {
		return errors.MissingArgument("request is required")
	}
	r.Page.Normalize()
	return nil
}
```

约定：

- `Normalize` 归一化页码和每页条数，校验通过后直接使用。
- `Limit`、`Offset` 计算查询条数上限和偏移量，可以直接传给查询构造器。
- `OutOfRange` 判断页码是否超出总页数。
- 分页参数内嵌在请求结构体中，客户端看到的字段仍然是 `page`、`page_size`。

## 日志

日志消息（msg）使用中文，字段名（key）保持英文，错误文案原样透传。

具体日志约定见 `logging` 包的注释。

## 上下文数据

通用的上下文数据在这里读写，例如请求 ID。

每份数据使用一个私有 key，写入方法使用 `WithXxx` 命名，读取方法使用 `XxxFrom` 命名。

```go
ctx = metadata.WithRequestId(ctx, rid)
rid := metadata.RequestIdFrom(ctx)
```

约定：

- 通用上下文数据放在 `internal/foundation/metadata`。
- 业务相关的上下文数据不放在这里，由 `internal/domain/metadata` 统一管理。
- 当前请求 ID 使用 `string` 类型。
- 请求 ID 由 HTTP 中间件写入，日志中由 logkit 自动注入，业务代码不需要自行拼接。