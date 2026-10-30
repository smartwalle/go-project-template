
# AGENTS.md

本文件是本仓库的编码规范，供 AI Coding Agent 和新同事直接照着执行。
这里定规则，各目录的 README 给骨架与用法，细则写在对应代码的注释里。

## 项目是什么

Go 后端服务模板：分层清晰、约定统一，业务模块自带一套完整示例。

## 目录与分层

```text
cmd/                       进程入口：加载配置、准备环境、启动服务
internal/app/              依赖组装与生命周期
internal/config/           配置加载与校验
internal/domain/<模块>/    领域契约：Entity、Model、Request/Response、Repository/Service 接口
internal/foundation/<能力>/ 平台能力：errors、database、logging、metadata、pagination
internal/repository/       数据访问实现
internal/service/          业务用例实现
internal/transport/http/   HTTP 适配层：handler、jsoncodec、中间件
```

依赖方向，只能从上往下：

```text
transport/http  →  service  →  domain
                      ↘        ↙
                     foundation（被所有层依赖）
```

- `domain` 不 import transport、service、repository 这些实现包，也不 import 数据库驱动；仓储契约基于 dbs 的通用 Repository
- `foundation` 不依赖业务代码
- `app` 只负责把各层组装起来，不写业务逻辑
- 全局共享的资源（配置、数据库连接等）挂在 `app.Environment` 上：`Prepare` 创建，`Clear` 释放
- 跨层调用靠接口：接口定义在 `domain`，实现放在 `repository` / `service`

## 命名

| 类型 | 约定 | 位置 |
| --- | --- | --- |
| Entity | 每张表一个 XXXEntity（表名驼峰 + Entity 后缀），带 sql tag，实现 TableName、PrimaryKey | `internal/domain/<模块>/entity.go` |
| Model | 整理、组合数据的结构体，不绑定数据库表 | `internal/domain/<模块>/model.go` |
| Request | XXXRequest，XXX 是用例名 | `internal/domain/<模块>/service.go` |
| Response | XXXResponse，XXX 是用例名 | `internal/domain/<模块>/service.go` |
| Repository | 每张表一个仓储接口 | `internal/domain/<模块>/repository.go` |
| Service | 每个方法是一个用例 | `internal/domain/<模块>/service.go` |

领域共用内容：所有领域共用的放 `internal/domain`（package domain），只在一个模块里用的放在模块自己的包里。

标识符（id）：

- Entity、Model 中的标识符类型统一使用 `uint64`
- Request 中接收的标识符可以使用 `int64`，以便校验负数
- 有效性判断统一用 `> 0` / `<= 0`，不要只用 `== 0`：类型改成有符号整型时，负数也要挡住

## 行为约定

**请求参数**

- 一律用指针传递：`Create(ctx context.Context, req *CreateRequest)`
- 每个 Request 必须实现 `Validate() error`
- Validate 做三件事：nil 判断、校验、归一化（清洗后写回字段）
- service 的第一个动作是 `req.Validate()`，校验不通过直接返回错误

**上下文**

- 业务相关的上下文数据（当前登录用户、租户等）的读写统一放 `internal/domain/metadata`，数据类型留在各自的领域包
- 键用私有空结构体，不同包的键不会互相冲突；取不到时返回零值，不返回错误

**错误**

- 数据访问（仓储）返回的错误原样透传，不包装、不转换
- 用例自己产生的错误用 `internal/foundation/errors` 里的业务错误：NotFound、AlreadyExists、InvalidArgument、MissingArgument 等
- 通用语义用这些构造函数；业务专属的语义在领域包里用 `New` 定义专属 code，code 带业务前缀（如 `user_not_found`）
- 判断错误语义按 code 比，不按值比：message 随业务变化，取出 `Error` 后比 `Code()`

**查询**

- 查询方法：会被复用的查询进仓储契约，一次性的在 service 里用 `SelectBuilder` 拼
- 条件一律用占位符传参，不要把用户输入拼进 SQL
- 存在性判断：查不到（nil）或主键 `<= 0`，都算没有这条记录

**事务**

- 多个写入要一起成功或一起失败时，放进一个事务：仓储的 `Transaction(ctx, func(txCtx context.Context) error { ... })`
- 事务回调里必须用回调收到的 ctx（txCtx）：用外层的 ctx，语句会跑到事务外
- 回调返回错误即回滚，返回 nil 即提交

**变量**

- 不要和外层变量重名，避免遮蔽，尤其是事务回调里的 `ctx`、`err`

**注释**

- 中文，以标识符开头，句末句号
- 需要解释"为什么"时另起一段
- 只留最终结论，不写讨论过程和方案取舍

## HTTP 接口约定

- 业务接口统一使用 POST，参数一律放在 JSON 请求体
- 健康检查等基础设施接口可以根据接口语义使用合适的 HTTP 方法
- 路径 = 资源 + 用例名，用例名与 Service 的方法名一致
    - 列表 `POST /users/list`
    - 创建 `POST /users/create`
    - 查询 `POST /users/info`
    - 更新 `POST /users/update`
- 业务接口请求体必须存在，空请求体按参数不合法处理
- handler 方法签名统一：`func (h *XxxHandler) Yyy(c *rx.Context, req *<模块>.XxxRequest) (*<模块>.Xxx, error)`
- handler 只做协议适配：绑定请求、调用用例、返回结果；业务逻辑一律在 service
- 业务错误的响应由 jsoncodec 统一输出，handler 不自己写状态码
- 业务接口统一用 jsoncodec 绑定请求、输出响应；对接第三方、报文格式不同时另建一套，不改 jsoncodec
- 返回结构有多种时，在自己的 handler 里处理错误并输出，不用 `AbortWithError`
- 路由在 handler 的 `RegisterRoutes` 里注册，路径前缀由 transport/http 统一加

## 日志

- 消息（msg）用中文，字段名（key）保持英文
- 错误文案原样透传
- 请求 id 通过 metadata 包在 context 里传递，日志里由 logkit 自动注入

## 新增一个用例

1. 领域包加 Request（实现 Validate）；要带额外信息就加 Response，否则直接返回模型
2. `internal/service` 实现用例：第一行 `req.Validate()`，然后取数、转换成模型、返回
3. `internal/transport/http/handler` 加方法，并在 `RegisterRoutes` 里注册路由
4. 验证：`gofmt -l .`、`go build ./...`、`go vet ./...`

## 新增一个业务模块

1. 新建 `internal/domain/<模块>/`，按需创建 entity.go、model.go、service.go、repository.go、errors.go 等文件
2. `internal/repository`、`internal/service` 各加一个文件
3. `internal/transport/http/handler` 加一个文件
4. 在 `internal/app/application.go` 里组装仓储与服务，在 `internal/transport/http/server.go` 里注册 Handler

## 不要做

以下是规则之外的工作约定：

- 不要往 cmd 里加业务代码
- 不要擅自改动通用层（`internal/foundation`、`internal/transport/http` 的公共组件）；要改先说明原因和影响
- 不要引入新的依赖或构建工具（Makefile、lint、CI），除非明确要求

## 细则索引

细则写在对应位置的注释里，改代码时一并看：

| 主题 | 位置 |
| --- | --- |
| 各层的职责与代码骨架 | 各层目录的 README.md |
| 分层与依赖方向 | 各包的包注释 |
| 领域共用的内容 | `internal/domain` |
| 领域模块的文件划分 | `internal/domain/<模块>/doc.go` 顶部 |
| Entity、Model、Request、Response 的定义与用法 | `internal/domain/<模块>/` 对应文件 |
| 分页参数与分页计算 | `internal/foundation/pagination` |
| 业务错误与错误码 | `internal/foundation/errors` |
| 日志文案约定 | `internal/foundation/logging` |
| 通用上下文数据 | `internal/foundation/metadata` |
| 业务上下文数据 | `internal/domain/metadata` |
| 请求 id | `internal/transport/http/middleware` |
| 请求绑定、错误映射、响应格式 | `internal/transport/http/jsoncodec` |
| 启动、依赖组装、优雅关闭 | `internal/app` |
| 配置结构与校验 | `internal/config` |
