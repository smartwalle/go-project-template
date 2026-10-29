package app

import (
	"context"
	"errors"
	"fmt"

	"project-template/internal/repository"
	"project-template/internal/service"

	"github.com/smartwalle/bootstrap"

	"project-template/internal/transport/http"
)

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
