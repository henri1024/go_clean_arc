package hash

import "golang.org/x/crypto/bcrypt"

type Hasher struct{}

func NewHasher() *Hasher {
	return &Hasher{}
}

func (h *Hasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 7)
	if err != nil {
		return "", err
	}
	password = string(hash)
	return password, nil
}

func (h *Hasher) ComparePassword(realPassowrd, inputPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(realPassowrd), []byte(inputPassword))
	if err != nil {
		return false
	}
	return true
}
