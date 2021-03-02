package data

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// type User struct {
// ID uuid
// Firstname string
// Lastname string
// Email string
// Age uint
// Created time.Time
// }

type User struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	Firstname string    `gorm:"size:256" json:"firstname" validate:"required,min=3"`
	Lastname  string    `gorm:"size:256" json:"lastname" validate:"required,min=3"`
	Email     string    `gorm:"size:256;unique" json:"email" validate:"required,email"`
	Age       uint      `json:"age" validate:"required,gte=0,lte=120"`
	Created   time.Time `json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	uuidNew, err := uuid.NewV4()
	if err != nil {
		return err
	}
	u.ID = uuidNew

	u.Created = time.Now()
	return nil
}
