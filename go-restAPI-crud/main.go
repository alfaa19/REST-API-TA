package main

import (
	"log"
	"net/http"

	statscontroller "github.com/alfaa19/go-restapi-crud/controllers/statsController"
	databases "github.com/alfaa19/go-restapi-crud/database"
	"github.com/gorilla/mux"
)

func main() {
	databases.ConnectDatabase()

	r := mux.NewRouter()

	r.HandleFunc("/stats", statscontroller.GetAll).Methods("GET")
	r.HandleFunc("/stats/{id}", statscontroller.GetOne).Methods("GET")
	r.HandleFunc("/stats", statscontroller.Create).Methods("POST")
	r.HandleFunc("/stats/{id}", statscontroller.Update).Methods("PUT")
	r.HandleFunc("/stats/{id}", statscontroller.Delete).Methods("DELETE")

	log.Fatal(http.ListenAndServe("192.168.56.1:8081", r))
}
