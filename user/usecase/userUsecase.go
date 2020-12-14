package usecase

import (
	"errors"
	"go_clean_arc/domain"
	"go_clean_arc/infrastructure/hash"
	"html"
	"regexp"
	"strings"
)

const (
	emailregexp    = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
	usernameregexp = "^[a-z0-9]+(?:_[a-zA-Z0-9]+)*$"
)

type userUsecase struct {
	userRepo domain.UserRepository
	hasher   *hash.Hasher
}

func NewUserUsecase(userRepo domain.UserRepository, hasher *hash.Hasher) domain.UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
		hasher:   hasher,
	}
}

func (uu *userUsecase) Prepare(u *domain.User) *domain.User {
	u.Email = strings.ToLower(html.EscapeString(strings.TrimSpace(u.Email)))
	u.Username = strings.ToLower(html.EscapeString(strings.TrimSpace(u.Username)))

	return u
}

func (uu *userUsecase) ValidateUserSignup(u *domain.User) (map[string]string, bool) {
	emailregexp := regexp.MustCompile(emailregexp)
	usernameregexp := regexp.MustCompile(usernameregexp)

	msg := make(map[string]string)

	u = uu.Prepare(u)

	if u.Email == "" {
		msg["email"] = "email is required"
	} else if !emailregexp.MatchString(u.Email) {
		msg["email"] = "invalid email address"
	}

	if u.Password == "" {
		msg["password"] = "password is required"
	} else if len(u.Password) < 6 {
		msg["password"] = "password must at least 6 character"
	} else if len(u.Password) > 12 {
		msg["password"] = "password must not more than 12 character"
	}

	if u.Username == "" {
		msg["username"] = "username is required"
	} else if !usernameregexp.MatchString(u.Username) {
		msg["username"] = "invalid username format"
	} else if len(u.Username) < 4 {
		msg["username"] = "username must at least 4 character"
	} else if len(u.Username) > 20 {
		msg["username"] = "username must not more than 20 character"
	}

	if len(msg) == 0 {
		return msg, true
	}
	return msg, false
}

func (uu *userUsecase) ValidateUserSignin(u *domain.User) (map[string]string, bool) {

	emailregexp := regexp.MustCompile(emailregexp)

	msg := make(map[string]string)

	u = uu.Prepare(u)

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

func (uu *userUsecase) SaveUser(u *domain.User) error {

	var (
		hashedPassword string
		err            error
	)
	if hashedPassword, err = uu.hasher.Hash(u.Password); err != nil {
		return err
	}

	u.Password = hashedPassword
	return uu.userRepo.SaveUser(u)
}

func (uu *userUsecase) GetUserByEmailAndPassword(email, password string) (*domain.PublicUser, error) {
	user, err := uu.userRepo.GetUserByEmail(email)

	if err != nil {
		return nil, err
	}

	if uu.hasher.ComparePassword(user.Password, password) {
		return uu.ToPublic(user), nil
	}

	return nil, errors.New("invalid password")
}

func (uu *userUsecase) ToPublic(u *domain.User) *domain.PublicUser {
	return &domain.PublicUser{
		ID:       u.ID,
		Email:    u.Email,
		Username: u.Username,
	}
}

func (uu *userUsecase) GetUserProfile(uid uint64) (*domain.User, error) {
	return uu.userRepo.GetUserById(uid)
}
