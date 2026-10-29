// Package service 实现业务用例。
//
// 每个业务模块一个文件，实现领域包定义的 Service 契约。
//
// 错误约定：数据访问（仓储）返回的错误原样透传，不做包装、不做转换，
// 保证调用方能按原始错误判断（比如数据库的唯一约束冲突）；
// 用例自己产生的错误用 foundation/errors 里的业务错误。
package service

import (
	"context"
	"fmt"
	"time"

	"project-template/internal/domain/userdomain"
	"project-template/internal/foundation/errors"
)

// UserService 实现用户模块的业务用例。
type UserService struct {
	userRepo userdomain.Repository
}

// NewUserService 创建 UserService。repo 由调用方（应用组装层）注入具体实现。
func NewUserService(repo userdomain.Repository) *UserService {
	return &UserService{userRepo: repo}
}

// List 返回全部用户。
func (s *UserService) List(ctx context.Context, req *userdomain.ListRequest) (*userdomain.ListResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	entityList, err := s.userRepo.FindList(ctx, "*", "id > ?", 0)
	if err != nil {
		return nil, err
	}

	var response = &userdomain.ListResponse{
		Total:    len(entityList),
		UserList: toModels(entityList),
	}
	return response, nil
}

// Get 按 Id 返回用户。用户不存在时返回 userdomain.ErrNotFound。
func (s *UserService) Get(ctx context.Context, id uint64) (*userdomain.User, error) {
	if id == 0 {
		return nil, errors.OutOfRange("id must be a positive integer")
	}

	entity, err := s.userRepo.Find(ctx, id, "*")
	if err != nil {
		return nil, err
	}
	if entity == nil || entity.Id <= 0 {
		return nil, userdomain.ErrNotFound
	}
	return toModel(entity), nil
}

// Create 创建用户。
//
// 参数校验由 Request.Validate 完成，这里在用例入口处调用它；
// 校验通过后请求里的字段已经是规范化后的值，直接使用。
//
// 邮箱唯一性最终由数据库的唯一约束保证：这里的前置查询只是为了给出更友好的错误。
func (s *UserService) Create(ctx context.Context, req *userdomain.CreateRequest) (*userdomain.User, error) {
	// Validate 保证了请求非空，校验通过后直接用请求里的字段。
	if err := req.Validate(); err != nil {
		return nil, err
	}

	existing, err := s.userRepo.FindOne(ctx, "*", "email = ?", req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil && existing.Id > 0 {
		return nil, errors.AlreadyExists(fmt.Sprintf("email %q", req.Email))
	}

	// Id 留空，由数据库生成。
	entity := &userdomain.UserEntity{
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: time.Now(),
	}
	if _, err = s.userRepo.Create(ctx, entity); err != nil {
		return nil, err
	}

	// dbs 的 Create 只返回 Result，而 PostgreSQL 驱动不支持 LastInsertId，
	// 拿不到数据库生成的 Id，这里回查一次。
	created, err := s.userRepo.FindOne(ctx, "*", "email = ?", req.Email)
	if err != nil {
		return nil, err
	}
	if created == nil || created.Id <= 0 {
		return nil, userdomain.ErrNotFound
	}
	return toModel(created), nil
}

// Update 更新用户信息。
//
// 这里的更新放在事务里执行，只是为了演示事务的用法：
// 当前只有一条更新语句，本身用不到事务；
// 以后要同时改多张表时，把操作加进同一个事务即可，不用改结构。
func (s *UserService) Update(ctx context.Context, req *userdomain.UpdateRequest) (*userdomain.User, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	var id = req.Id

	existing, err := s.userRepo.Find(ctx, id, "*")
	if err != nil {
		return nil, err
	}
	if existing == nil || existing.Id <= 0 {
		return nil, userdomain.ErrNotFound
	}

	// 邮箱唯一性：排除自己，只看有没有被别的用户占用。
	conflict, err := s.userRepo.FindOne(ctx, "*", "email = ? AND id <> ?", req.Email, id)
	if err != nil {
		return nil, err
	}
	if conflict != nil && conflict.Id > 0 {
		return nil, errors.AlreadyExists(fmt.Sprintf("email %q", req.Email))
	}

	// 事务示例：这里的更新放在事务里，只是为了演示事务的用法。
	// Transaction 会把事务对象放进回调收到的 ctx 里，
	// 回调里要用这个 ctx 继续操作，语句才落在同一个事务上，
	// 所以这里用 txCtx 命名，不和外层的 ctx 混用。
	if err = s.userRepo.Transaction(ctx, func(txCtx context.Context) error {
		if _, updateErr := s.userRepo.Update(txCtx, id, map[string]any{
			"name":  req.Name,
			"email": req.Email,
		}); updateErr != nil {
			return updateErr
		}
		return nil
	}, nil); err != nil {
		return nil, err
	}

	updated, err := s.userRepo.Find(ctx, id, "*")
	if err != nil {
		return nil, err
	}
	if updated == nil || updated.Id <= 0 {
		return nil, userdomain.ErrNotFound
	}
	return toModel(updated), nil
}

// toModel 把实体转换成对外模型。
//
// Entity 对应数据库结构，不直接对外输出，只在这里转换成模型。
func toModel(entity *userdomain.UserEntity) *userdomain.User {
	if entity == nil {
		return nil
	}
	return &userdomain.User{
		Id:        entity.Id,
		Name:      entity.Name,
		Email:     entity.Email,
		CreatedAt: entity.CreatedAt,
	}
}

// toModels 批量把实体转换成对外模型。
func toModels(entityList []*userdomain.UserEntity) []*userdomain.User {
	if entityList == nil {
		return nil
	}
	var models = make([]*userdomain.User, 0, len(entityList))
	for _, entity := range entityList {
		models = append(models, toModel(entity))
	}
	return models
}
