package repository

import (
	"errors"
	"fmt"
	"go_clean_arc/domain"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

type authRepository struct {
	rdb *redis.Client
}

func NewAuthRepository(rdb *redis.Client) domain.AuthRepository {
	return &authRepository{
		rdb: rdb,
	}
}

func (tr *authRepository) SaveToken(uid uint, token *domain.Token) error {
	at := time.Unix(token.AccessExpired, 0)
	rt := time.Unix(token.RefreshExpired, 0)
	now := time.Now()

	atCreated, err := tr.rdb.Set(token.TokenUuid, strconv.Itoa(int(uid)), at.Sub(now)).Result()
	if err != nil {
		fmt.Println("1", err)
		return err
	}
	rtCreated, err := tr.rdb.Set(token.RefreshUuid, strconv.Itoa(int(uid)), rt.Sub(now)).Result()
	if err != nil {
		fmt.Println("2", err)
		return err
	}
	if atCreated == "0" || rtCreated == "0" {
		return errors.New("no record inserted")
	}
	return nil
}
