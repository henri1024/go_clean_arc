package entity

import (
	"html"
	"regexp"
	"strings"

	"clean_arc/infrastructure/security"

	"github.com/jinzhu/gorm"
)

// User is a data structure (Model)
type User struct {
	gorm.Model
	Email    string `json:"email" gorm:"UNIQUE"`
	Password string `json:"password"`
	Username string `json:"username" gorm:"UNIQUE"`
}

const (
	emailregexp    = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
	usernameregexp = "^[a-z0-9]+(?:_[a-zA-Z0-9]+)*$"
)

type PublicUser struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

func (u *User) ToPublic() *PublicUser {
	return &PublicUser{
		ID:       u.ID,
		Email:    u.Email,
		Username: u.Username,
	}
}

func (u *User) Prepare() {
	u.Email = strings.ToLower(html.EscapeString(strings.TrimSpace(u.Email)))
	u.Username = strings.ToLower(html.EscapeString(strings.TrimSpace(u.Username)))
}

func (u *User) ValidSave() (map[string]string, bool) {

	emailregexp := regexp.MustCompile(emailregexp)
	usernameregexp := regexp.MustCompile(usernameregexp)

	msg := make(map[string]string)

	u.Prepare()

	if u.Email == "" {
		msg["email_required"] = "email is required"
	} else if !emailregexp.MatchString(u.Email) {
		msg["invalid_email"] = "invalid email address"
	}

	if u.Password == "" {
		msg["password_required"] = "password is required"
	} else if len(u.Password) < 6 {
		msg["invalid_password"] = "password must at least 6 character"
	} else if len(u.Password) > 12 {
		msg["invalid_password"] = "password must not more than 6 character"
	}

	if u.Username == "" {
		msg["username_required"] = "username is required"
	} else if !usernameregexp.MatchString(u.Username) {
		msg["invalid_username"] = "invalid username format"
	}

	if len(msg) == 0 {
		err := u.HashPassword()
		if err != nil {
			msg["invalid_password"] = "cant hash your password, try another"
			return msg, false
		}
		return nil, true
	}
	return msg, false
}

func (u *User) ValidGetByEmailAndPassword() (map[string]string, bool) {

	emailregexp := regexp.MustCompile(emailregexp)

	msg := make(map[string]string)

	u.Prepare()

	if u.Email == "" {
		msg["email_required"] = "email is required"
	} else if !emailregexp.MatchString(u.Email) {
		msg["invalid_email"] = "invalid email address"
	}

	if u.Password == "" {
		msg["password_required"] = "password is required"
	} else if len(u.Password) < 6 {
		msg["invalid_password"] = "password must at least 6 character"
	} else if len(u.Password) > 12 {
		msg["invalid_password"] = "password must not more than 6 character"
	}

	if len(msg) == 0 {
		return nil, true
	}

	return msg, false
}

func (u *User) HashPassword() error {
	password, err := security.Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = password
	return nil
}

func (u *User) ComparePassword(password string) bool {
	return security.ComparePassword(u.Password, password)
}
