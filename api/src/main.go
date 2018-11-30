package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	. "app/dao"
	. "app/controllers"
)

var dao = PersonsDAO{}

func init() {
	dao.Server = os.Getenv("SERVER")
	dao.Database = os.Getenv("DATABASE")
	dao.Connect()
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/persons", AllPersonsEndPoint).Methods("GET")
	r.HandleFunc("/persons", CreatePersonEndPoint).Methods("POST")
	r.HandleFunc("/persons", DeleteAllEndPoint).Methods("DELETE")
	r.HandleFunc("/persons/{id}", FindPersonEndpoint).Methods("GET")
	r.HandleFunc("/persons/{id}", UpdatePersonEndPoint).Methods("PUT")
	r.HandleFunc("/persons/{id}", UpdateValueOfPersonEndPoint).Methods("PATCH")
	r.HandleFunc("/persons/{id}", DeletePersonEndPoint).Methods("DELETE")
	r.HandleFunc("/persons/{id}/parents", FindParentsEndPoint).Methods("GET")
	r.HandleFunc("/persons/{id}/children", FindChildrenEndPoint).Methods("GET")
	r.HandleFunc("/persons/{parent_id}/isparentof/{child_id}", RelateParentToChildEndPoint).Methods("PATCH")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}
