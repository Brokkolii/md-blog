package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"text/template"

	"github.com/labstack/echo"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"

	_ "github.com/lib/pq"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Template {
	return &Template{
		templates: template.Must(template.ParseGlob("./views/*.html")),
	}
}

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

type Post struct {
	Content string
}

func newPost() Post {
	return Post{
		Content: "",
	}
}

func getPost(c echo.Context) error {
	id := c.Param("id")

	var md string
	err := db.QueryRow("SELECT content FROM posts WHERE id = $1", id).Scan(&md)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	html := mdToHtml([]byte(md))
	var post = newPost()
	post.Content = string(html)

	return c.Render(http.StatusOK, "post", post)
}

func getPosts(c echo.Context) error {

	rows, err := db.Query("SELECT content FROM posts")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var markdownContent string
		err := rows.Scan(&markdownContent)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		html := mdToHtml([]byte(markdownContent))
		var post = newPost()
		post.Content = string(html)
		posts = append(posts, post)
	}

	return c.Render(http.StatusOK, "posts", posts)
}

func mdToHtml(content []byte) []byte {
	unsafeHtmlContent := blackfriday.Run(content)
	saveHtmlContent := bluemonday.UGCPolicy().SanitizeBytes(unsafeHtmlContent)

	return saveHtmlContent
}

func main() {

	initDB()
	defer db.Close()

	e := echo.New()
	e.Renderer = newTemplate()

	// statics
	e.Static("/assets", "static")

	// api endpoints
	e.GET("/posts", func(c echo.Context) error {
		return getPosts(c)
	})
	e.GET("/posts/:id", func(c echo.Context) error {
		return getPost(c)
	})
	e.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "page", nil)
	})

	e.Start(":12345")
}
