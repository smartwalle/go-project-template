// Package service 实现业务用例。
//
// 每个业务模块一个文件，实现领域包定义的 Service 契约。
package service

import (
	"context"
	"fmt"
	"strings"
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
func (s *UserService) List(ctx context.Context) ([]*userdomain.User, error) {
	users, err := s.userRepo.FindList(ctx, "*", "id > ?", 0)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	return users, nil
}

// Get 按 Id 返回用户。用户不存在时返回 userdomain.ErrNotFound。
func (s *UserService) Get(ctx context.Context, id int64) (*userdomain.User, error) {
	if id <= 0 {
		return nil, errors.OutOfRange("id must be a positive integer")
	}

	user, err := s.userRepo.Find(ctx, id, "*")
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, userdomain.ErrNotFound
	}
	return user, nil
}

// Create 创建用户。
//
// 邮箱唯一性最终由数据库的唯一约束保证：这里的前置查询只是为了给出更友好的错误。
func (s *UserService) Create(ctx context.Context, req userdomain.CreateRequest) (*userdomain.User, error) {
	name := strings.TrimSpace(req.Name)
	email := strings.TrimSpace(req.Email)

	if name == "" {
		return nil, errors.MissingArgument("name is required")
	}
	if email == "" {
		return nil, errors.MissingArgument("email is required")
	}
	if !strings.Contains(email, "@") {
		return nil, errors.InvalidArgument(fmt.Sprintf("email %q is not a valid address", email))
	}

	existing, err := s.userRepo.FindOne(ctx, "*", "email = ?", email)
	if err != nil {
		return nil, fmt.Errorf("look up user by email: %w", err)
	}
	if existing != nil {
		return nil, errors.AlreadyExists(fmt.Sprintf("email %q", email))
	}

	// Id 留空，由数据库生成。
	entity := &userdomain.User{
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
	}
	if _, err = s.userRepo.Create(ctx, entity); err != nil {
		return nil, err
	}

	// dbs 的 Create 只返回 Result，而 PostgreSQL 驱动不支持 LastInsertId，
	// 拿不到数据库生成的 Id，这里回查一次。
	created, err := s.userRepo.FindOne(ctx, "*", "email = ?", email)
	if err != nil {
		return nil, err
	}
	if created == nil {
		return nil, userdomain.ErrNotFound
	}
	return created, nil
}
