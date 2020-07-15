package mocks

import (
	"store/domain"
	"store/infrastructure/security"

	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) SaveUser(user *domain.User) error {

	args := m.Called(user)

	return args.Error(0)

}

func (m *MockUserRepository) GetUserByEmail(email string) (*domain.User, error) {
	args := m.Called(email)

	res1 := args.Get(0).(*domain.User)
	res1.Password, _ = security.Hash(res1.Password)
	res2 := args.Error(1)

	return res1, res2
}
