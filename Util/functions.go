package Util

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

var Svc *dynamodb.DynamoDB

type Item struct {
	ID   string `json:"id"`
	Data string `json:"data"`
}

func GetList(w http.ResponseWriter, r *http.Request) {
	input := &dynamodb.ScanInput{
		TableName: aws.String("picus"),
	}

	result, err := Svc.Scan(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	items := []Item{}
	for _, i := range result.Items {
		item := Item{
			ID:   *i["id"].S,
			Data: *i["data"].S,
		}
		items = append(items, item)
	}
	/*	fmt.Println(items)
	 */
	json.NewEncoder(w).Encode(items)
}

func InitDynamoDB() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return err
	}
	awsSession, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), ""),
		Region:      aws.String(os.Getenv("AWS_REGION")),
	})
	if err != nil {
		log.Fatal(err)
		return err
	}
	Svc = dynamodb.New(awsSession)
	return nil

}
func PutItem(w http.ResponseWriter, r *http.Request) {
	var item Item
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("DB_NAME")),
		Item: map[string]*dynamodb.AttributeValue{
			"id":   {S: aws.String(item.ID)},
			"data": {S: aws.String(item.Data)},
		},
	}
	fmt.Println(input)

	_, err = Svc.PutItem(input)
	//temp,err:=Svc.PutItem(input) for debug
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	/*fmt.Println(temp)
	fmt.Println(temp.Attributes)*/

	w.WriteHeader(http.StatusOK)
}

func GetItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	input := &dynamodb.GetItemInput{
		TableName: aws.String("picus"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(key)},
		},
	}

	result, err := Svc.GetItem(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.Item == nil {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	item := Item{
		ID:   *result.Item["id"].S,
		Data: *result.Item["data"].S,
	}

	json.NewEncoder(w).Encode(item)
}
