package main

import (
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/cwiegleb/pdc-services/pdc-cashbox-service/handler"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/cashboxes", handler.PostHandler).Methods("POST")
	r.HandleFunc("/cashboxes", handler.GetsHandler).Methods("GET")
	r.HandleFunc("/cashboxes/{id}", handler.GetHandler).Methods("GET")
	r.HandleFunc("/cashboxes/{id}", handler.DeleteHandler).Methods("DELETE")
	r.HandleFunc("/cashboxes/{id}", handler.PutHandler).Methods("PUT")
	http.ListenAndServe(":8080", r)
}
