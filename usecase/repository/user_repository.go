package repository

import (
	"clean_arc/domain/entity"
)

type UserRepository interface {
	SaveUser(u *entity.User) (*entity.User, error)
}
