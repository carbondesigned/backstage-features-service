package database

import (
	"fmt"
	"log"
	"os"

	"github.com/carbondesinged/backstage-features-service/models"
	"github.com/ilyakaznacheev/cleanenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Dbinstance struct {
	Db *gorm.DB
}

var DB Dbinstance

type ConfigDatabase struct {
	Port     string `yaml:"port" env:"PORT" env-default:"5432"`
	Host     string `yaml:"host" env:"HOST" env-default:"localhost"`
	Name     string `yaml:"name" env:"NAME" env-default:"postgres"`
	User     string `yaml:"user" env:"USER" env-default:"user"`
	Password string `yaml:"password" env:"PASSWORD"`
}

var cfg ConfigDatabase

func ConnectDb() {
	err := cleanenv.ReadConfig(".env", &cfg)
	if err != nil {
		log.Fatal("error with .env", err)
	}
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_PORT := os.Getenv("DB_PORT")
	DB_USER := os.Getenv("DB_USER")
	DB_NAME := os.Getenv("DB_NAME")
	DB_HOST := os.Getenv("DB_HOST")
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v", DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT)
	// dsn := fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=require", DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME)
	println(dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
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
