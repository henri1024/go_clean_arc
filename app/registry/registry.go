package registry

import (
	"userauth/user/controller"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
)

type registry struct {
	db  *gorm.DB
	rdb *redis.Client
}

type Registry interface {
	NewUserController() controller.UserController
}

func NewRegistry(db *gorm.DB, rdb *redis.Client) Registry {
	return &registry{
		db:  db,
		rdb: rdb,
	}
}
