package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"gotemplate/config"
	handler "gotemplate/handler"
	"gotemplate/logger"
	repo "gotemplate/repo/postgres"
	r "gotemplate/route"
	"gotemplate/supertest"

	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	db               *repo.DB
	router           *handler.Router
	log              *logger.Logger
	err              error
	c                config.Econfig
	validatorService *handler.ValidatorService
)

func validation() error {
	tagToNumber := handler.GetTagToNumberMap()
	errordbMap := handler.GetErrordbMap()
	var err error
	validatorService, err = handler.NewValidatorService(tagToNumber, errordbMap)
	if err != nil {
		return err
	}
	err = validatorService.RegisterCustomValidation("myvalidate", handler.Myvalidate, "must be equal to 10", "CST1")
	if err != nil {
		return err
	}

	err = validatorService.RegisterCustomValidation("hourvalidate", handler.HourValidate, "Pass correct hour", "CST2")
	if err != nil {
		return err
	}
	return nil
}

func setupTest() {
	// Initialize config
	log = logger.New()
	c = config.Load(log)
	log.SetLevel(c.LogLevel())

	
	err = validation()
	if err != nil {
		log.Error("Error initialising validation:", err.Error())

	}
	//validatorService, _ := handler.NewValidatorService(handler.TagToNumber, handler.ErrordbMap)
	// Initialize database
	ctx := context.Background()
	db, err = repo.NewDB(ctx, c)
	if err != nil {
		log.Warn("error in db connection %s", err)
	}

	log.Info("Successfully connected to the database %s", os.Getenv("DB_CONNECTION"))

	// Initialize router
	router, err = getrouter(db, log, validatorService)
	if err != nil {
		log.Error("error in router")
	}

}

func getrouter(db *repo.DB, log *logger.Logger, validatorService *handler.ValidatorService) (router *handler.Router, err error) {

	router, err = r.Routes(db, log, c, validatorService)
	return

}

func TestMain(m *testing.M) {
	setupTest()

	defer db.Close()
	//testdb
	exitCode := m.Run()

	// Clean up test data
	//teardownTestData(db)
	os.Exit(exitCode)
}
func teardownTestData(db *repo.DB) {
	// Clean up test data from the 'users' table
	_, err := db.Exec(context.Background(), `
		DELETE FROM public.users 
	`)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deleting test data: %v\n", err)
		os.Exit(1)
	}

	// Close the database connection pool
	//db.Close()
}

func Positiveassertions(t *testing.T, rr *httptest.ResponseRecorder, expectedStruct interface{}, method string) {
	var dataJSON []byte
	var response handler.Response
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success, fmt.Sprintf("Expected status code %d, got %d", http.StatusOK, rr.Code))
	assert.Equal(t, "Success", response.Message)
	if method != "delete" {
		//log.Debug("came inside, Response data is:", response.Data)
		dataJSON, err = json.Marshal(response.Data)
		assert.NoError(t, err, "Failed to marshal data")
		decoder := json.NewDecoder(bytes.NewReader(dataJSON))
		decoder.DisallowUnknownFields()
		err := decoder.Decode(expectedStruct)
		assert.NoError(t, err, "Failed to unmarshal data")
	}
}

type testerrordbResponse struct {
	Success bool     `json:"success" example:"false"`
	Message []string `json:"message" example:"Error message"`
	Errorno []string `json:"errorno"`
}

func Negativeassertions(t *testing.T, rr *httptest.ResponseRecorder, expectedStruct interface{}, method string) {
	var response testerrordbResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success, fmt.Sprintf("Expected status code %d, got %d", http.StatusOK, rr.Code))
}

func sendAndAssertPostRequest(t *testing.T, test *supertest.Supertest, url string, payload gin.H, expectedStruct interface{}, method string) {

	switch method {
	case "post":
		test.Post(url)
		test.Send(payload)
	case "put":
		test.Put(url)
		test.Send(payload)
	case "delete":
		test.Delete(url)
		test.Send(nil)
	case "get":
		test.Get(url)
		test.Send(nil)

	default:
		test.Get(url)
		test.Send(nil)
	}

	test.Set("Content-Type", "application/json")
	test.End(func(req *http.Request, rr *httptest.ResponseRecorder) {

		assert.Equal(t, http.StatusOK, rr.Code, fmt.Sprintf("Expected status code %d, got %d", http.StatusOK, rr.Code))
		if rr.Code == http.StatusOK {
			Positiveassertions(t, rr, expectedStruct, method)

		}
	})
}

func sendAndAssertErrorRequest(t *testing.T, test *supertest.Supertest, url string, payload gin.H, expectedStruct interface{}, method string, expectedStatusCode int) {
	switch method {
	case "post":
		test.Post(url)
		test.Send(payload)
	case "put":
		test.Put(url)
		test.Send(payload)
	case "delete":
		test.Delete(url)
		test.Send(nil)
	case "get":
		test.Get(url)
		test.Send(nil)
	default:
		test.Get(url)
		test.Send(nil)
	}

	if method == "unsupported" {

	} else {
		test.Set("Content-Type", "application/json")
	}

	test.End(func(req *http.Request, rr *httptest.ResponseRecorder) {
		assert.Equal(t, expectedStatusCode, rr.Code, fmt.Sprintf("Expected status code %d, got %d", expectedStatusCode, rr.Code))

		Negativeassertions(t, rr, expectedStruct, method)

	})
}
