package domain

import (
	"context"
)

type User struct {
	ID   int64  `sql:"id"`
	Name string `sql:"name"`
}

type UserRepository interface {
	Create(ctx context.Context, user *User) (err error)

	FindByID(ctx context.Context, id int64) (user *User, err error)
}
