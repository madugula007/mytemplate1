package supertest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	//"github.com/gin-gonic/gin"
	handler "gotemplate/handler"
)

/**
* @description -> core funcionality
 */

type Options struct {
	Key string
	Value interface{}
}

type payload struct {
	path   string
	method string
}

type response struct {
	httpResponse *httptest.ResponseRecorder
}

type request struct {
	httpRequest *http.Request
}

type Supertest struct {
	router *handler.Router
	test   *testing.T
	payload
	response
	request
}

/**
* @description -> parent core funcionality
 */

func NewSuperTest(router *handler.Router, test *testing.T) *Supertest {
	return &Supertest{router: router, test: test}
}

/**
* @description -> http client for get request
 */

func (ctx *Supertest) Get(url string) {
	ctx.payload.path = url
	ctx.payload.method = http.MethodGet
}

/**
* @description -> http client for post request
 */

func (ctx *Supertest) Post(url string) {
	ctx.payload.path = url
	ctx.payload.method = http.MethodPost
}

/**
* @description -> http client for delete request
 */

func (ctx *Supertest) Delete(url string) {
	ctx.payload.path = url
	ctx.payload.method = http.MethodDelete
}

/**
* @description -> http client for put request
 */

func (ctx *Supertest) Put(url string) {
	ctx.payload.path = url
	ctx.payload.method = http.MethodPut
}

/**
* @description -> http client for patch request
 */

func (ctx *Supertest) Patch(url string) {
	ctx.payload.path = url
	ctx.payload.method = http.MethodPatch
}

/**
* @description -> http client for head request
 */

func (ctx *Supertest) Head(url string) {
	ctx.payload.path = url
	ctx.payload.method = http.MethodHead
}

/**
* @description -> http client for options request
 */

func (ctx *Supertest) Options(url string) {
	ctx.payload.path = url
	ctx.payload.method = http.MethodOptions
}
