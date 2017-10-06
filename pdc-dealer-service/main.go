package main

import (
	"net/http"

	"github.com/cwiegleb/pdc-services/pdc-dealer-service/handler"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/dealers", handler.PostHandler).Methods("POST")
	r.HandleFunc("/dealers/{id}", handler.GetHandler).Methods("GET")
	r.HandleFunc("/dealers", handler.GetsHandler).Methods("GET")
	r.HandleFunc("/dealers/{id}", handler.DeleteHandler).Methods("DELETE")
	r.HandleFunc("/dealers/{id}", handler.PutHandler).Methods("PUT")
	http.ListenAndServe(":9003", r)
}
