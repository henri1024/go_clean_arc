package repository

import (
	"clean_arc/domain/authtoken"
	"errors"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

type AuthRepository interface {
	SaveToken(uint, *authtoken.Token) error
}

type authRepository struct {
	redisDB *redis.Client
}

func NewAuthRepository(redisDB *redis.Client) AuthRepository {
	return &authRepository{
		redisDB: redisDB,
	}
}

func (ar *authRepository) SaveToken(userid uint, token *authtoken.Token) error {
	at := time.Unix(token.AccessExpired, 0)
	rt := time.Unix(token.RefreshExpired, 0)
	now := time.Now()

	atCreated, err := ar.redisDB.Set(token.TokenUuid, strconv.Itoa(int(userid)), at.Sub(now)).Result()
	if err != nil {
		return err
	}
	rtCreated, err := ar.redisDB.Set(token.RefreshUuid, strconv.Itoa(int(userid)), rt.Sub(now)).Result()
	if err != nil {
		return err
	}
	if atCreated == "0" || rtCreated == "0" {
		return errors.New("no record inserted")
	}
	return nil
}
