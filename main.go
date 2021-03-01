package main

import (
	"log"

	"github.com/Planutim/simple-api/db"
	"github.com/Planutim/simple-api/handlers"
	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
	echo "github.com/labstack/echo/v4"
)

// сделать REST API на Go для создания/удаления/редактирования юзеров. Любой framework (или без него). Запушить код на github. В идеале с unit тестами. БД - PostgreSQL. Запросы:

// POST /users - create user
// GET /users/<id> - get user
// PUT /users/<id> - edit user

// type User struct {
// ID uuid
// Firstname string
// Lastname string
// Email string
// Age uint
// Created time.Time
// }

// ID / Created генерим сами. Остальные - обязательны и валидируем на входе.

// Лимита по времени нет, но чем быстрее сделаешь - тем больше вероятность, что возьмем )
func main() {
	godotenv.Load()
	e := echo.New()

	db, err := db.NewHelper()
	if err != nil {
		log.Fatal(err)
	}
	validator := validator.New()
	sh := handlers.NewSimpleHandler(db, validator)

	e.POST("/users", sh.CreateUser)
	e.GET("/home", sh.Home)
	e.Logger.Fatal(e.Start(":8000"))
}
