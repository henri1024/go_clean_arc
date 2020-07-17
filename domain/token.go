package domain

import "net/http"

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
	DeleteTokens(accessDetails *AccessDetails) error
	IsValid(r *http.Request) error
}

type TokenRepository interface {
	SaveToken(uint, *Token) error
	DeleteTokens(accessDetails *AccessDetails) error
}

func (t *Token) ToPublic() *PublicToken {
	return &PublicToken{
		AccessToken:  t.AccessToken,
		RefreshToken: t.RefreshToken,
	}
}
