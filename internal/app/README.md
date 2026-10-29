# internal/app

依赖组装与生命周期：把各层实现接到一起，交给 bootstrap 管理启动与关闭。

| 文件 | 职责 |
| --- | --- |
| `application.go` | 组装依赖、注册到 bootstrap，管理应用的运行生命周期 |
| `environment.go` | 持有全局资源：创建和释放运行期资源 |

## Run：组装与生命周期

`Run` 是应用运行的核心入口，做两件事：**按依赖方向创建实现、把实现注入给使用方**，
再把服务交给 bootstrap 管理启动与关闭。

它阻塞运行，直到收到停止信号、上下文取消或启动失败。

```go
// Run 组装依赖并启动服务，阻塞直到收到停止信号、上下文取消或启动失败。
//
// Environment 持有数据库等资源，由它负责创建与释放；
// 这里只负责把实现组装起来，并交给 bootstrap 管理生命周期。
func Run(ctx context.Context, env *Environment) error {
	if env == nil {
		return errors.New("environment is required")
	}
	// 新增业务模块时，在这里接入它的仓储与服务；Handler 在 transport/http 里注册。
	var userRepo = repository.NewUserRepository(env.Database)
	var userService = service.NewUserService(userRepo)

	// HTTP 服务的监听与启停由传输层的 Server 管理，应用层只负责把它注册到生命周期框架。
	var server = http.NewServer(env.Config.HTTP, userService)

	var application = bootstrap.New(
		bootstrap.WithContext(ctx),
		bootstrap.WithServers(server),
		bootstrap.WithStopTimeout(env.Config.App.ShutdownTimeout),
	)

	if err := application.Run(); err != nil {
		return fmt.Errorf("run application: %w", err)
	}
	return nil
}
```

### 新增业务模块在这里接入

三步，顺序不能颠倒：先仓储、再服务、最后把服务交给使用方。

1. 创建仓储，注入 `env.Database`；
2. 用仓储创建服务；
3. 把服务实例传给使用方（HTTP 的话是 `http.NewServer`）。

路由不在这里注册——Handler 在 `internal/transport/http` 里注册，
`Run` 只负责把服务实例递过去。

## Environment：全局资源都在这里

`Environment` 是整个应用的运行环境，**全局共享的资源都挂在它上面**——
配置、数据库连接，以及以后要加的缓存、消息队列客户端之类的。

```go
// Environment 是应用运行期需要的资源，由它负责创建与释放。
type Environment struct {
	// Config 是应用的运行时配置。
	Config *config.Config

	// Database 是数据库连接池。
	Database *dbs.DB
}
```

资源的生命周期也在这里闭环：

```go
// Prepare 校验配置并初始化应用依赖（日志、数据库）。
//
// 任何一步失败都会返回错误，已经创建的资源会被释放。
func Prepare(cfg *config.Config) (*Environment, error) {
	if cfg == nil {
		return nil, errors.New("config is required")
	}

	var env = &Environment{}
	if err := env.setup(cfg); err != nil {
		return nil, err
	}
	return env, nil
}

// Clear 释放应用占用的资源。
func (env *Environment) Clear() {
	if env.Database != nil {
		if err := env.Database.Close(); err != nil {
			slog.Error("数据库关闭异常", slog.Any("error", err))
		}
	}
}
```

**新增一个全局资源时，三处都要动**：

1. 在 `Environment` 上加字段；
2. 在 `setup` 里加一个 `initXxx` 完成创建，失败就返回错误；
3. 在 `Clear` 里释放它。

## 共同约定

- 应用层不写业务逻辑
- 全局资源由 `Environment` 持有：`Prepare` 创建，`Clear` 释放；入口那边 `defer env.Clear()`
- 业务模块不自己管连接，需要什么从 `Environment` 拿
- 服务的启停交给 bootstrap，`Run` 只负责注册
