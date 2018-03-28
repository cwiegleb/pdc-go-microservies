package main

import (
	"net/http"

	"github.com/cwiegleb/pdc-services/pdc-dealer-service/handler"
	"github.com/gorilla/handlers"
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
	r.HandleFunc("/dealers/{id}/invoice", handler.GetInvoiceHandler).Methods("GET")
	r.HandleFunc("/dealers-invoices", handler.GetsInvoiceHandler).Methods("GET")
	r.HandleFunc("/dealers-transfer", handler.GetsHibiscusTransferHandler).Methods("GET")

	allowedHeaders := handlers.AllowedHeaders([]string{"content-type"})
	exposedHeaders := handlers.ExposedHeaders([]string{"Location"})
	allowedMethods := handlers.AllowedMethods([]string{"PUT", "POST", "GET", "HEAD"})

	http.ListenAndServe(":9006", handlers.CORS(allowedHeaders, exposedHeaders, allowedMethods)(r))
}
