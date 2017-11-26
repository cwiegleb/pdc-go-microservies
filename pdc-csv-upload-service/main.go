package main

import (
	"net/http"

	"github.com/cwiegleb/pdc-services/pdc-csv-upload-service/handler"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/dealers-upload", handler.PostHandler).Methods("POST")

	allowedHeaders := handlers.AllowedHeaders([]string{"content-type"})
	allowedMethods := handlers.AllowedMethods([]string{"POST", "HEAD"})

	http.ListenAndServe(":9005", handlers.CORS(allowedHeaders, allowedMethods)(r))
}
