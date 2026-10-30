# go-project-template

Go 后端服务模板：分层清晰、约定统一，业务模块自带一套完整示例。

## 快速开始

```bash
cp cmd/config.example.yaml cmd/config.yaml   # 改成自己的数据库连接
cd cmd && go run .
```

配置在 `cmd/config.yaml`，该文件不入库；`cmd/config.example.yaml` 是样例。

服务默认监听 `127.0.0.1:8080`。

```bash
curl -X POST 'http://127.0.0.1:8080/api/v1/users/list' \
  -H 'Content-Type: application/json' \
  -d '{"page":1,"page_size":10}'
```

## 分层

```text
cmd/                       进程入口
internal/app/              依赖组装与生命周期
internal/config/           配置加载与校验
internal/domain/<模块>/    领域契约：Entity、Model、Request/Response、Repository/Service 接口
internal/foundation/<能力>/ 通用技术能力：errors、database、logging、metadata、pagination
internal/repository/       数据访问实现
internal/service/          业务用例实现
internal/transport/        外部入口与协议适配
internal/transport/http/   HTTP 适配层
```

依赖约定：

- `transport` 调用 `service` 提供的业务用例。
- `service` 实现 `domain` 中定义的 Service 契约，并通过 Repository 契约访问数据，不依赖 Repository 的具体实现。
- `repository` 实现 `domain` 中定义的 Repository 契约。
- `foundation` 提供跨业务通用技术能力，由需要相关能力的层按需依赖。
- `app` 负责创建、组装各层实现，并管理应用生命周期。
- `domain` 定义业务模型和契约，不依赖具体实现层和外部传输协议。

各层职责和依赖约定详见 [AGENTS.md](AGENTS.md)。

## 一次请求的完整链路

以 `POST /api/v1/users/list` 为例：

| 环节 | 位置 | 负责什么 |
| --- | --- | --- |
| 路由 | `internal/transport/http/server.go` | 注册中间件、路由分组 |
| 绑定 | `internal/transport/http/jsoncodec` | 把 JSON 请求体绑定成 Request；错误统一编码成响应 |
| 适配 | `internal/transport/http/handler` | 取出参数、调用用例、返回结果 |
| 校验 | `internal/domain/<模块>` 的 Request.Validate | 校验请求参数，并把清洗结果写回字段 |
| 用例 | `internal/service` | 业务逻辑：取数、转换成模型、组装响应 |
| 数据 | `internal/repository` → dbs | 执行 SQL，错误原样返回 |
| 响应 | `internal/transport/http/jsoncodec` | 输出统一响应；成功时 `code`、`message` 为空，`data` 按 `omitempty` 规则输出；失败时输出错误 `code`、`message`，不输出 `data` |

每一层只做自己的事：Handler 不写业务逻辑，Service 不处理 HTTP 细节，Domain 不依赖具体实现包。

## 文档

| 文档 | 内容 |
| --- | --- |
| [AGENTS.md](AGENTS.md) | 统一开发规范：分层、依赖方向、命名、行为约定和开发流程 |
| [PROJECT.md](PROJECT.md) | 项目初始化：如何基于模板创建具体项目 |
| [internal/domain/README.md](internal/domain/README.md) | 领域契约的类型、职责与写法 |
| [internal/service/README.md](internal/service/README.md) | 业务用例的写法、查询放法、事务 |
| [internal/repository/README.md](internal/repository/README.md) | 数据访问实现 |
| [internal/transport/http/README.md](internal/transport/http/README.md) | HTTP 适配层、路由与 Handler 写法 |
| [internal/foundation/README.md](internal/foundation/README.md) | 通用技术能力：错误、分页、日志、请求 ID |
| [internal/app/README.md](internal/app/README.md) | 依赖组装与生命周期 |
| [internal/config/README.md](internal/config/README.md) | 配置结构与加载 |
| [cmd/README.md](cmd/README.md) | 进程入口 |

文档职责：

- `AGENTS.md` 维护跨目录、跨层的统一开发规范。
- 各层 README 说明本层职责、代码组织和使用方式。
- `doc.go` 与代码注释说明具体包和实现的细节。

避免在多份文档中重复维护同一条规则；调整规范时，应同步检查受影响的文档和代码注释。