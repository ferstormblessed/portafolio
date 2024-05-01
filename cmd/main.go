package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplates() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

type Data struct {
	Message string
}

func newData() Data {
	return Data{
		Message: "",
	}
}

func getUser(c echo.Context) error {
	// User ID from path `users/:id`
	id := c.Param("id")
	return c.String(http.StatusOK, id)
}

func show(c echo.Context) error {
	// Get team and member from the query string
	team := c.QueryParam("team")
	member := c.QueryParam("member")
	return c.String(http.StatusOK, "team:"+team+", member:"+member)
}

func main() {
	e := echo.New()
	e.Renderer = NewTemplates()

	e.Static("/images", "images")
	e.Static("/css", "css")

	data := newData()
	data.Message = "hello, world"

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", data)
	})

	// e.POST("/users", saveUser)
	e.GET("/users/:id", getUser)
	e.GET("/show", show)
	// e.PUT("/users/:id", updateUser)
	// e.DELETE("/users/:id", deleteUser)

	e.Logger.Fatal(e.Start(":42069"))
}
