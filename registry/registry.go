package registry

import (
	"clean_arc/interface/controller"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
)

type registry struct {
	db      *gorm.DB
	redisDB *redis.Client
}

type Registry interface {
	NewAppController() controller.AppController
}

func NewRegistry(db *gorm.DB, rdb *redis.Client) Registry {
	return &registry{
		db:      db,
		redisDB: rdb,
	}
}

func (r *registry) NewAppController() controller.AppController {
	return r.NewUserController()
}
