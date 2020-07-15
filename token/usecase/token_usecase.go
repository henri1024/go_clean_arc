package usecase

import (
	"os"
	"store/domain"
	"store/infrastructure/uuidgenerator"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type tokenUsecase struct {
	tokenRepository domain.TokenRepository
}

func NewTokenUsecase(tokenRepository domain.TokenRepository) domain.TokenUsecase {
	return &tokenUsecase{
		tokenRepository: tokenRepository,
	}
}

func (tu *tokenUsecase) CreateToken(uid uint) (*domain.Token, error) {
	secretAccessKey := []byte(os.Getenv("SECRET_ACCESS_KEY"))
	secretRefreshKey := []byte(os.Getenv("SECRET_REFRESH_KEY"))

	token := &domain.Token{}

	token.AccessExpired = time.Now().Add(10 * time.Minute).Unix()
	token.TokenUuid = uuidgenerator.NewUuid()

	token.RefreshExpired = time.Now().Add(7 * time.Hour * 24).Unix()
	token.RefreshUuid = uuidgenerator.NewUuid()

	// Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = token.TokenUuid
	atClaims["user_id"] = uid
	atClaims["exp"] = token.AccessExpired

	unsignedAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	accessToken, err := unsignedAccessToken.SignedString(secretAccessKey)
	if err != nil {
		return nil, err
	}
	token.AccessToken = accessToken

	// Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["access_uuid"] = token.RefreshUuid
	rtClaims["user_id"] = uid
	rtClaims["exp"] = token.RefreshExpired

	unsignedRefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	refreshToken, err := unsignedRefreshToken.SignedString(secretRefreshKey)
	if err != nil {
		return nil, err
	}
	token.RefreshToken = refreshToken

	return token, nil
}

func (tu *tokenUsecase) SaveToken(uid uint, token *domain.Token) error {
	return tu.tokenRepository.SaveToken(uid, token)
}
