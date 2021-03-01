package db

import (
	"fmt"
	"os"

	"github.com/Planutim/simple-api/data"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//Helper is a wrapper for gorm.db
type Helper struct {
	DB *gorm.DB
}

//NewHelper returns new instance of DbHelper
func NewHelper() (*Helper, error) {
	err := godotenv.Load("../.env")
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s	password=%s dbname=%s sslmode=disable",
		os.Getenv("TEST_PG_HOST"),
		os.Getenv("TEST_PG_PORT"),
		os.Getenv("TEST_PG_USER"),
		os.Getenv("TEST_PG_PASSWORD"),
		os.Getenv("TEST_PG_DBNAME"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&data.User{})
	return &Helper{
		db,
	}, nil
}

func NewTestHelper() (*Helper, error) {
	err := godotenv.Load("../.env")
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s	password=%s dbname=%s sslmode=disable",
		os.Getenv("TEST_PG_HOST"),
		os.Getenv("TEST_PG_PORT"),
		os.Getenv("TEST_PG_USER"),
		os.Getenv("TEST_PG_PASSWORD"),
		os.Getenv("TEST_PG_DBNAME"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	//delete table
	db.Exec("DROP TABLE IF EXISTS users;")
	db.AutoMigrate(&data.User{})
	return &Helper{
		db,
	}, nil
}
