package datastore

import (
	"clean_arc/domain/entity"
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func NewDB() *gorm.DB {
	DBMS := os.Getenv("DB_DRIVER")
	DBURL := generateDBURL()

	db, err := gorm.Open(DBMS, DBURL)

	if err != nil {
		log.Fatalln(err)
	}

	db.Debug().DropTableIfExists(&entity.User{})
	db.Debug().AutoMigrate(&entity.User{})

	return db
}

func generateDBURL() string {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", host, port, user, name, password)
}

func PopulateDB(db *gorm.DB) error {
	users := []entity.User{
		{
			Email:    "first@example.com",
			Password: "password",
			Username: "username1",
		},
		{
			Email:    "second@example.com",
			Password: "password",
			Username: "username2",
		},
		{
			Email:    "third@example.com",
			Password: "password",
			Username: "username3",
		},
	}

	for _, user := range users {
		user.HashPassword()
		err := db.Model(&entity.User{}).Save(&user).Error
		if err != nil {
			return err
		}
	}

	fmt.Println("Success populate DB")
	return nil
}
