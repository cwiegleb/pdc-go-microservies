package main

import (
	"net/http"

	"github.com/cwiegleb/pdc-services/pdc-cashbox-service/handler"
	"github.com/gorilla/handlers"
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
	r.HandleFunc("/cashboxes/{id}/accounting", handler.GetAccountingHandler).Methods("GET")

	allowedHeaders := handlers.AllowedHeaders([]string{"content-type"})
	exposedHeaders := handlers.ExposedHeaders([]string{"Location"})
	allowedMethods := handlers.AllowedMethods([]string{"PUT", "POST", "GET", "HEAD"})

	http.ListenAndServe(":9002", handlers.CORS(allowedHeaders, exposedHeaders, allowedMethods)(r))
}
