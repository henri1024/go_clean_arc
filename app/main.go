package main

import (
	"log"
	"store/app/registry"
	"store/infrastructure/authdatastore"
	"store/infrastructure/datastore"
	"store/infrastructure/router"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("no env gotten")
	}
}

func main() {
	psqlDB := datastore.NewDB()
	psqlDB.LogMode(true)
	defer psqlDB.Close()

	err := datastore.PopulateDB(psqlDB)
	if err != nil {
		log.Fatal(err)
	}

	redisDB, err := authdatastore.NewRedisDB()
	if err != nil {
		log.Fatal(err)
	}

	registry := registry.NewRegistry(psqlDB, redisDB)

	router := router.NewRouter(registry)
	router.Run()
}
