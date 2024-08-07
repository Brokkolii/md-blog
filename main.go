package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"

	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB() {
	var err error
	connStr := "user=mdblog password=mdblogdb dbname=mdblog sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to the database!")
}

func getPostEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Cache-Control", "public, max-age=10")
	params := mux.Vars(r)
	id := params["id"]

	// Static Markdown text
	var markdownContent string
	err := db.QueryRow("SELECT content FROM posts WHERE id = $1", id).Scan(&markdownContent)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(markdownContent)

	w.Write(postToHtml([]byte(markdownContent)))
}

func getPostsEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	rows, err := db.Query("SELECT content FROM posts")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	htmlContent := []byte("")
	for rows.Next() {
		var markdownContent string
		err := rows.Scan(&markdownContent)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		htmlContent = append(htmlContent, postToHtml([]byte(markdownContent))...)
	}

	w.Write(htmlContent)
}

func postToHtml(content []byte) []byte {
	unsafeHtmlContent := blackfriday.Run(content)
	saveHtmlContent := bluemonday.UGCPolicy().SanitizeBytes(unsafeHtmlContent)

	prefix := []byte("<div class='post'>")
	postfix := []byte("</div>")
	html := append(append(prefix, saveHtmlContent...), postfix...)

	return html
}

func main() {

	initDB()
	defer db.Close()

	router := mux.NewRouter()

	// api endpoints
	router.HandleFunc("/posts", getPostsEndpoint).Methods("GET")
	router.HandleFunc("/posts/{id}", getPostEndpoint).Methods("GET")

	// static pages
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	http.ListenAndServe(":12345", router)
}
