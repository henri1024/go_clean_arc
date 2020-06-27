package main

import (
	"clean_arc/infrastructure/authdatastore"
	"clean_arc/infrastructure/datastore"
	"clean_arc/infrastructure/router"
	"clean_arc/registry"
	"fmt"
	"log"

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

	redisDB := authdatastore.NewRedisDB()
	fmt.Println(redisDB)

	registry := registry.NewRegistry(psqlDB)

	router := router.NewRouter(registry.NewAppController())

	router.Run()
}
