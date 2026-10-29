# internal/transport/http

HTTP 适配层：创建 Engine、组装路由、配置中间件，负责请求绑定和响应输出。

| 文件 | 职责 |
| --- | --- |
| `server.go` | Engine、中间件、路由分组；HTTP 服务的监听与启停 |
| `middleware/` | 中间件，一个中间件一个文件 |
| `handler/` | 业务接口的适配层，一个业务模块一个文件 |
| `jsoncodec/` | 请求绑定、业务错误的响应输出、统一响应格式 |

## 接口约定

- 全部用 POST，参数放在 JSON 请求体
- 路径 = 资源 + 用例名，用例名与 Service 的方法名一致
  - 列表 `POST /users/list`
  - 创建 `POST /users/create`
  - 查询 `POST /users/info`
  - 更新 `POST /users/update`
- 请求体必须存在，空请求体按参数不合法处理
- 路径前缀（`/api/v1`）由这一层统一加，Handler 只声明自己的资源路径

## handler 骨架

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

handler 只做协议适配：绑定请求、调用用例、返回结果；业务逻辑一律在 service。

## jsoncodec：业务通用的编解码

项目业务统一的编解码方式：请求体绑定成 Request、业务错误统一输出成响应、响应体统一成 `{code, message, data}`。
业务接口都用它，不要再写一套。

**如果要对第三方提供接口，且要满足第三方的要求**（报文格式不同、字段命名不同、要签名或加密），
在 jsoncodec 之外另建一个包自定义编解码方式，不要改 jsoncodec 本身——
它被所有业务接口共用，改了会影响全部接口。

## 统一响应

```json
{"code": "", "message": "", "data": {}}
```

- 失败时 code、message 说明原因，data 为空
