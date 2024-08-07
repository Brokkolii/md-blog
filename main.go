package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

func getPostEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	// Static Markdown text
	markdownText := "# Hello World\n## Second Hello World\n\nThis is a static markdown text.\n\n- Item 1\n- Item 2\n- Item 3\n\n```go\nfmt.Println(\"Hello, world!\")\n```"

	// Convert Markdown to HTML
	unsafeHtmlContent := blackfriday.Run([]byte(markdownText))
	saveHtmlContent := bluemonday.UGCPolicy().SanitizeBytes(unsafeHtmlContent)

	w.Write(saveHtmlContent)
}

func main() {
	router := mux.NewRouter()

	// api endpoints
	router.HandleFunc("/posts/{id}", getPostEndpoint).Methods("GET")

	// static pages
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	http.ListenAndServe(":12345", router)
}
