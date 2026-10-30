# internal/transport/http

HTTP 适配层：创建 Engine、组装路由、配置中间件，负责请求绑定和响应输出。

| 文件 | 职责 |
| --- | --- |
| `server.go` | Engine、中间件、路由分组；HTTP 服务的监听与启停 |
| `middleware/` | 中间件，一个中间件一个文件 |
| `handler/` | 业务接口的适配层，一个业务模块一个文件 |
| `jsoncodec/` | 请求绑定、业务错误的响应输出、统一响应格式 |

## 接口约定

- 业务接口统一使用 POST；健康检查等基础设施接口按语义选择 HTTP 方法，业务参数放在 JSON 请求体。
- 路径根据资源和接口语义确定，通常由资源名称和操作语义组成。
- 路径中的操作名称不要求与 Service 方法名完全一致，应以接口语义清晰为准。
- Handler 注册的是相对路径，`/api/v1` 等统一前缀由 `server.go` 中的路由分组添加。
    - 列表：`POST /api/v1/users/list`
    - 创建：`POST /api/v1/users/create`
    - 查询：`POST /api/v1/users/info`
    - 更新：`POST /api/v1/users/update`
- 请求体必须存在，空请求体按参数不合法处理。

## Handler 骨架

```go
// RegisterRoutes 声明本 UserHandler 负责的路由。
func (h *UserHandler) RegisterRoutes(r router.Router) {
	r.POST("/users/list", jsoncodec.Handler(h.List))
	r.POST("/users/info", jsoncodec.Handler(h.Get))
	r.POST("/users/create", jsoncodec.Handler(h.Create))
	r.POST("/users/update", jsoncodec.Handler(h.Update))
}

// Get 按 Id 返回用户。
func (h *UserHandler) Get(c *rx.Context, req *userdomain.GetRequest) (*userdomain.User, error) {
	return h.service.Get(c.Context(), req)
}
```

Handler 只做协议适配：绑定请求、调用用例、返回结果；业务逻辑一律放在 Service。

## jsoncodec：业务通用的编解码

项目业务统一使用 `jsoncodec` 完成请求绑定、业务错误响应和统一响应输出。

- 请求体绑定为对应的 Request。
- 请求绑定失败时，统一作为参数错误处理。
- 业务函数返回的错误交给 `ErrorEncoder` 统一编码。
- 业务函数成功返回时，返回值作为响应中的 `data`。
- 当前 `jsoncodec` 只对 POST、PUT、PATCH 请求解码请求体；其他 HTTP 方法不解码。

业务接口统一使用它，不要再为相同协议重复实现一套编解码逻辑。

**如果要对第三方提供接口，且需要满足第三方的协议要求**（例如报文格式不同、字段命名不同、需要签名或加密），请求绑定、响应结构和错误编码应作为一套协议适配处理。

在 `jsoncodec` 之外另建对应的编解码包，不要为了第三方协议修改业务通用的 `jsoncodec`，避免影响所有业务接口。

返回结构有特殊要求时，可以在对应 Handler 中明确处理错误并输出响应。

不要使用 `c.AbortWithError` 处理需要特殊响应格式的错误。错误编码器是 `rx.Engine` 级别的，调用该方法产生的错误会进入这里的 `ErrorEncoder`，最终仍按统一响应格式输出。

也不要在通用 `ErrorEncoder` 中按具体业务或第三方接口类型分支。通用层不应该感知第三方接口的特殊协议。

## 统一响应

当前统一响应结构由 `jsoncodec.Response` 定义：

```go
type Response struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}
```

成功时 HTTP 状态码为 `200 OK`。当前实现只设置 `Data`，因此 `code` 和 `message` 是空字符串；`data` 为 nil 时会因 `omitempty` 被省略。

成功响应示例：

```json
{
  "code": "",
  "message": "",
  "data": {
    "id": 1
  }
}
```

失败时由 `ErrorEncoder` 统一处理：

- `internal/foundation/errors.Error` 类型的错误返回 HTTP `400 Bad Request`，`code` 和 `message` 分别取错误对象中的值。
- 其他错误返回 HTTP `500 Internal Server Error`，`code` 为 `internal_server_error`，`message` 为 `internal server error`。
- 当前失败响应不设置 `data`，因此 JSON 中不会输出 `data` 字段。
- 响应体中的 `code` 是错误标识，与 HTTP 状态码没有对应关系。

业务错误响应示例：

```json
{
  "code": "invalid_argument",
  "message": "invalid argument"
}
```

服务端内部错误响应示例：

```json
{
  "code": "internal_server_error",
  "message": "internal server error"
}
```