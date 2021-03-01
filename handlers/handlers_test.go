package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Planutim/simple-api/data"
	"github.com/Planutim/simple-api/db"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	e := echo.New()

	users := []*data.User{
		{
			Firstname: "Artem",
			Lastname:  "Grafkin",
			Email:     "artemgrafkin@mail.ru",
			Age:       26,
		},
		{
			Firstname: "Artem",
			Lastname:  "Grafkin",
			Email:     "artemgrafkin@mail.ru",
			Age:       26,
		},
		{
			Firstname: "",
			Lastname:  "Grafkin",
			Email:     "newmail@mail.ru",
			Age:       16,
		},
	}
	validator := validator.New()
	dbHelper, err := db.NewTestHelper()
	if err != nil {
		t.Errorf("Could not create dbhelper: %s", err.Error())
	}
	h := NewSimpleHandler(dbHelper, validator)

	for i, user := range users {
		userSerialized, err := json.Marshal(&user)
		if err != nil {
			t.Errorf("Marshal error: %s", err.Error())
		}

		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(userSerialized))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rr := httptest.NewRecorder()

		c := e.NewContext(req, rr)

		err = h.CreateUser(c)
		switch i {
		case 0:
			if assert.NoError(t, err) {
				var userCreated data.User
				err = json.Unmarshal(rr.Body.Bytes(), &userCreated)
				if err != nil {
					t.Errorf("Unmarshal error: %s", err.Error())
				}
				assert.Equal(t, http.StatusCreated, rr.Code)
				assert.Equal(t, userCreated.Firstname, user.Firstname)
				assert.Equal(t, userCreated.Lastname, user.Lastname)
				assert.Equal(t, userCreated.Email, user.Email)
				assert.Equal(t, userCreated.Age, user.Age)
			}
		case 1:
			assert.Equal(t, http.StatusInternalServerError, rr.Code)
		case 2:
			assert.Equal(t, http.StatusBadRequest, rr.Code)
		}

	}
}
