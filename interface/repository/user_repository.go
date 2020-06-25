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
	GetUserByEmail(email string) error
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

func (ur *userRepository) GetUserByEmail(email string) error {
	err := ur.db.Model(&entity.User{}).First(&entity.User{Email: email}).Error
	return err
}
