package main

import (
	"net/http"

	"github.com/cwiegleb/pdc-services/pdc-order-service/handler"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/cashboxes/{id}/orders", handler.PostHandler).Methods("POST")
	r.HandleFunc("/cashboxes/{id}/orders", handler.GetsHandler).Methods("GET")
	r.HandleFunc("/cashboxes/{id}/orders/{order-id}", handler.GetHandler).Methods("GET")
	r.HandleFunc("/cashboxes/{id}/orders/{order-id}", handler.DeleteHandler).Methods("DELETE")
	r.HandleFunc("/cashboxes/{id}/orders/{order-id}", handler.PutHandler).Methods("PUT")

	allowedHeaders := handlers.AllowedHeaders([]string{"content-type"})
	exposedHeaders := handlers.ExposedHeaders([]string{"Location"})
	allowedMethods := handlers.AllowedMethods([]string{"PUT", "POST", "GET", "HEAD"})

	http.ListenAndServe(":9004", handlers.CORS(allowedHeaders, exposedHeaders, allowedMethods)(r))
}
