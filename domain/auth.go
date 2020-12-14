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

type AuthUsecase interface {
	CreateToken(uint) (*Token, error)
	SaveToken(uint, *Token) error
	ToPublic(*Token) *PublicToken
	VerifyRequest(*http.Request, string) (*jwt.Token, error)
	ExtractTokenMetadata(*http.Request, ...string) (*AccessDetails, error)
	IsValidRequest(*http.Request) error
	IsValidToken(*jwt.Token) error
}

type AuthRepository interface {
	SaveToken(uint, *Token) error
}
