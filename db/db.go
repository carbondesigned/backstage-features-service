package db

import (
	"fmt"
	"log"
	"os"

	"github.com/digitalocean/sample-golang/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Dbinstance struct {
	Db *gorm.DB
}

var DB Dbinstance

func ConnectDb() {
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_PORT := os.Getenv("DB_PORT")
	DB_USER := os.Getenv("DB_USER")
	DB_ADDR := os.Getenv("DB_HOST")
	DB_NAME := os.Getenv("DB_NAME")

	connectionString := fmt.Sprintf(
		"postgresql://%v:%v@%v:%v/%v?sslmode=require",
		DB_USER,
		DB_PASSWORD,
		DB_ADDR,
		DB_PORT,
		DB_NAME,
	)

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database \n", err)
		os.Exit(2)
	}
	log.Println("connected")
	db.Logger = logger.Default.LogMode(logger.Info)
	log.Println("running migrations")
	db.AutoMigrate(&models.Author{}, &models.Post{})

	DB = Dbinstance{Db: db}
}
