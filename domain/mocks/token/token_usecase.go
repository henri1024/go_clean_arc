package token

import (
	"userauth/domain"

	"github.com/stretchr/testify/mock"
)

type MockTokenUsecase struct {
	mock.Mock
}

// type TokenUsecase interface {
// 	CreateToken(uint) (*Token, error)
// 	SaveToken(uint, *Token) error
// }

func (m *MockTokenUsecase) CreateToken(userid uint) (*domain.Token, error) {
	args := m.Called(userid)

	token := args.Get(0).(*domain.Token)
	err := args.Error(1)

	return token, err
}

func (m *MockTokenUsecase) SaveToken(id uint, token *domain.Token) error {
	args := m.Called(id, token)

	return args.Error(0)
}
