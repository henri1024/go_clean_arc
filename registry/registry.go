package registry

import (
	"clean_arc/interface/controller"

	"github.com/jinzhu/gorm"
)

type registry struct {
	db *gorm.DB
}

type Registry interface {
	NewAppController() controller.AppController
}

func NewRegistry(db *gorm.DB) Registry {
	return &registry{
		db: db,
	}
}

func (r *registry) NewAppController() controller.AppController {
	return r.NewUserController()
}
