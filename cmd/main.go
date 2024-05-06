package main

import (
	"bufio"
	"html/template"
	"io"
	"log"
	"os"

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
	Text []string
}

func newData() Data {
	return Data{
		Text: make([]string, 0),
	}
}

func (d *Data) getContentText(filepath string) {
    file, err := os.Open(filepath)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        d.Text = append(d.Text, scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
}

func main() {
	e := echo.New()
	e.Renderer = NewTemplates()

	e.Static("/images", "images")
	e.Static("/css", "css")

    data := newData()
    filepath := "text/content.md"
    data.getContentText(filepath)

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", data)
	})

	e.Logger.Fatal(e.Start(":42069"))
}
