package jwtAuth

import (
	"os"

	"github.com/dgrijalva/jwt-go"
)

type JwtWidget struct {
	secretAccessKey  []byte
	secretRefreshKey []byte
}

func NewJwtWidget() *JwtWidget {
	return &JwtWidget{
		secretAccessKey:  []byte(os.Getenv("SECRET_ACCESS_KEY")),
		secretRefreshKey: []byte(os.Getenv("SECRET_REFRESH_KEY")),
	}
}

func (jw *JwtWidget) CreateToken(flag, tokenID string, uid uint, tokenExp int64) (string, error) {
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = tokenID
	atClaims["user_id"] = uid
	atClaims["exp"] = tokenExp

	if flag == "ACCESS" {
		unsignedAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
		return unsignedAccessToken.SignedString(jw.secretAccessKey)
	}
	unsignedAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	return unsignedAccessToken.SignedString(jw.secretRefreshKey)

}
