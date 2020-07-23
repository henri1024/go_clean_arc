package usecase

import (
	"fmt"
	"net/http"
	"os"
	"userauth/domain"
	"userauth/infrastructure/uuidgenerator"
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
	rtClaims["refresh_uuid"] = token.RefreshUuid
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

func (tu *tokenUsecase) DeleteToken(tokenstring string) error {
	return tu.tokenRepository.DeleteToken(tokenstring)
}

func extractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func (tu *tokenUsecase) ExtractTokenMetadata(r *http.Request) (*domain.AccessDetails, error) {
	tokenString := extractToken(r)
	token, err := tu.VerifyToken(tokenString, "access")
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

func (tu *tokenUsecase) IsValidRequest(r *http.Request) error {
	tokenString := extractToken(r)
	token, err := tu.VerifyToken(tokenString, "access")
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

func (tu *tokenUsecase) IsValidToken(token *jwt.Token) error {

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return fmt.Errorf("unvalid token")
	}
	return nil
}

func (tu *tokenUsecase) VerifyToken(tokenstring, mode string) (*jwt.Token, error) {
	if mode == "access" {
		token, err := jwt.Parse(
			tokenstring,
			func(token *jwt.Token) (interface{}, error) {
				//Make sure that the token method conform to "SigningMethodHMAC"
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(os.Getenv("SECRET_ACCESS_KEY")), nil
			})
		return token, err
	} else {
		token, err := jwt.Parse(
			tokenstring,
			func(token *jwt.Token) (interface{}, error) {
				//Make sure that the token method conform to "SigningMethodHMAC"
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(os.Getenv("SECRET_REFRESH_KEY")), nil
			})
		return token, err
	}
}

func (tu *tokenUsecase) RefreshToken(token *jwt.Token) error {
	return nil
}
