# internal/foundation

平台能力：与业务无关、被所有层复用的技术设施。这里不出现业务概念。

| 包 | 职责 |
| --- | --- |
| `errors/` | 业务错误类型和通用错误码 |
| `database/` | 数据库连接、连接池参数、关闭；唯一允许出现驱动和 DSN 的地方 |
| `logging/` | 全局 slog logger 的初始化、SQL 日志接入 |
| `metadata/` | 通用的上下文数据管理（请求 id 等） |
| `pagination/` | 分页参数与分页计算 |

## 业务错误

```go
// 用例自己产生的错误
return nil, errors.NotFound("user not found")
```

- code 表达语义、message 描述细节，调用方按 code 判断
- 通用语义用这里的构造函数；业务专属的语义在领域包里用 `New` 定义专属 code，code 带业务前缀
- 判断错误语义按 code 比，不按值比：message 随业务变化

## 分页

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

- `Normalize` 归一化页码和每页条数，校验通过后直接用
- `Limit`、`Offset` 算条数上限和偏移量，直接喂给查询构造器
- `OutOfRange` 判断页码是否超出总页数
- 分页参数内嵌在请求结构体里，客户端看到的字段还是 page、page_size

## 日志

消息（msg）用中文，字段名（key）保持英文，错误文案原样透传；细则见 `logging` 包的注释。

## 上下文数据

通用的上下文数据在这里读写：每份数据一个私有键，写入用 `WithXxx`，读取用 `XxxFrom`。

```go
ctx = metadata.WithRequestId(ctx, rid)
rid := metadata.RequestIdFrom(ctx)
```

请求 id 由 HTTP 中间件写入，日志里由 logkit 自动注入，不用业务代码自己拼。
业务相关的上下文数据不放这里，由 `internal/domain/metadata` 统一管理。
