package psqldb

import (
	"fmt"
	"go_clean_arc/domain"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func NewDB() (*gorm.DB, error) {
	DBMS := os.Getenv("DB_DRIVER")
	DBURL := generateDBURL()

	var (
		db  *gorm.DB
		err error
	)

	if db, err = gorm.Open(DBMS, DBURL); err != nil {
		return nil, err
	}

	mdl := &domain.User{}

	if err = db.Debug().DropTableIfExists(mdl).Error; err != nil {
		return nil, err
	}
	if err = db.Debug().AutoMigrate(mdl).Error; err != nil {
		return nil, err
	}

	log.Printf("Connected to postgresql\n")

	return db, nil
}

func generateDBURL() string {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	user = os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", host, port, user, name, password)

	// dbUrl := os.Getenv("DATABASE_URL")
	// return dbUrl
}
