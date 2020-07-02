package entity

import (
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestUserPrepare(t *testing.T) {

	samples := []struct {
		input  User
		result User
	}{
		{
			User{
				Email: "    TEst1@example.com 		",
				Username: "\t\t<username1\n",
			},
			User{
				Email:    "test1@example.com",
				Username: "&lt;username1",
			},
		},
		{
			User{
				Email: "    <TEst2@example.com\n 		",
				Username: "\t\tmyusername2\n		",
			},
			User{
				Email:    "&lt;test2@example.com",
				Username: "myusername2",
			},
		},
	}

	for _, sample := range samples {
		sample.input.Prepare()

		assert.Equal(t, sample.result.Email, sample.input.Email)
		assert.Equal(t, sample.result.Username, sample.input.Username)

	}

}

func TestUserSaveValidation(t *testing.T) {
	samples := []struct {
		input  User
		result map[string]string
	}{
		{
			User{
				Email:    "",
				Username: "",
				Password: "",
			},
			map[string]string{
				"email_required":    "email is required",
				"password_required": "password is required",
				"username_required": "username is required",
			},
		},
		{
			User{
				Email:    "validemail@test.com",
				Username: "valid_name",
				Password: "validPass123",
			},
			map[string]string{},
		},
		{
			User{
				Email:    "not valid@test.com",
				Username: "no",
				Password: "fail",
			},
			map[string]string{
				"invalid_email":    "invalid email address",
				"invalid_username": "username must at least 4 character",
				"invalid_password": "password must at least 6 character",
			},
		},
		{
			User{
				Email:    "@test.com",
				Username: "this one fail",
				Password: "failbecausetolong",
			},
			map[string]string{
				"invalid_email":    "invalid email address",
				"invalid_username": "invalid username format",
				"invalid_password": "password must not more than 12 character",
			},
		},
		{
			User{
				Email:    "notrailing",
				Username: "thisistolongusernamecuzmorethan20",
				Password: "ValidPass123",
			},
			map[string]string{
				"invalid_email":    "invalid email address",
				"invalid_username": "username must not more than 20 character",
			},
		},
	}

	for _, sample := range samples {
		sample.input.Prepare()
		msg, _ := sample.input.ValidSave()

		assert.Equal(t, sample.result, msg)

	}
}

func TestUserGetByEmailAndPasswordValidation(t *testing.T) {
	samples := []struct {
		input  User
		result map[string]string
	}{
		{
			User{
				Email:    "",
				Password: "",
			},
			map[string]string{
				"email_required":    "email is required",
				"password_required": "password is required",
			},
		},
		{
			User{
				Email:    "valid@test.com",
				Password: "ValidPass",
			},
			map[string]string{},
		},
		{
			User{
				Email:    "in valid@test@test.com",
				Password: "pass",
			},
			map[string]string{
				"invalid_email":    "invalid email format",
				"invalid_password": "password must at least 6 character",
			},
		},
		{
			User{
				Email:    "in%&^validtest@test.com",
				Password: "passwordistolong",
			},
			map[string]string{
				"invalid_email":    "invalid email format",
				"invalid_password": "password must not more than 12 character",
			},
		},
	}

	for _, sample := range samples {
		sample.input.Prepare()
		msg, _ := sample.input.ValidGetByEmailAndPassword()

		assert.Equal(t, sample.result, msg)

	}
}

func TestUserToPublicUser(t *testing.T) {
	samples := []struct {
		input  User
		result *PublicUser
	}{
		{
			User{
				Email:    "email@test.com",
				Password: "password",
				Username: "username",
				Model: gorm.Model{
					ID:        1,
					CreatedAt: time.Now(),
					DeletedAt: nil,
					UpdatedAt: time.Now(),
				},
			},
			&PublicUser{
				ID:       1,
				Email:    "email@test.com",
				Username: "username",
			},
		},
	}

	for _, sample := range samples {
		publicUser := sample.input.ToPublic()

		assert.IsType(t, sample.result, publicUser)

		assert.Equal(t, sample.result.Email, publicUser.Email)
		assert.Equal(t, sample.result.ID, publicUser.ID)
		assert.Equal(t, sample.result.Username, publicUser.Username)
	}

}

func TestUserHashPassword(t *testing.T) {
	samples := []struct {
		input User
	}{
		{
			User{
				Email:    "email@test.com",
				Password: "password",
				Username: "username",
			},
		},
	}

	for _, sample := range samples {
		recentPassword := sample.input.Password
		sample.input.HashPassword()

		assert.NotEqual(t, recentPassword, sample.input.Password)
	}
}

func TestUserComparePassword(t *testing.T) {
	samples := []struct {
		input         User
		inputPassword string
		result        bool
	}{
		{
			User{
				Email:    "email@test.com",
				Password: "password",
				Username: "username",
			},
			"password",
			true,
		},
		{
			User{
				Email:    "email@test.com",
				Password: "password",
				Username: "username",
			},
			"wrongpass",
			false,
		},
	}

	for _, sample := range samples {
		sample.input.HashPassword()

		log := sample.input.ComparePassword(sample.inputPassword)

		assert.True(t, log == sample.result)
	}
}
