package usecase

import (
	"errors"
	"fmt"
	"go_clean_arc/domain"
	"go_clean_arc/infrastructure/jwtAuth"
	uuid "go_clean_arc/infrastructure/uuid"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type authUsecase struct {
	tokenRepository domain.AuthRepository
	idGen           *uuid.UuidGenerator
	jwt             *jwtAuth.JwtWidget
}

type tokenString struct {
	token string
	err   error
	mode  string
}

func NewAuthUsecase(tokenRepository domain.AuthRepository, idGen *uuid.UuidGenerator, jwtauth *jwtAuth.JwtWidget) domain.AuthUsecase {
	return &authUsecase{
		tokenRepository: tokenRepository,
		idGen:           idGen,
		jwt:             jwtauth,
	}
}

const (
	ACCESSEXPTIME  = 15 * time.Minute
	REFRESHEXPTIME = 3 * 24 * time.Hour
)

func (au *authUsecase) CreateToken(uid uint) (*domain.Token, error) {

	token := &domain.Token{
		AccessExpired:  time.Now().Add(ACCESSEXPTIME).Unix(),
		TokenUuid:      au.idGen.NewId(),
		RefreshExpired: time.Now().Add(REFRESHEXPTIME).Unix(),
		RefreshUuid:    au.idGen.NewId(),
	}

	ch := make(chan *tokenString)
	defer close(ch)

	go func() {
		tkn, err := au.jwt.CreateToken("ACCESS", token.TokenUuid, uid, token.AccessExpired)
		ch <- &tokenString{
			token: tkn,
			err:   err,
			mode:  "ACCESS",
		}
	}()

	go func() {
		tkn, err := au.jwt.CreateToken("REFRESH", token.RefreshUuid, uid, token.RefreshExpired)
		ch <- &tokenString{
			token: tkn,
			err:   err,
			mode:  "REFRESH",
		}
	}()

	for i := 0; i < 2; i++ {
		temp := <-ch
		if temp.err != nil {
			return nil, temp.err
		} else if temp.mode == "ACCESS" {
			token.AccessToken = temp.token
		} else if temp.mode == "REFRESH" {
			token.RefreshToken = temp.token
		}
	}

	// var (
	// 	accessToken  string
	// 	refreshToken string
	// 	err          error
	// )

	// if accessToken, err = au.jwt.CreateToken("ACCESS", token.TokenUuid, uid, token.AccessExpired); err != nil {
	// 	return nil, err
	// }

	// if refreshToken, err = au.jwt.CreateToken("REFRESH", token.RefreshUuid, uid, token.RefreshExpired); err != nil {
	// 	return nil, err
	// }

	// token.AccessToken = accessToken
	// token.RefreshToken = refreshToken

	return token, nil
}

func (au *authUsecase) SaveToken(uid uint, token *domain.Token) error {
	return au.tokenRepository.SaveToken(uid, token)
}

func (au *authUsecase) ToPublic(token *domain.Token) *domain.PublicToken {
	return &domain.PublicToken{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
	}
}

func extractToken(r *http.Request) (string, error) {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1], nil
	}
	return "", errors.New("Invalid token bearer format")
}

func (au *authUsecase) VerifyRequest(r *http.Request, mode string) (*jwt.Token, error) {
	var (
		tokenstr string
		err      error
	)

	if tokenstr, err = extractToken(r); err != nil {
		return nil, err
	}

	if mode == "access" {
		token, err := jwt.Parse(
			tokenstr,
			func(token *jwt.Token) (interface{}, error) {
				//Make sure that the token method conform to "SigningMethodHMAC"
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(os.Getenv("SECRET_ACCESS_KEY")), nil
			})
		return token, err
	}

	return nil, errors.New("Unexpected mode")
}

func (au *authUsecase) ExtractTokenMetadata(r *http.Request, modearr ...string) (*domain.AccessDetails, error) {
	var (
		mode   string
		token  *jwt.Token
		err    error
		claims jwt.MapClaims
		ok     bool
	)

	if len(modearr) == 0 {
		mode = "access"
	} else if modearr[0] == "" {
		mode = "access"
	}

	if token, err = au.VerifyRequest(r, mode); err != nil {
		return nil, err
	}

	claims, ok = token.Claims.(jwt.MapClaims)
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

	return nil, errors.New("unable to extract token")
}

func (au *authUsecase) IsValidToken(token *jwt.Token) error {
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return fmt.Errorf("unvalid token")
	}
	return nil
}

func (au *authUsecase) IsValidRequest(r *http.Request) error {

	var (
		token *jwt.Token
		err   error
	)

	if token, err = au.VerifyRequest(r, "access"); err != nil {
		return err
	}

	if err = au.IsValidToken(token); err != nil {
		return err
	}

	return nil
}
