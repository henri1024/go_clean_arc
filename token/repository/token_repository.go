package repository

import (
	"errors"
	"fmt"
	"store/domain"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

type tokenRepository struct {
	rdb *redis.Client
}

func NewTokenRepository(rdb *redis.Client) domain.TokenRepository {
	return &tokenRepository{
		rdb: rdb,
	}
}

func (tr *tokenRepository) SaveToken(uid uint, token *domain.Token) error {
	at := time.Unix(token.AccessExpired, 0)
	rt := time.Unix(token.RefreshExpired, 0)
	now := time.Now()

	atCreated, err := tr.rdb.Set(token.TokenUuid, strconv.Itoa(int(uid)), at.Sub(now)).Result()
	if err != nil {
		return err
	}
	rtCreated, err := tr.rdb.Set(token.RefreshUuid, strconv.Itoa(int(uid)), rt.Sub(now)).Result()
	if err != nil {
		return err
	}
	if atCreated == "0" || rtCreated == "0" {
		return errors.New("no record inserted")
	}
	return nil
}

func (tr *tokenRepository) DeleteToken(tokenstring string) error {

	atDelete, err := tr.rdb.Del(tokenstring).Result()
	if err != nil {
		return err
	}

	if atDelete != 1 {
		return fmt.Errorf("something went wrong %v", atDelete)
	}
	return nil
}
