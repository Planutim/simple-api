package helpers

import (
	"fmt"
	"os"

	"github.com/Planutim/simple-api/data"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

//Helper is a wrapper for gorm.db
type Helper struct {
	db *gorm.DB
}

//NewHelper returns new instance of DbHelper
func NewHelper() (*Helper, error) {

	dsn := fmt.Sprintf("host=%s port=%s user=%s	password=%s dbname=%s sslmode=disable",
		os.Getenv("PG_HOST"),
		os.Getenv("PG_PORT"),
		os.Getenv("PG_USER"),
		os.Getenv("PG_PASSWORD"),
		os.Getenv("PG_DBNAME"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&data.User{})
	return &Helper{
		db,
	}, nil
}

func NewTestHelper() (*Helper, error) {

	dsn := fmt.Sprintf("host=%s port=%s user=%s	password=%s dbname=%s sslmode=disable",
		os.Getenv("TEST_PG_HOST"),
		os.Getenv("TEST_PG_PORT"),
		os.Getenv("TEST_PG_USER"),
		os.Getenv("TEST_PG_PASSWORD"),
		os.Getenv("TEST_PG_DBNAME"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
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

//GetUsers get the list of users
func (h *Helper) GetUsers() ([]*data.User, error) {
	var users []*data.User

	err := h.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

//CreateUser creates user
func (h *Helper) CreateUser(user *data.User) (*data.User, error) {
	err := h.db.Model(data.User{}).Create(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

//ListUser gets User with specified id
func (h *Helper) ListUser(uuidStr string) (*data.User, error) {
	var u data.User
	// err := h.db.Where("id = ?", uuidStr).Take(&u).Error
	err := h.db.First(&u, "id = ?", uuidStr).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

//UpdateUser updates user
func (h *Helper) UpdateUser(user *data.User) (*data.User, error) {
	var userFound data.User

	err := h.db.First(&userFound, "id = ?", user.ID).Error
	if err != nil {
		return nil, err
	}
	err = h.db.Save(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
