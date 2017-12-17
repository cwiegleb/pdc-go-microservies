package main

import (
	"net/http"

	"github.com/cwiegleb/pdc-services/pdc-article-service/handler"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/articles//dealers/{id}", handler.PostHandler).Methods("POST")
	r.HandleFunc("/articles/dealers/{id}", handler.GetsHandler).Methods("GET")
	r.HandleFunc("/articles/{article-id}/dealers/{id}", handler.PutHandler).Methods("PUT")
	r.HandleFunc("/articles/{article-id}/dealers/{id}", handler.GetHandler).Methods("GET")
	r.HandleFunc("/articles/{article-id}/dealers/{id}", handler.DeleteHandler).Methods("DELETE")

	allowedHeaders := handlers.AllowedHeaders([]string{"content-type"})
	exposedHeaders := handlers.ExposedHeaders([]string{"Location"})
	allowedMethods := handlers.AllowedMethods([]string{"PUT", "POST", "GET", "HEAD"})

	http.ListenAndServe(":9001", handlers.CORS(allowedHeaders, exposedHeaders, allowedMethods)(r))
}
