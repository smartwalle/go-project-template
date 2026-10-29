// Package handler 是业务接口的 HTTP 适配层。
//
// 这里只做协议适配：绑定请求、调用用例、输出结果。
// 业务错误到状态码的映射由 jsoncodec 统一完成。
//
// 接口命名规则：
//   - 全部用 POST，参数一律放在 JSON 请求体；
//   - 路径是"资源 + 用例名"，例如：
//     列表 POST /users/list，创建 POST /users/create，
//     单个查询 POST /users/info，更新 POST /users/update。
package handler

import (
	"project-template/internal/transport/http/jsoncodec"

	"github.com/smartwalle/rx"
	"github.com/smartwalle/rx/router"

	"project-template/internal/domain/userdomain"
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
	r.POST("/users/list", jsoncodec.Handler(h.List))
	r.POST("/users/info", jsoncodec.Handler(h.Get))
	r.POST("/users/create", jsoncodec.Handler(h.Create))
	r.POST("/users/update", jsoncodec.Handler(h.Update))
}

// List 返回用户列表。
//
// 分页参数和过滤条件都在请求体里：page 是页码，page_size 是每页条数，
// name 和 email 是模糊匹配条件。
func (h *UserHandler) List(c *rx.Context, req *userdomain.ListRequest) (*userdomain.ListResponse, error) {
	return h.service.List(c.Context(), req)
}

// Get 按 Id 返回用户。
func (h *UserHandler) Get(c *rx.Context, req *userdomain.GetRequest) (*userdomain.User, error) {
	return h.service.Get(c.Context(), req)
}

// Create 创建用户。
func (h *UserHandler) Create(c *rx.Context, req *userdomain.CreateRequest) (*userdomain.User, error) {
	return h.service.Create(c.Context(), req)
}

// Update 更新用户信息。
func (h *UserHandler) Update(c *rx.Context, req *userdomain.UpdateRequest) (*userdomain.User, error) {
	return h.service.Update(c.Context(), req)
}
