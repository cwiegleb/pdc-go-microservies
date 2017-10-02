package main

import (
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/cwiegleb/pdc-services/pdc-article-service/handler"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/dealers/{id}/articles", handler.PostHandler).Methods("POST")
	r.HandleFunc("/dealers/{id}/articles", handler.GetsHandler).Methods("GET")
	r.HandleFunc("/dealers/{id}/articles/{article-id}", handler.PutHandler).Methods("PUT")
	r.HandleFunc("/dealers/{id}/articles/{article-id}", handler.GetHandler).Methods("GET")
	r.HandleFunc("/dealers/{id}/articles/{article-id}", handler.DeleteHandler).Methods("DELETE")
	http.ListenAndServe(":8082", r)
}
