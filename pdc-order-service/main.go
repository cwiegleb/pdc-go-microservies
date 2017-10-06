package main

import (
	"net/http"

	"github.com/cwiegleb/pdc-services/pdc-order-service/handler"
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
	http.ListenAndServe(":9004", r)
}
