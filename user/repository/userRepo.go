package repository

import (
	"go_clean_arc/domain"

	"github.com/jinzhu/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepo{
		db: db,
	}
}

func (ur *userRepo) SaveUser(user *domain.User) error {
	return ur.db.Model(&domain.User{}).Save(user).Error
}

func (ur *userRepo) GetUserByEmail(email string) (*domain.User, error) {
	user := &domain.User{}

	err := ur.db.Model(&domain.User{}).First(user, "email = ?", email).Error

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (ur *userRepo) GetUserById(id uint64) (*domain.User, error) {
	user := &domain.User{}

	err := ur.db.Model(&domain.User{}).First(user, id).Error

	if err != nil {
		return nil, err
	}

	return user, nil
}
