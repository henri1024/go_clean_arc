package domain

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmailRegexp(t *testing.T) {
	samples := []struct {
		text   string
		result bool
	}{
		{
			text:   "validemail@gmail.com",
			result: true,
		},
		{
			text:   "validemail@gmail",
			result: true,
		},
		{
			text:   "unvalidemail.com",
			result: false,
		},
		{
			text:   "unvalid.gmail.com",
			result: false,
		},
	}
	emailregexp := regexp.MustCompile(emailregexp)

	for _, sample := range samples {
		t.Run(sample.text, func(t *testing.T) {
			assert.Equal(t, sample.result, emailregexp.MatchString(sample.text))
		})
	}
}

func TestUsernameRegexp(t *testing.T) {
	samples := []struct {
		text   string
		result bool
	}{
		{
			text:   "valid",
			result: true,
		},
		{
			text:   "__",
			result: false,
		},
		{
			text:   "unvalid one",
			result: false,
		},
		{
			text:   "notvalid!@#$",
			result: false,
		},
	}
	usernameregexp := regexp.MustCompile(usernameregexp)

	for _, sample := range samples {
		t.Run(sample.text, func(t *testing.T) {
			assert.Equal(t, sample.result, usernameregexp.MatchString(sample.text))
		})
	}
}

func TestPrepareUser(t *testing.T) {
	samples := []struct {
		test   User
		result User
	}{
		{
			test: User{
				Email:    "CAPITALEMAIL@GMAIL.COM",
				Username: "CAPITALUSERNAME",
			},
			result: User{
				Email:    "capitalemail@gmail.com",
				Username: "capitalusername",
			},
		},
		{
			test: User{
				Email:    "semiCAPITALEMAIL@GMAIL.COM",
				Username: "semiCAPITALUSERNAME",
			},
			result: User{
				Email:    "semicapitalemail@gmail.com",
				Username: "semicapitalusername",
			},
		},
	}

	for _, sample := range samples {
		t.Run(sample.test.Email, func(t *testing.T) {
			sample.test.Prepare()
			assert.Equal(t, sample.result.Email, sample.test.Email)
			assert.Equal(t, sample.result.Username, sample.test.Username)
		})
	}
}

func TestSaveUserValidation(t *testing.T) {
	samples := []struct {
		user     User
		expected bool
	}{
		{
			user: User{
				Email:    "validemail@gmail.com",
				Username: "username",
				Password: "password",
			},
			expected: true,
		},
		{
			user: User{
				Email:    "",
				Username: "username",
				Password: "password",
			},
			expected: false,
		},
		{
			user: User{
				Email:    "validemail@gmail.com",
				Username: "username",
				Password: "",
			},
			expected: false,
		},
		{
			user: User{
				Email:    "validemail@gmail.com",
				Username: "",
				Password: "password",
			},
			expected: false,
		},
		{
			user: User{
				Email:    "validemail@gmail.com",
				Username: "username",
				Password: "sht",
			},
			expected: false,
		},
		{
			user: User{
				Email:    "validemail@gmail.com",
				Username: "username",
				Password: "passwordtolong",
			},
			expected: false,
		},
		{
			user: User{
				Email:    "validemail@gmail.com",
				Username: "sht",
				Password: "password",
			},
			expected: false,
		},
		{
			user: User{
				Email:    "validemail@gmail.com",
				Username: "usernametolonglonglong",
				Password: "password",
			},
			expected: false,
		},
	}

	for _, sample := range samples {
		t.Run(sample.user.Username, func(t *testing.T) {
			_, tested := sample.user.ValidSave()
			assert.Equal(t, sample.expected, tested)
		})
	}
}

func TestGetUserValidation(t *testing.T) {
	samples := []struct {
		user     User
		expected bool
	}{
		{
			user: User{
				Email:    "validemail@gmail.com",
				Password: "password",
			},
			expected: true,
		},
		{
			user: User{
				Email:    "",
				Password: "password",
			},
			expected: false,
		},
		{
			user: User{
				Email:    "validemail@gmail.com",
				Password: "",
			},
			expected: false,
		},
		{
			user: User{
				Email:    "validemail@gmail.com",
				Password: "sht",
			},
			expected: false,
		},
		{
			user: User{
				Email:    "validemail@gmail.com",
				Password: "passwordtolong",
			},
			expected: false,
		},
	}

	for _, sample := range samples {
		t.Run(sample.user.Email, func(t *testing.T) {
			_, tested := sample.user.ValidGetByEmailAndPassword()
			assert.Equal(t, sample.expected, tested)
		})
	}
}
