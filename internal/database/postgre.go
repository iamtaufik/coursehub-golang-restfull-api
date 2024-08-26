package database

import (
	"fmt"
	"log"

	"github.com/iamtaufik/coursehub-golang-restfull-api/internal/config"
	"github.com/iamtaufik/coursehub-golang-restfull-api/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func StartDB() *gorm.DB {

	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable",
		config.Config("DB_HOST"),
		config.Config("DB_USER"),
		config.Config("DB_PASSWORD"),
		config.Config("DB_NAME"),
		config.Config("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("Error:", err)
	}

	if err := db.AutoMigrate(
		&model.User{},
		&model.Profile{},
		&model.Category{},
		&model.Module{},
		&model.Chapter{},
		&model.Course{},
	); err != nil {
		log.Fatalln("Error:", err)
	}
	

	return db
}