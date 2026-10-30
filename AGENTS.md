# AGENTS.md

本文件是本仓库的编码规范，供 AI Coding Agent 和新同事直接照着执行。
这里定规则，各目录的 README 给骨架与用法，细则写在对应代码的注释里。

## 项目初始化

当用户要求基于本模板创建新项目、初始化项目结构或将模板转换为具体项目时，
必须先阅读根目录的 `PROJECT.md`，并遵循其中的初始化流程。

`PROJECT.md` 仅负责模板初始化，不替代本文档中的日常开发约定。

## 项目是什么

Go 后端服务项目结构模板：分层清晰、约定统一，提供一套业务示例帮助理解各层职责。

## 目录与分层

```text
cmd/                       进程入口：加载配置、准备环境、启动应用
internal/app/              依赖组装与生命周期
internal/config/           配置加载与校验
internal/domain/<模块>/    领域契约：Entity、Model、Request/Response、Repository/Service 接口
internal/foundation/<能力>/ 通用技术能力：errors、database、logging、metadata、pagination
internal/repository/       数据访问实现
internal/service/          业务用例实现
internal/transport/        外部入口与协议适配
internal/transport/http/   HTTP 适配层：handler、jsoncodec、中间件