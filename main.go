package main

import (
	"log"

	"github.com/Planutim/simple-api/handlers"
	"github.com/Planutim/simple-api/helpers"
	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	db, err := helpers.NewHelper()
	if err != nil {
		log.Fatal(err)
	}
	validator := validator.New()
	sh := handlers.NewSimpleHandler(db, validator)

	e.Pre(middleware.RemoveTrailingSlash())
	e.GET("/users", sh.ListUsers)
	e.GET("/users/:id", sh.ListUser)
	e.POST("/users", sh.CreateUser)
	e.PUT("/users/:id", sh.UpdateUser)

	e.Logger.Fatal(e.Start(":8000"))
}
