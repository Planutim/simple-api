package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/Planutim/simple-api/data"
	"github.com/Planutim/simple-api/helpers"
	"github.com/Planutim/simple-api/jsonerrors"
	"github.com/go-playground/validator"
	"github.com/gofrs/uuid"
	echo "github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

//SimpleHandler is a router handler for my app
type SimpleHandler struct {
	DBHelper  *helpers.Helper
	Validator *validator.Validate
}

//NewSimpleHandler creates an instance of SimpleHandler
func NewSimpleHandler(db *helpers.Helper, validator *validator.Validate) *SimpleHandler {
	return &SimpleHandler{db, validator}
}

//CreateUser creates user for POST /users
func (sh *SimpleHandler) CreateUser(c echo.Context) error {
	var u data.User
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusBadRequest, &jsonerrors.GenericError{
			Message: "Invalid json body!",
		})
	}
	err := sh.Validator.Struct(u)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &jsonerrors.GenericError{
			Message: "Wrong parameters in request!",
		})
	}

	user, err := sh.DBHelper.CreateUser(&u)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &jsonerrors.GenericError{
			Message: "Could not create user!",
		})
	}

	// return c.String(http.StatusCreated, "GOT IT!")
	return c.JSON(http.StatusCreated, user)
}

//ListUsers returns JSON response for GET /users
func (sh *SimpleHandler) ListUsers(c echo.Context) error {
	users, err := sh.DBHelper.GetUsers()
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) == false {
		return c.JSON(http.StatusInternalServerError, &jsonerrors.GenericError{
			Message: "Could not get users!",
		})
	}
	return c.JSON(http.StatusOK, users)
}

//ListUser returns JSON response for GET /users/:id
func (sh *SimpleHandler) ListUser(c echo.Context) error {
	uuidStr := c.Param("id")
	log.Println(uuidStr)
	if _, err := uuid.FromString(uuidStr); err != nil {
		return c.JSON(http.StatusBadRequest, &jsonerrors.GenericError{
			Message: "Wrong id format!",
		})
	}
	user, err := sh.DBHelper.ListUser(uuidStr)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &jsonerrors.GenericError{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, user)
}

//UpdateUser updates user for PUT /users/:id
func (sh *SimpleHandler) UpdateUser(c echo.Context) error {
	var user data.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, &jsonerrors.GenericError{
			Message: "Invalid json body!",
		})
	}
	userID, err := uuid.FromString(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, &jsonerrors.GenericError{
			Message: "Wrong id format!",
		})
	}
	user.ID = userID
	err = sh.Validator.Struct(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &jsonerrors.GenericError{
			Message: "Wrong parameters in request!",
		})
	}
	userUpdated, err := sh.DBHelper.UpdateUser(&user)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, &jsonerrors.GenericError{
				Message: "No user with such id!",
			})
		}
		return c.JSON(http.StatusInternalServerError, &jsonerrors.GenericError{
			Message: "Could not update user!",
		})
	}
	return c.JSON(http.StatusOK, userUpdated)
}
