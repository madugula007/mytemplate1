package tests

import (
	"gotemplate/core/domain"
	handler "gotemplate/handler"
	"gotemplate/supertest"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestUserRegister(t *testing.T) {

	test := supertest.NewSuperTest(router, t)
	url := "/v1/users/"

	payload := gin.H{
		"email":    "99939999te588099@gmail.com",
		"password": "fghjklhjgf",
		"name":     "sawerr",
		"check":    10,
	}
	var userResponse handler.UserResponse
	//var userResponse meta
	sendAndAssertPostRequest(t, test, url, payload, &userResponse, "post")

}

func TestGetUserByID(t *testing.T) {

	test := supertest.NewSuperTest(router, t)
	url := "/v1/users/391"
	payload := gin.H{}
	var userResponse handler.UserResponse
	//var userResponse meta
	sendAndAssertPostRequest(t, test, url, payload, &userResponse, "get")

}

type meta struct {
	Total uint64 `json:"total,omitempty" example:"100"`
	Limit uint64 `json:"limit,omitempty" example:"10"`
	Skip  uint64 `json:"skip,omitempty" example:"0"`
}

func TestGetUsers(t *testing.T) {

	test := supertest.NewSuperTest(router, t)
	url := "/v1/users/?skip=1&limit=5"
	payload := gin.H{}
	type GetAll struct {
		meta  `json:"meta"`
		Users []domain.User `json:"users"`
	}
	var userall GetAll
	sendAndAssertPostRequest(t, test, url, payload, &userall, "get")
}

func TestUserUpdate(t *testing.T) {

	test := supertest.NewSuperTest(router, t)
	url := "/v1/users/391"

	payload := gin.H{
		"email":    "3t10test28890oo@gmail.com",
		"password": "fghjklhjgf",
		"name":     "sawerr4",
		"check":    10,
	}
	var userResponse handler.UserResponse
	sendAndAssertPostRequest(t, test, url, payload, &userResponse, "put")

}

func TestUserDelete(t *testing.T) {

	test := supertest.NewSuperTest(router, t)
	url := "/v1/users/373"
	payload := gin.H{}

	sendAndAssertPostRequest(t, test, url, payload, "delete", "delete")

}

func TestUserInvalidError(t *testing.T) {
	test := supertest.NewSuperTest(router, t)
	url := "/v1/users/invalid/"
	payload := gin.H{}

	sendAndAssertErrorRequest(t, test, url, payload, "", "get", http.StatusNotFound)

}

func TestUserdatanotfoundError(t *testing.T) {
	test := supertest.NewSuperTest(router, t)
	url := "/v1/users/123344/"
	payload := gin.H{}

	sendAndAssertErrorRequest(t, test, url, payload, "", "get", http.StatusNotFound)

}

func TestUserUnsupportedError(t *testing.T) {
	test := supertest.NewSuperTest(router, t)
	url := "/v1/users/123344/"
	payload := gin.H{}

	sendAndAssertErrorRequest(t, test, url, payload, "", "unsupported", http.StatusUnsupportedMediaType)

}

func TestUserUnprocessedError(t *testing.T) {
	test := supertest.NewSuperTest(router, t)
	url := "/v1/users/"
	payload := gin.H{"email": "32t11test288990@gmail.com",
		"password": "fghjklhjgf",
		"name":     "sawerr",
		"check":    "300"}

	sendAndAssertErrorRequest(t, test, url, payload, "", "post", http.StatusUnprocessableEntity)

}

func TestUserdbError(t *testing.T) {
	test := supertest.NewSuperTest(router, t)
	url := "/v1/users/"
	payload := gin.H{"email": "3248r11test288990@gmail.com",
		"password": "fghjklhjgf",
		"name":     "sawerr",
		"check":    10}

	sendAndAssertErrorRequest(t, test, url, payload, "", "post", http.StatusInternalServerError)

}
