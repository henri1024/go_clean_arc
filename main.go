package main

import (
	"go_clean_arc/app"
	"go_clean_arc/infrastructure/hash"
	"go_clean_arc/infrastructure/jwtAuth"
	"go_clean_arc/infrastructure/psqldb"
	"go_clean_arc/infrastructure/redisdb"
	"go_clean_arc/infrastructure/router"
	uuid "go_clean_arc/infrastructure/uuid"

	"log"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("No .env file found\n")
	}
}

func main() {
	var (
		db  *gorm.DB
		rdb *redis.Client
		err error
	)

	if db, err = psqldb.NewDB(); err != nil {
		log.Fatal("Error starting app : \v\n", err)
	}

	if rdb, err = redisdb.NewRedisDB(); err != nil {
		log.Fatal("Error starting app : \v\n", err)
	}

	userApp := app.NewControllers(
		db,
		hash.NewHasher(),
		uuid.NewUuidGenerator(),
		rdb,
		jwtAuth.NewJwtWidget(),
	)

	r := router.NewRouter(userApp)

	log.Fatal(r.Run())
}
