package domain

import (
	"html"
	"regexp"
	"strings"
	"time"
)

type User struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	Email     string     `json:"email" gorm:"UNIQUE"`
	Password  string     `json:"password"`
	Username  string     `json:"username" gorm:"UNIQUE"`
	IsStaff   bool
}

type UserRepository interface {
	SaveUser(u *User) error
	GetUserByEmail(email string) (*User, error)
}

type UserUsecase interface {
	SaveUser(u *User) error
	GetUserByEmailAndPassword(email, password string) (*PublicUser, error)
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
		msg["invalid_password"] = "password must not more than 12 character"
	}

	if u.Username == "" {
		msg["username_required"] = "username is required"
	} else if !usernameregexp.MatchString(u.Username) {
		msg["invalid_username"] = "invalid username format"
	} else if len(u.Username) < 4 {
		msg["invalid_username"] = "username must at least 4 character"
	} else if len(u.Username) > 20 {
		msg["invalid_username"] = "username must not more than 20 character"
	}

	if len(msg) == 0 {
		return msg, true
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
		msg["invalid_email"] = "invalid email format"
	}

	if u.Password == "" {
		msg["password_required"] = "password is required"
	} else if len(u.Password) < 6 {
		msg["invalid_password"] = "password must at least 6 character"
	} else if len(u.Password) > 12 {
		msg["invalid_password"] = "password must not more than 12 character"
	}

	if len(msg) == 0 {
		return msg, true
	}

	return msg, false
}
