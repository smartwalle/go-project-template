// Package http 是 HTTP 传输层：创建 Engine、组装路由、配置中间件，
// 并用 Server 管理 HTTP 服务的监听与启停。
//
// 统一响应格式与业务错误的映射在 jsonrender 里完成。
package http

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"project-template/internal/transport/http/jsoncodec"
	"time"

	"project-template/internal/domain/userdomain"
	"project-template/internal/foundation/errors"
	"project-template/internal/transport/http/handler"

	"github.com/smartwalle/rx"
	"github.com/smartwalle/rx/router"
)

// Server 管理 HTTP 服务：rx Engine、路由与中间件，以及 HTTP Server 的监听与启停。
//
// 它满足应用生命周期框架的服务契约（Start/Stop），由 internal/app 注册到 bootstrap。
type Server struct {
	config Config
	server *http.Server
}

// NewServer 创建 HTTP 服务。
//
// config 提供地址与超时参数，userService 是用户模块的业务契约（定义在领域包）。
func NewServer(config Config, userService userdomain.Service) *Server {
	var engine = rx.New()
	engine.Use(MidRecovery())
	engine.Use(MidRequestId())
	engine.Use(MidLogger())
	engine.UseErrorEncoder(jsoncodec.ErrorEncoder())
	engine.UseNotFoundHandler(jsoncodec.NotFoundHandler())

	// 新增业务接口时，在这里注册对应的 Handler。
	var r = router.New()
	handler.NewHealthHandler().RegisterRoutes(r)

	var api = r.Group("/api/v1")
	handler.NewUserHandler(userService).RegisterRoutes(api)

	engine.UseRouter(r)

	return &Server{
		config: config,
		server: &http.Server{
			Addr:         config.Address,
			Handler:      engine,
			ReadTimeout:  config.ReadTimeout,
			WriteTimeout: config.WriteTimeout,
			IdleTimeout:  config.IdleTimeout,
		},
	}
}

// Start 监听地址并提供服务，阻塞直到服务被关闭。
//
// 监听失败（地址非法、端口被占用）时返回错误，调用方据此判定启动失败。
func (s *Server) Start(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.config.Address)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", s.config.Address, err)
	}
	// 请求的 ctx 不跟随生命周期 ctx 取消：收到停止信号后，在途请求仍要跑完。
	s.server.BaseContext = func(net.Listener) context.Context {
		return context.WithoutCancel(ctx)
	}
	if err = s.server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("serve http: %w", err)
	}
	return nil
}

// Stop 优雅关闭：停止接收新请求，并等待在途请求结束。
//
// 等待上限由调用方给的期限决定，到点后仍未结束的请求会被中断。
func (s *Server) Stop(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			slog.Warn("http graceful shutdown timed out", "timeout", s.ShutdownTimeout())
		}
		return fmt.Errorf("shutdown http server: %w", err)
	}
	return nil
}

// ShutdownTimeout 返回优雅关闭的等待上限，也就是排空在途请求的最长时间。
func (s *Server) ShutdownTimeout() time.Duration {
	if s.config.ShutdownTimeout > 0 {
		return s.config.ShutdownTimeout
	}
	return defaultShutdownTimeout
}
