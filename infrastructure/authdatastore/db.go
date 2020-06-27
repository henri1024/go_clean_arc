package authdatastore

import (
	"fmt"
	"os"
	"strconv"

	"github.com/go-redis/redis"
)

func NewRedisDB() *redis.Client {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")
	db, _ := strconv.ParseInt(os.Getenv("REDIS_DB"), 10, 64)

	redisDB := redis.NewClient(
		&redis.Options{
			Addr:     fmt.Sprintf("%s:%s", host, port),
			Password: password,
			DB:       int(db),
		},
	)

	return redisDB
}
