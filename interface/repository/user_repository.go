package repository

import (
	"clean_arc/domain/entity"

	"github.com/jinzhu/gorm"
)

type userRepository struct {
	db *gorm.DB
}

type UserRepository interface {
	SaveUser(u *entity.User) (*entity.User, error)
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (ur *userRepository) SaveUser(user *entity.User) (*entity.User, error) {
	err := ur.db.Model(&entity.User{}).Save(user).Error

	if err != nil {
		return nil, err
	}
	return user, nil
}
