package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Planutim/simple-api/data"
	"github.com/Planutim/simple-api/helpers"
	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var valdator *validator.Validate
var dbHelper helpers.Helper
var h *SimpleHandler

func init() {
	godotenv.Load("../.env")
	valdator = validator.New()
	dbHelper, err := helpers.NewTestHelper()
	if err != nil {
		log.Fatal(err)
	}
	h = NewSimpleHandler(dbHelper, valdator)
}
func TestCreateUser(t *testing.T) {
	e := echo.New()

	users := []*data.User{
		{
			//ok
			Firstname: "General",
			Lastname:  "Kenobi",
			Email:     "GeneralKenobi@mail.ru",
			Age:       26,
		},
		{
			// duplicate email
			Firstname: "General",
			Lastname:  "Kenobi",
			Email:     "GeneralKenobi@mail.ru",
			Age:       26,
		},
		{
			// no firstname
			Firstname: "",
			Lastname:  "Kenobi",
			Email:     "newmail@mail.ru",
			Age:       16,
		},
		{
			// no lastname
			Firstname: "General",
			Lastname:  "",
			Email:     "newmail@mail.ru",
			Age:       25,
		},
		{
			// short name
			Firstname: "Ge",
			Lastname:  "Kenobi",
			Email:     "newmail@mail.ru",
			Age:       25,
		},
	}
	var err error
	h.DBHelper, err = helpers.NewTestHelper()
	if err != nil {
		t.Errorf("Error creating dbhelper! %v", err)
	}

	for i, user := range users {
		userSerialized, err := json.Marshal(&user)
		if err != nil {
			t.Errorf("Marshal error: %v", err)
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
					t.Errorf("Unmarshal error: %v", err)
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
		case 3:
			assert.Equal(t, http.StatusBadRequest, rr.Code)
		case 4:
			assert.Equal(t, http.StatusBadRequest, rr.Code)
		}

	}
}

func TestGetUserById(t *testing.T) {
	e := echo.New()

	user := &data.User{
		Firstname: "General",
		Lastname:  "Kenobi",
		Email:     "starwars@gmail.com",
		Age:       26,
	}
	var err error
	h.DBHelper, err = helpers.NewTestHelper()
	if err != nil {
		t.Errorf("Error creating dbhelper! %v", err)
	}

	userSerialized, err := json.Marshal(user)
	if err != nil {
		t.Errorf("Marshal error: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(userSerialized))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rr := httptest.NewRecorder()

	c := e.NewContext(req, rr)

	h.CreateUser(c)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var userWithID *data.User
	err = json.Unmarshal(rr.Body.Bytes(), &userWithID)
	if err != nil {
		t.Errorf("Unmarshal error: %v", err)
	}

	id := userWithID.ID

	req2 := httptest.NewRequest(http.MethodGet, "/users", nil)
	rr2 := httptest.NewRecorder()

	c2 := e.NewContext(req2, rr2)
	c2.SetParamNames("id")
	c2.SetParamValues(id.String())

	h.ListUser(c2)
	assert.Equal(t, http.StatusOK, rr2.Code)
	assert.Equal(t, rr.Body.String(), rr2.Body.String())

	req3 := httptest.NewRequest(http.MethodGet, "/users", nil)
	rr3 := httptest.NewRecorder()

	c3 := e.NewContext(req3, rr3)

	c3.SetParamNames("id")
	c3.SetParamValues("wrong_uuid")

	h.ListUser(c3)

	assert.Equal(t, http.StatusBadRequest, rr3.Code)
	assert.Equal(t, `{"message":"Wrong id format!"}`, strings.TrimRight(rr3.Body.String(), "\n "))

}

func TestUpdateUserById(t *testing.T) {
	e := echo.New()

	user := &data.User{
		Firstname: "General",
		Lastname:  "Kenobi",
		Email:     "starwars@gmail.com",
		Age:       26,
	}

	var err error
	h.DBHelper, err = helpers.NewTestHelper()
	if err != nil {
		t.Errorf("Error creating dbhelper! %v", err)
	}

	userSerialized, err := json.Marshal(&user)
	if err != nil {
		t.Errorf("Marshal error: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(userSerialized))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rr := httptest.NewRecorder()

	c := e.NewContext(req, rr)
	c.SetParamNames("id")

	h.CreateUser(c)
	assert.Equal(t, http.StatusCreated, rr.Code)
	var userCreated *data.User
	err = json.Unmarshal(rr.Body.Bytes(), &userCreated)
	if err != nil {
		t.Errorf("Unmarshal error: %v", err)
	}
	id := userCreated.ID
	users := []*data.User{
		{
			Firstname: "Emperor",
			Lastname:  "Palpatine",
			Email:     "newstarwars@gmail.com",
			Age:       50,
		},
		{
			Firstname: "Em",
			Lastname:  "Palpatine",
			Email:     "newstarwars@gmail.com",
			Age:       50,
		},
		{
			Firstname: "Emperor",
			Lastname:  "Palpatine",
			Age:       50,
		},
	}

	for i, user := range users {
		userUSerialized, err := json.Marshal(&user)
		if err != nil {
			t.Errorf("Marshal error: %v", err)
		}

		reqPut := httptest.NewRequest(http.MethodPut, "/users", bytes.NewBuffer(userUSerialized))
		reqPut.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rrPut := httptest.NewRecorder()

		cPut := e.NewContext(reqPut, rrPut)
		cPut.SetParamNames("id")
		cPut.SetParamValues(id.String())
		h.UpdateUser(cPut)

		switch i {
		case 0:
			assert.Equal(t, http.StatusOK, rrPut.Code)
		case 1:
			assert.Equal(t, http.StatusBadRequest, rrPut.Code)
		case 2:
			assert.Equal(t, http.StatusBadRequest, rrPut.Code)
		}
	}

}
