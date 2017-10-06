package main

import (
	"net/http"

	"github.com/cwiegleb/pdc-services/pdc-cashbox-service/handler"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/cashboxes", handler.PostHandler).Methods("POST")
	r.HandleFunc("/cashboxes", handler.GetsHandler).Methods("GET")
	r.HandleFunc("/cashboxes/{id}", handler.GetHandler).Methods("GET")
	r.HandleFunc("/cashboxes/{id}", handler.DeleteHandler).Methods("DELETE")
	r.HandleFunc("/cashboxes/{id}", handler.PutHandler).Methods("PUT")
	http.ListenAndServe(":9002", r)
}
