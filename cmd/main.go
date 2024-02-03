package main

import (
	"html/template"
	"io"

	// "fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
	templates *template.Template
}

type Count struct {
	Count int
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// newTemplate returns an instance of Templates, parsing all templates present in the views directory
func newTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

func main() {

	e := echo.New()
	e.Use(middleware.Logger())

	count := Count{Count: 0}
	e.Renderer = newTemplate()
	e.GET("/", func(c echo.Context) error {
		count.Count++
		// http.StatusOK
		return c.Render(200, "index", count)
	})

	e.Logger.Fatal(e.Start(":6969"))
}
