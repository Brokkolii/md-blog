package main

import (
	"fmt"
	"io"
	"md-blog-api/database"
	"net/http"
	"text/template"

	"github.com/labstack/echo"

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

func main() {

	database.InitDB()
	defer database.Db.Close()

	e := echo.New()
	e.Renderer = newTemplate()

	// statics
	e.Static("/assets", "static")

	// api endpoints
	e.GET("api/posts", func(c echo.Context) error {
		posts, err := database.GetPosts()
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		c.Response().Header().Set("Cache-Control", "no-cache")
		return c.Render(http.StatusOK, "posts", posts)
	})

	e.POST("api/posts/create", func(c echo.Context) error {
		content := c.FormValue("content")
		post, err := database.CreatePost(content)
		if err != nil {
			fmt.Println("got here", err)
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		c.Response().Header().Set("Cache-Control", "no-cache")
		return c.Render(http.StatusOK, "post", post)
	})

	e.GET("api/posts/:id", func(c echo.Context) error {
		id := c.Param("id")
		post, err := database.GetPost(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		c.Response().Header().Set("Cache-Control", "public, max-age=86400")
		return c.Render(http.StatusOK, "post", post)
	})

	e.GET("/", func(c echo.Context) error {
		c.Response().Header().Set("Cache-Control", "no-cache")
		return c.Render(http.StatusOK, "page", nil)
	})

	e.Start(":12345")
}
