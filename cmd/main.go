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

type Text struct {
	Title  string
	Middle []string
	Bottom string
}

func newText() Text {
	return Text{
		Title:  "",
		Middle: make([]string, 0),
		Bottom: "",
	}
}

type Data struct {
	Text Text
}

func newData() Data {
	return Data{
		Text: newText(),
	}
}

func (d *Data) getContentText(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// read all the lines of file
	tmp_text := make([]string, 0)
	for scanner.Scan() {
		tmp_text = append(tmp_text, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	d.Text.Title = tmp_text[0]
	d.Text.Bottom = tmp_text[len(tmp_text)-1]
	for i := 1; i < len(tmp_text)-1; i++ {
		d.Text.Middle = append(d.Text.Middle, tmp_text[i])
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
