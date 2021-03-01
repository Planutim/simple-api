package handlers

import (
	"net/http"

	"github.com/Planutim/simple-api/data"
	"github.com/Planutim/simple-api/db"
	"github.com/go-playground/validator"
	echo "github.com/labstack/echo/v4"
)

//SimpleHandler is a router handler for my app
type SimpleHandler struct {
	DBHelper  *db.Helper
	Validator *validator.Validate
}

//NewSimpleHandler creates an instance of SimpleHandler
func NewSimpleHandler(db *db.Helper, validator *validator.Validate) *SimpleHandler {
	return &SimpleHandler{db, validator}
}

//Home returns  response for GET /
func (sh *SimpleHandler) Home(c echo.Context) error {
	return c.String(200, "This is home page!")
}

//CreateUser creates user for POST /users
func (sh *SimpleHandler) CreateUser(c echo.Context) error {
	var u data.User
	if err := c.Bind(&u); err != nil {
		return err
	}
	err := sh.Validator.Struct(u)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err = sh.DBHelper.DB.Create(&u).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// return c.String(http.StatusCreated, "GOT IT!")
	return c.JSON(http.StatusCreated, u)
}
