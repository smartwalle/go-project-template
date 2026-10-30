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

```
cmd/                       进程入口
internal/app/              依赖组装与生命周期
internal/config/           配置加载与校验
internal/domain/<模块>/    领域契约：Entity、Model、Request/Response、Repository/Service 接口
internal/foundation/<能力>/ 平台能力：errors、database、logging、metadata、pagination
internal/repository/       数据访问实现
internal/service/          业务用例实现
internal/transport/http/   HTTP 适配层
```

依赖方向：`transport/http → service → domain`，`foundation` 被所有层依赖，`app` 负责组装。

## 一次请求的完整链路

以 `POST /api/v1/users/list` 为例：

| 环节 | 位置 | 负责什么 |
| --- | --- | --- |
| 路由 | `internal/transport/http/server.go` | 注册中间件、路由分组 |
| 绑定 | `internal/transport/http/jsoncodec` | 把 JSON 请求体绑成 Request；错误统一编码成响应 |
| 适配 | `internal/transport/http/handler` | 取出参数、调用用例、返回结果 |
| 校验 | `internal/domain/<模块>` 的 Request.Validate | 校验请求参数，并把清洗结果写回字段 |
| 用例 | `internal/service` | 业务逻辑：取数、转换成模型、组装响应 |
| 数据 | `internal/repository` → dbs | 执行 SQL，错误原样返回 |
| 响应 | `internal/transport/http/jsoncodec` | 统一响应体 `{code, message, data}` |

每一层只做自己的事：handler 不写业务逻辑，service 不碰 HTTP 细节，domain 不依赖任何实现包。

## 文档

| 文档 | 内容 |
| --- | --- |
| [AGENTS.md](AGENTS.md) | 编码规范：分层、命名、行为约定、HTTP 接口约定、改动清单 |
| [PROJECT.md](PROJECT.md) | 项目初始化：如何基于模板创建具体项目 |
| [internal/domain/README.md](internal/domain/README.md) | 领域契约的六种类型与写法 |
| [internal/service/README.md](internal/service/README.md) | 业务用例的写法、查询放法、事务 |
| [internal/repository/README.md](internal/repository/README.md) | 数据访问实现 |
| [internal/transport/http/README.md](internal/transport/http/README.md) | HTTP 适配层、路由与 handler 写法 |
| [internal/foundation/README.md](internal/foundation/README.md) | 平台能力：错误、分页、日志、请求 id |
| [internal/app/README.md](internal/app/README.md) | 依赖组装与生命周期 |
| [internal/config/README.md](internal/config/README.md) | 配置结构与加载 |
| [cmd/README.md](cmd/README.md) | 进程入口 |

细则写在对应代码的注释里，README 给出骨架，AGENTS.md 列出规则。
