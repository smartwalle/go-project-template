package domain

import (
	"github.com/smartwalle/dbs"
	"project-template/internal/domain/entity"
)

type UserRepository interface {
	dbs.Repositoy[entity.User]
}
