package token

import (
	"userauth/domain"

	"github.com/stretchr/testify/mock"
)

type MockTokenRepository struct {
	mock.Mock
}

func (m *MockTokenRepository) SaveToken(uid uint, token *domain.Token) error {
	args := m.Called(uid, token)

	return args.Error(0)
}

func (m *MockTokenRepository) DeleteToken(tokenstring string) error {
	args := m.Called(tokenstring)

	return args.Error(0)
}
