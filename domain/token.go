package domain

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

type Token struct {
	AccessToken    string
	RefreshToken   string
	TokenUuid      string
	RefreshUuid    string
	AccessExpired  int64
	RefreshExpired int64
}

type PublicToken struct {
	AccessToken  string
	RefreshToken string
}

type AccessDetails struct {
	AccessUuid string
	UserId     uint64
}

type TokenUsecase interface {
	CreateToken(uint) (*Token, error)
	SaveToken(uint, *Token) error
	ExtractTokenMetadata(*http.Request) (*AccessDetails, error)
	DeleteToken(string) error
	IsValidRequest(*http.Request) error
	IsValidToken(*jwt.Token) error
	VerifyToken(string, string) (*jwt.Token, error)
}

type TokenRepository interface {
	SaveToken(uint, *Token) error
	DeleteToken(string) error
}

func (t *Token) ToPublic() *PublicToken {
	return &PublicToken{
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
	}
}
