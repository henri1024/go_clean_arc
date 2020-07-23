package repository

import (
	"userauth/domain"

	"github.com/jinzhu/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (ur *userRepository) SaveUser(user *domain.User) error {
	err := ur.db.Model(&domain.User{}).Save(user).Error

	if err != nil {
		return err
	}

	return nil
}

func (ur *userRepository) GetUserByEmail(email string) (*domain.User, error) {
	user := &domain.User{}

	err := ur.db.Model(&domain.User{}).First(user, "email = ?", email).Error

	if err != nil {
		return nil, err
	}

	return user, nil
}
