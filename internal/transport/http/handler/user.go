// Package handler 是业务接口的 HTTP 适配层。
//
// 这里只做协议适配：绑定请求、调用用例、输出结果。
// 业务错误到状态码的映射由 jsonrender 统一完成。
package handler

import (
	"strconv"

	"project-template/internal/transport/http/jsoncodec"

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
	r.GET("/users", jsoncodec.Handler(h.List))
	r.POST("/users", jsoncodec.Handler(h.Create))
	r.GET("/users/:id", jsoncodec.Handler(h.Get))
	r.PUT("/users/:id", jsoncodec.Handler(h.Update))
}

// List 返回用户列表。
func (h *UserHandler) List(c *rx.Context, _ struct{}) (*userdomain.ListResponse, error) {
	return h.service.List(c.Context(), &userdomain.ListRequest{})
}

// Create 创建用户。
func (h *UserHandler) Create(c *rx.Context, req *userdomain.CreateRequest) (*userdomain.User, error) {
	return h.service.Create(c.Context(), req)
}

// Get 按 Id 返回用户。
func (h *UserHandler) Get(c *rx.Context, _ struct{}) (*userdomain.User, error) {
	id, err := parseId(c)
	if err != nil {
		return nil, err
	}
	return h.service.Get(c.Context(), id)
}

// Update 更新用户信息。
func (h *UserHandler) Update(c *rx.Context, req *userdomain.UpdateRequest) (*userdomain.User, error) {
	id, err := parseId(c)
	if err != nil {
		return nil, err
	}
	if req == nil {
		return nil, errors.MissingArgument("request is required")
	}
	// id 只来自路径参数，请求体里不带。
	req.Id = id
	return h.service.Update(c.Context(), req)
}

// parseId 解析路径里的用户标识。
func parseId(c *rx.Context) (uint64, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		return 0, errors.InvalidArgument("id must be a positive integer")
	}
	return id, nil
}
