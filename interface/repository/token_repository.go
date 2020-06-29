package repository

import (
	"clean_arc/domain/authtoken"
	"clean_arc/infrastructure/uuidgenerator"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type tokenRepository struct {
}

type TokenRepository interface {
	CreateToken(uint) (*authtoken.Token, error)
}

func NewTokenRepository() TokenRepository {
	return &tokenRepository{}
}

func (tr *tokenRepository) CreateToken(userid uint) (*authtoken.Token, error) {

	secretAccessKey := []byte(os.Getenv("SECRET_ACCESS_KEY"))
	secretRefreshKey := []byte(os.Getenv("SECRET_REFRESH_KEY"))

	token := &authtoken.Token{}

	token.AccessExpired = time.Now().Add(10 * time.Minute).Unix()
	token.TokenUuid = uuidgenerator.NewUuid()

	token.RefreshExpired = time.Now().Add(7 * time.Hour * 24).Unix()
	token.RefreshUuid = uuidgenerator.NewUuid()

	// Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = token.TokenUuid
	atClaims["user_id"] = userid
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
	rtClaims["user_id"] = userid
	rtClaims["exp"] = token.RefreshExpired

	unsignedRefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	refreshToken, err := unsignedRefreshToken.SignedString(secretRefreshKey)
	if err != nil {
		return nil, err
	}
	token.RefreshToken = refreshToken

	return token, nil
}
