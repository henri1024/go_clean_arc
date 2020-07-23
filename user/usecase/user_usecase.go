package usecase

import (
	"errors"
	"userauth/domain"
	"userauth/infrastructure/security"
)

type userUsecase struct {
	userRepository domain.UserRepository
}

func NewUserUsecase(userRepository domain.UserRepository) domain.UserUsecase {
	return &userUsecase{
		userRepository: userRepository,
	}
}

func (uu *userUsecase) SaveUser(user *domain.User) error {

	hashedPass, err := security.Hash(user.Password)
	if err != nil {
		return errors.New("failed to hash password")
	}

	user.Password = hashedPass

	err = uu.userRepository.SaveUser(user)

	if err != nil {
		return err
	}

	return nil
}

func (uu *userUsecase) GetUserByEmailAndPassword(email, password string) (*domain.PublicUser, error) {
	user, err := uu.userRepository.GetUserByEmail(email)

	if err != nil {
		return nil, err
	}

	if security.ComparePassword(user.Password, password) {
		return user.ToPublic(), nil

	}

	return nil, errors.New("invalid password")
}
