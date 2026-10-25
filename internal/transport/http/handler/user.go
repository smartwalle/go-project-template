// Package handler 是业务接口的 HTTP 适配层。
//
// 这里只做协议适配：绑定请求、调用用例、输出结果。
// 业务错误到状态码的映射由 jsonrender 统一完成。
package handler

import (
	"strconv"

	"project-template/internal/transport/http/jsonrender"

	"github.com/smartwalle/rx"
	"github.com/smartwalle/rx/router"

	"project-template/internal/domain/userdomain"
	"project-template/internal/foundation/errors"
)

// UserHandler 处理用户相关的 HTTP 请求。
//
// service 是领域包定义的 userdomain.Service 契约，具体实现由应用组装层注入。
type UserHandler struct {
	service userdomain.Service
}

// NewUserHandler 创建 UserHandler。service 由应用组装层注入。
func NewUserHandler(service userdomain.Service) *UserHandler {
	return &UserHandler{service: service}
}

// RegisterRoutes 声明本 UserHandler 负责的路由。
//
// 路由前缀由 HTTP 层统一创建，Handler 只声明自己的资源路径。
func (h *UserHandler) RegisterRoutes(r router.Router) {
	r.GET("/users", jsonrender.HandleQuery(h.List))
	r.POST("/users", jsonrender.HandleJSON(h.Create))
	r.GET("/users/:id", jsonrender.HandleQuery(h.Get))
}

// List 返回用户列表。
func (h *UserHandler) List(c *rx.Context) ([]*userdomain.User, error) {
	return h.service.List(c.Context())
}

// Create 创建用户。
func (h *UserHandler) Create(c *rx.Context, req userdomain.CreateRequest) (*userdomain.User, error) {
	return h.service.Create(c.Context(), req)
}

// Get 按 Id 返回用户。
func (h *UserHandler) Get(c *rx.Context) (*userdomain.User, error) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		return nil, errors.InvalidArgument("id must be a positive integer")
	}
	return h.service.Get(c.Context(), id)
}
