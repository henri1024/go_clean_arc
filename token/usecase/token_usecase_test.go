package usecase

import (
	"errors"
	"userauth/domain"
	tokenmock "userauth/domain/mocks/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateToken(t *testing.T) {
	mockRepo := new(tokenmock.MockTokenRepository)

	t.Run("success", func(t *testing.T) {

		var uid uint = 12
		testUsecase := NewTokenUsecase(mockRepo)

		_, err := testUsecase.CreateToken(uid)

		mockRepo.AssertExpectations(t)
		assert.Equal(t, nil, err)
	})
}

func generateToken() *domain.Token {
	var uid uint = 12

	testUsecase := NewTokenUsecase(nil)
	token, _ := testUsecase.CreateToken(uid)

	return token
}

func TestSaveToken(t *testing.T) {
	mockRepo := new(tokenmock.MockTokenRepository)

	t.Run("success", func(t *testing.T) {
		token := generateToken()

		mockRepo.On("SaveToken", mock.Anything, mock.Anything).Return(nil).Once()
		testUsecase := NewTokenUsecase(mockRepo)

		result := testUsecase.SaveToken(12, token)

		mockRepo.AssertExpectations(t)
		assert.Equal(t, nil, result)
	})

	t.Run("failed", func(t *testing.T) {

		err := errors.New("some error")
		token := generateToken()

		mockRepo.On("SaveToken", mock.Anything, mock.Anything).Return(err).Once()
		testUsecase := NewTokenUsecase(mockRepo)

		result := testUsecase.SaveToken(12, token)

		mockRepo.AssertExpectations(t)
		assert.Equal(t, err, result)
	})
}
