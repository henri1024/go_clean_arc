package domain

type Token struct {
	AccessToken    string
	RefreshToken   string
	TokenUuid      string
	RefreshUuid    string
	AccessExpired  int64
	RefreshExpired int64
}

type TokenUsecase interface {
	CreateToken(uint) (*Token, error)
	SaveToken(uint, *Token) error
}

type TokenRepository interface {
	SaveToken(uint, *Token) error
}
