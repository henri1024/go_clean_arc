package usecase

import (
	"fmt"
	"net/http"
	"os"
	"store/domain"
	"store/infrastructure/uuidgenerator"
	"strconv"
	"strings"
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

func (tu *tokenUsecase) DeleteTokens(accessDetails *domain.AccessDetails) error {
	return tu.tokenRepository.DeleteTokens(accessDetails)
}

func extractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func verifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := extractToken(r)
	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			//Make sure that the token method conform to "SigningMethodHMAC"
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("SECRET_ACCESS_KEY")), nil
		})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (tu *tokenUsecase) ExtractTokenMetadata(r *http.Request) (*domain.AccessDetails, error) {
	token, err := verifyToken(r)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}
		return &domain.AccessDetails{
			AccessUuid: accessUuid,
			UserId:     userId,
		}, nil
	}
	return nil, err
}

func (tu *tokenUsecase) IsValid(r *http.Request) error {
	token, err := verifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}
