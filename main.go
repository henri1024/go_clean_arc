package main

import (
	"clean_arc/infrastructure/datastore"
	"clean_arc/infrastructure/router"
	"clean_arc/registry"
	"log"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("no env gotten")
	}
}
func main() {
	db := datastore.NewDB()
	db.LogMode(true)
	defer db.Close()

	err := datastore.PopulateDB(db)
	if err != nil {
		log.Fatal(err)
	}

	registry := registry.NewRegistry(db)

	router := router.NewRouter(registry.NewAppController())

	router.Run()
}
