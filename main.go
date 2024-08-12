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

	// api

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

	// content

	e.GET("content/blog", func(c echo.Context) error {
		// get posts from db
		posts, err := database.GetPosts()
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		c.Response().Header().Set("Cache-Control", "no-cache")
		return c.Render(http.StatusOK, "blog.content", posts)
	})

	e.GET("content/post/create", func(c echo.Context) error {
		c.Response().Header().Set("Cache-Control", "public, max-age=6400")
		return c.Render(http.StatusOK, "create-post.content", nil)
	})

	e.GET("content/post/:id", func(c echo.Context) error {
		id := c.Param("id")
		post, err := database.GetPost(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		c.Response().Header().Set("Cache-Control", "public, max-age=6400")
		return c.Render(http.StatusOK, "post.content", post)
	})

	e.GET("content/home", func(c echo.Context) error {
		c.Response().Header().Set("Cache-Control", "public, max-age=6400")
		return c.Render(http.StatusOK, "home.content", nil)
	})

	// pages

	e.GET("/blog", func(c echo.Context) error {
		// get posts from db
		posts, err := database.GetPosts()
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		c.Response().Header().Set("Cache-Control", "no-cache")
		return c.Render(http.StatusOK, "blog.page", posts)
	})

	e.GET("/post/create", func(c echo.Context) error {
		c.Response().Header().Set("Cache-Control", "public, max-age=6400")
		return c.Render(http.StatusOK, "create-post.page", nil)
	})

	e.GET("/post/:id", func(c echo.Context) error {
		id := c.Param("id")
		post, err := database.GetPost(id)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		c.Response().Header().Set("Cache-Control", "public, max-age=6400")
		return c.Render(http.StatusOK, "post.page", post)
	})

	e.GET("/home", func(c echo.Context) error {
		c.Response().Header().Set("Cache-Control", "public, max-age=6400")
		return c.Render(http.StatusOK, "home.page", nil)
	})

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusPermanentRedirect, "/home")
	})

	// statics
	e.Static("/assets", "static")

	e.Start(":12345")
}
