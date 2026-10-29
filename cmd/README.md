# cmd

进程入口。

**这个目录基本是固定的**：内容只有加载配置、准备运行环境、启动服务三件事，
新增业务模块、接口、配置项都不需要动这里的代码——
业务接在哪几个目录，见各目录的 README。

| 文件 | 职责 |
| --- | --- |
| `main.go` | 加载配置 → 准备运行环境 → 启动服务 → 优雅关闭 |
| `config.example.yaml` | 配置样例，复制成 `config.yaml` 再改（`config.yaml` 不入库） |

会变的只有配置文件本身：新增配置项时改 `config.yaml`（本地）和 `config.example.yaml`（样例）。

```go
func main() {
	if err := run(); err != nil {
		slog.Error("应用运行失败", "error", err)
		os.Exit(1)
	}
}

// run 加载配置、准备运行环境并启动服务。
//
// 环境在这里创建并释放
func run() error {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		slog.Error("配置加载失败", slog.Any("error", err))
		return fmt.Errorf("load config: %w", err)
	}

	env, err := app.Prepare(cfg)
	if err != nil {
		slog.Error("运行环境准备失败", slog.Any("error", err))
		return fmt.Errorf("prepare environment: %w", err)
	}
	defer env.Clear()

	if err = app.Run(context.Background(), env); err != nil {
		slog.Error("应用启动失败", slog.Any("error", err))
		return err
	}
	return nil
}
```

- 入口只做编排，不写业务逻辑
- 资源在 `run` 里创建，退出前 `defer env.Clear()` 释放
- 配置文件路径由这里传给 `config.Load`
- 只有进程级的事（启动参数、退出处理）才需要改 `main.go`
