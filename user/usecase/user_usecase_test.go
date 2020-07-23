package usecase

import (
	"errors"
	"userauth/domain"
	usermock "userauth/domain/mocks/user"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSaveUser(t *testing.T) {

	mockRepo := new(usermock.MockUserRepository)
	user := &domain.User{
		ID:       1,
		Email:    "test@email.com",
		IsStaff:  true,
		Password: "password",
		Username: "user",
	}

	// expectation
	mockRepo.On("SaveUser", mock.Anything).Return(nil).Once()
	testUsecase := NewUserUsecase(mockRepo)

	result := testUsecase.SaveUser(user)

	// Mock Assertion
	mockRepo.AssertExpectations(t)

	// Data Assertion
	assert.Equal(t, nil, result)

}

func TestGetUserByEmailAndPassword(t *testing.T) {

	mockRepo := new(usermock.MockUserRepository)
	user := &domain.User{
		ID:       1,
		Email:    "test@email.com",
		IsStaff:  true,
		Password: "password",
		Username: "user",
	}

	t.Run("success", func(t *testing.T) {

		// expectation
		mockRepo.On("GetUserByEmail", mock.Anything).Return(user, nil).Once()
		testUsecase := NewUserUsecase(mockRepo)

		result, _ := testUsecase.GetUserByEmailAndPassword(user.Email, user.Password)

		mockRepo.AssertExpectations(t)
		assert.Equal(t, user.ToPublic(), result)
	})

	t.Run("incorrect password", func(t *testing.T) {

		passError := errors.New("invalid password")

		// expectation
		mockRepo.On("GetUserByEmail", mock.Anything).Return(user, nil).Once()
		testUsecase := NewUserUsecase(mockRepo)

		// test invalid input password
		_, result := testUsecase.GetUserByEmailAndPassword(user.Email, user.Password+"1")

		mockRepo.AssertExpectations(t)
		assert.Equal(t, passError, result)
	})

	t.Run("invalid crud", func(t *testing.T) {

		crudError := errors.New("some error")

		// expectation
		mockRepo.On("GetUserByEmail", mock.Anything).Return(&domain.User{}, crudError).Once()
		testUsecase := NewUserUsecase(mockRepo)

		// test invalid input password
		_, err := testUsecase.GetUserByEmailAndPassword(user.Email, user.Password)

		mockRepo.AssertExpectations(t)
		assert.Equal(t, crudError, err)
	})

}
