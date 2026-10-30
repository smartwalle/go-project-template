# AGENTS.md

本文件是本仓库的统一开发规范，供 AI Coding Agent 和开发人员遵循。

本文档规定跨目录、跨层的开发约定。各目录的 README 负责说明本层职责、代码组织与使用方式，具体实现细则以对应代码中的注释为准。

## 项目初始化

当用户要求基于本模板创建新项目、初始化项目结构或将模板转换为具体项目时，必须先阅读根目录的 `PROJECT.md`，并遵循其中的初始化流程。

`PROJECT.md` 负责模板初始化，不替代本文档中的日常开发约定。

## 项目定位

本项目是 Go 后端服务项目模板，提供清晰的分层结构、统一的编码约定和基础业务示例。

## 目录与分层

```text
cmd/                       进程入口：加载配置、准备环境、启动应用
internal/app/              依赖组装与应用生命周期
internal/config/           配置加载与校验
internal/domain/<模块>/    领域契约：Entity、Model、Request/Response、Repository/Service 接口
internal/foundation/<能力>/ 通用技术能力：errors、database、logging、metadata、pagination
internal/repository/       数据访问实现
internal/service/          业务用例实现
internal/transport/        外部入口与协议适配
internal/transport/http/   HTTP 适配层：handler、jsoncodec、中间件
```

### 分层职责

- `cmd`：进程入口，负责加载配置、准备应用环境并启动应用，不包含业务逻辑。
- `app`：负责依赖创建、组装、注入和应用生命周期管理，不包含业务逻辑。
- `config`：负责配置结构定义、配置加载与校验。
- `domain`：定义领域模型和业务契约，不负责具体实现。
- `foundation`：提供可跨业务复用的通用技术能力，不依赖具体业务领域，不包含业务专属模型、规则和实现。
- `repository`：实现领域包中定义的 Repository 契约，负责数据访问。
- `service`：实现领域包中定义的 Service 契约，负责业务用例。
- `transport`：负责外部入口与协议适配，例如 HTTP、消息消费、定时任务等。根据实际需求组织，不要求每个业务用例都必须提供 HTTP 接口。

### 依赖方向

- `domain` 不依赖 `transport`、`service`、`repository` 等实现层，也不依赖具体数据库驱动。
- `service` 依赖 `domain` 中定义的业务契约和数据访问契约，不依赖 `repository` 的具体实现。
- `repository` 依赖 `domain` 中定义的实体和 Repository 契约，不依赖 `service`。
- `transport` 调用 `service` 提供的业务用例，不直接访问数据库或实现业务规则。
- `foundation` 不依赖具体业务领域，可以被需要相关通用能力的层按需依赖。
- `app` 负责组装各层实现，不承担业务逻辑。
- 各层之间通过接口隔离具体实现：接口定义在 `domain`，实现分别放在 `service`、`repository` 等对应实现层。
- 不需要某项基础能力的包，不应为了遵循分层形式而引入对应依赖。
- 全局共享资源（配置、数据库连接等）由 `app.Environment` 管理：`Prepare` 创建和初始化，`Clear` 释放资源。

## 命名约定

| 类型 | 约定 | 位置 |
| --- | --- | --- |
| Entity | 每张表一个 XXXEntity，名称为表名驼峰形式加 Entity 后缀 | `internal/domain/<模块>/entity.go` |
| Model | 业务场景所需的数据结构，不绑定数据库表 | `internal/domain/<模块>/model.go` |
| Request | XXXRequest，XXX 为用例名 | `internal/domain/<模块>/service.go` |
| Response | XXXResponse，XXX 为用例名 | `internal/domain/<模块>/service.go` |
| Repository | 每张表一个仓储接口 | `internal/domain/<模块>/repository.go` |
| Service | 定义业务用例接口 | `internal/domain/<模块>/service.go` |

### Entity 与 Model

- Entity 是数据库表的映射结构，一个 Entity 对应一张表。
- Entity 使用 SQL tag 映射数据库列，并按需实现 `TableName`、`PrimaryKey` 等数据库访问所需的方法。
- Model 表示业务场景所需的数据结构，不绑定数据库表，主要用于整理和表达业务数据。
- Entity 不直接作为业务接口的返回对象。需要返回业务数据时，应转换为对应的 Model。
- Model 可以根据实际业务需要组织数据，不要求与 Entity 一一对应，也不要求一定由多个 Entity 组合而成。
- 不要仅仅为了区分类型而创建没有实际意义的 Model。

### Request 与 Response

- Request 表示一个业务用例的请求参数。
- Response 用于表达需要额外组织或返回多项信息的结果。
- 如果返回结构本身已经能够清楚表达结果，可以直接返回 Model，不必额外定义 Response。
- Request、Response 定义在领域包中，不放在 HTTP Handler 或 Service 实现包中。
- Request 统一使用指针传递，并实现 `Validate() error`。

### 领域共用内容

- 所有领域共用的定义放在 `internal/domain`（`package domain`）。
- 只在一个领域中使用的定义放在对应领域包中。
- 不要仅仅因为某个类型看起来通用，就将它放到领域根包；是否共用应以实际使用范围为依据。

### 标识符（ID）

- Entity、Model 中的标识符类型统一使用 `uint64`。
- Request 中接收的标识符可以使用 `int64`，以便校验负数。
- 有效性判断统一使用 `> 0` / `<= 0`，不要只用 `== 0`，避免有符号整型出现负数时漏检。

## 行为约定

### 请求参数校验

- 每个 Request 必须实现 `Validate() error`。
- `Validate` 负责 nil 判断、参数校验和数据归一化。
- 数据归一化是指对请求数据进行必要清洗，并将清洗结果写回字段。
- Service 用例的入口处首先调用 `req.Validate()`，校验失败直接返回错误。
- 校验通过后，Service 直接使用归一化后的请求字段，不重复执行相同校验。
- 不需要请求参数的用例，不必为了形式创建空 Request。

### 上下文数据

上下文数据按通用性区分：

- 业务相关的数据（当前登录用户、租户等）统一放在 `internal/domain/metadata` 管理，具体数据类型留在各自领域包。
- 通用的数据（请求 ID 等）统一放在 `internal/foundation/metadata` 管理。

约定：

- 每份数据使用私有空结构体作为 context key，避免与其他包的 key 冲突。
- 写入方法使用 `WithXxx` 命名，读取方法使用 `XxxFrom` 命名。
- 读取不到时返回零值，不返回错误。
- 不要直接使用字符串等基础类型作为 context key。

### 错误处理

- 数据访问（Repository）返回的错误原样透传，不随意包装或转换。
- 用例自己产生的通用错误使用 `internal/foundation/errors` 中定义的错误，例如 `NotFound`、`AlreadyExists`、`InvalidArgument`、`MissingArgument` 等。
- 业务专属的错误语义在对应领域包中定义，使用带业务前缀的错误 code，例如 `user_not_found`。
- 判断错误语义按 code 比较，不按错误值比较。错误 message 可以随业务变化，调用方应通过错误类型提供的 `Code()` 获取 code 后判断。
- 不要为了包装错误而重复包装；需要增加上下文时，应确保不会改变调用方判断错误语义的方式。

### 数据查询

- 会被复用的查询放进 Repository 契约。
- 一次性的查询可以在 Service 中使用 `SelectBuilder` 拼接。
- SQL 条件一律使用占位符传参，不要把用户输入直接拼进 SQL。
- 存在性判断：查不到（nil）或主键 `<= 0`，都视为没有这条记录。
- 逻辑删除的表，查询时必须带上未删除条件。
- 查询条件、排序字段等动态 SQL 片段必须来自受控定义，不得直接使用未经校验的用户输入。

### 事务

- 多个写入操作需要一起成功或一起失败时，放进同一个事务。
- 使用 Repository 提供的 `Transaction` 方法。
- 事务回调中必须使用回调收到的 `txCtx` 执行事务内的数据库操作，不能使用外层的 `ctx`，否则操作可能运行在事务之外。
- 回调返回错误即回滚，返回 nil 即提交。
- 不要为了单条普通写入无意义地增加事务。

事务调用示例：

```go
err := repo.Transaction(ctx, func(txCtx context.Context) error {
// 使用 txCtx 执行事务内的数据库操作。
return nil
}, nil)
```

### 变量

- 不要和外层变量重名，避免变量遮蔽。
- 尤其注意事务回调里的 `ctx`、`err` 等变量，应使用清晰且不同的名称。
- 优先使用含义明确的变量名，避免无意义的缩写。

### 注释

- 使用中文，以标识符开头，句末使用句号。
- 需要解释原因或背景时另起一段。
- 只保留最终结论，不写讨论过程和方案取舍。
- 注释应说明代码的职责、约束或必要背景，不重复代码本身已经表达的内容。
- 导出类型、函数、方法应按照 Go 文档注释规范编写注释。

## HTTP 接口约定

以下约定适用于项目内的业务 HTTP 接口，不限制其他 Transport 的协议设计。

### 请求与路径

- 业务接口统一使用 POST，参数一律放在 JSON 请求体。
- 健康检查等基础设施接口可以根据接口语义使用合适的 HTTP 方法。
- 业务接口请求体必须存在，空请求体按参数不合法处理。
- 路径根据资源和接口语义确定，不要求路径动作名称与 Service 方法名完全一致。

例如：

| 业务语义 | HTTP 路径 | Service 方法 |
| --- | --- | --- |
| 列表 | `POST /users/list` | `List` |
| 创建 | `POST /users/create` | `Create` |
| 查询 | `POST /users/info` | `Get` |
| 更新 | `POST /users/update` | `Update` |

### Handler

- Handler 只做协议适配：绑定请求、调用业务用例、返回结果。
- 业务逻辑一律放在 Service，不要在 Handler 中实现业务规则。
- Handler 方法签名遵循项目统一约定，接收 HTTP Context 和对应的 Request，返回业务结果及错误。
- 业务错误的响应由 `jsoncodec` 统一处理，Handler 不自行定义通用业务错误的 HTTP 响应格式。
- 业务接口统一使用 `jsoncodec` 绑定请求、输出响应。
- 对接第三方且报文格式不同时，另建对应的编解码实现，不要修改业务通用的 `jsoncodec`。
- 返回结构有特殊处理需求时，在对应 Handler 中明确处理，不要为了复用通用响应逻辑而滥用 `AbortWithError`。
- 路由在 Handler 的 `RegisterRoutes` 中注册，路径前缀由 `transport/http` 统一添加。

## 日志

- 日志消息（msg）使用中文，字段名（key）保持英文。
- 字段名应保持稳定，不要随意修改，以免影响日志检索和统计。
- SQL 语句、配置路径、错误码、第三方库名等机器可读内容保持原样，不做翻译。
- 错误文案原样透传，不要随意改写。
- 请求 ID 通过 `foundation/metadata` 在 context 中传递，由日志组件按项目约定注入。
- SQL 日志不得记录密码、邮箱、手机号等敏感参数的原始值。

## 新增一个用例

1. 在领域包中定义 Request，并实现 `Validate()`。
2. 根据返回结果的实际需要定义 Model 或 Response；不需要额外返回信息时，直接返回 Model。
3. 在 `internal/service` 中实现领域包定义的 Service 契约。
4. Service 用例入口处首先调用 `req.Validate()`，然后完成业务逻辑、数据访问和模型转换。
5. 如果需要对外提供 HTTP 接口，在 `internal/transport/http/handler` 中增加 Handler 方法，并在 `RegisterRoutes` 中注册路由。
6. 如果需要通过其他协议或任务入口调用，放在对应的 Transport 实现中，不要为了复用业务逻辑而把协议处理放进 Service。
7. 不需要请求参数的用例，不必创建空 Request；按实际契约设计方法签名。

## 新增一个业务模块

1. 新建 `internal/domain/<模块>/`，按需创建 `entity.go`、`model.go`、`service.go`、`repository.go`、`errors.go` 等文件。
2. 在 `internal/repository` 中增加对应的 Repository 实现。
3. 在 `internal/service` 中增加对应的 Service 实现。
4. 根据实际需要，在 `internal/transport` 下增加对应的入口或协议适配实现。
5. 在 `internal/app` 中组装 Repository 与 Service，并将实现注入对应的使用方。
6. 不需要对外提供 HTTP 接口的业务模块，不要为了目录形式额外创建 HTTP Handler。
7. 不要为了满足目录结构创建没有实际用途的空文件、空接口或抽象。

## 不要做

- 不要往 `cmd` 中添加业务代码。
- 不要把业务逻辑放进 `app`、`repository` 或 `transport`。
- 不要让 `domain` 依赖具体实现层或外部传输协议。
- 不要让 `foundation` 依赖具体业务领域。
- 不要让 `service` 直接依赖 Repository 的具体实现。
- 不要擅自修改通用层（`internal/foundation`、`internal/transport` 的公共组件）；需要修改时，先说明原因和影响。
- 不要引入新的依赖或构建工具（Makefile、lint、CI 等），除非用户明确要求。
- 不要为了满足目录形式创建没有实际意义的抽象。
- 不要把协议处理逻辑放进 Service，也不要把业务规则放进 Handler。

## 细则索引

具体实现细则写在对应位置的 README 和代码注释中，修改代码时应一并阅读：

| 主题 | 位置 |
| --- | --- |
| 项目初始化流程 | `PROJECT.md` |
| 各层的职责与代码骨架 | 各层目录的 `README.md` |
| 分层与依赖方向 | 各包的包注释 |
| 领域共用的内容 | `internal/domain` |
| 领域模块的文件划分 | `internal/domain/<模块>/doc.go` |
| Entity、Model、Request、Response 的定义与用法 | `internal/domain/<模块>/` 对应文件 |
| 分页参数与分页计算 | `internal/foundation/pagination` |
| 通用业务错误与错误码 | `internal/foundation/errors` |
| 日志约定 | `internal/foundation/logging` |
| 通用上下文数据 | `internal/foundation/metadata` |
| 业务上下文数据 | `internal/domain/metadata` |
| 请求 ID | `internal/transport/http/middleware` |
| 请求绑定、错误映射、响应格式 | `internal/transport/http/jsoncodec` |
| 启动、依赖组装、优雅关闭 | `internal/app` |
| 配置结构与校验 | `internal/config` |
| 进程入口 | `cmd` |