package test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"dynomo/Util"
)

func setupTestData() {
	item := map[string]*dynamodb.AttributeValue{
		"id":   {S: aws.String("123")},
		"data": {S: aws.String("test data")},
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("picus"),
		Item:      item,
	}

	_, err := Util.Svc.PutItem(input)
	if err != nil {
		log.Fatalf("Failed to put test data: %v", err)
	}
}

func TestGetList(t *testing.T) {
	req, err := http.NewRequest("GET", "/picus/list", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Util.GetList)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expected := ""

	assert.NotEqual(t, expected, rr.Body.String())
}

func TestPutItem(t *testing.T) {
	item := Util.Item{ID: "1001", Data: "test data"}
	jsonData, err := json.Marshal(item)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/picus/put", bytes.NewBuffer(jsonData))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Util.PutItem)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expected := ""
	assert.Equal(t, expected, rr.Body.String())
}

func TestGetItem(t *testing.T) {
	setupTestData()

	req, err := http.NewRequest("GET", "/picus/get/123", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/picus/get/{key}", Util.GetItem)
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expected := `{"id":"123","data":"test data"}`
	actual := strings.TrimSpace(rr.Body.String())
	assert.Equal(t, expected, actual)
}
func TestDeleteItem(t *testing.T) {
	setupTestData()

	req, err := http.NewRequest("DELETE", "/picus/123", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	router := mux.NewRouter()

	router.HandleFunc("/picus/{key}", Util.DeleteItem).Methods("DELETE")

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expected := "Item deleted successfully"
	actual := strings.TrimSpace(rr.Body.String())

	assert.Equal(t, expected, actual)
}

func TestMain(m *testing.M) {
	Util.InitDynamoDB()

	m.Run()
}
