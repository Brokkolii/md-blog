package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func getPostEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	htmlPost := `
	<h1>Hello World</h1>
    <p>This is a static HTML snippet.</p>
    `
	w.Write([]byte(htmlPost))
}

func main() {
	router := mux.NewRouter()

	// api endpoints
	router.HandleFunc("/posts/{id}", getPostEndpoint).Methods("GET")

	// static pages
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	http.ListenAndServe(":12345", router)
}
