package main

import (
	"dynomo/Util"
	_ "fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {

	err := Util.InitDynamoDB()
	if err != nil {
		log.Fatal(err)
		return
	}

	r := mux.NewRouter()
	r.HandleFunc("/picus/list", Util.GetList).Methods("GET")
	r.HandleFunc("/picus/put", Util.PutItem).Methods("POST")
	r.HandleFunc("/picus/get/{key}", Util.GetItem).Methods("GET")
	r.HandleFunc("/picus/{key}", Util.DeleteItem).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}
