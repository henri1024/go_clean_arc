package mocks

import (
	"store/domain"

	"github.com/stretchr/testify/mock"
)

type MockUserUsecase struct {
	mock.Mock
}

func (m *MockUserUsecase) SaveUser(user *domain.User) error {
	args := m.Called(user)

	return args.Error(0)
}

func (m *MockUserUsecase) GetUserByEmailAndPassword(email, password string) (*domain.PublicUser, error) {
	args := m.Called(email, password)

	pUser := args.Get(0).(*domain.PublicUser)
	err := args.Error(1)

	return pUser, err
}
