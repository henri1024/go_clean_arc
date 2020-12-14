package domain

import "time"

type User struct {
	ID         uint `gorm:"primary_key"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time `sql:"index"`
	Email      string     `json:"email" gorm:"UNIQUE"`
	Password   string     `json:"password"`
	Username   string     `json:"username" gorm:"UNIQUE"`
	IsStaff    bool
	IsVerified bool
}

type PublicUser struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type UserRepository interface {
	SaveUser(*User) error
	GetUserByEmail(string) (*User, error)
	GetUserById(uint64) (*User, error)
}

type UserUsecase interface {
	ValidateUserSignup(*User) (map[string]string, bool)
	ValidateUserSignin(*User) (map[string]string, bool)
	SaveUser(*User) error
	ToPublic(*User) *PublicUser
	GetUserByEmailAndPassword(string, string) (*PublicUser, error)
	GetUserProfile(uint64) (*User, error)
}
