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

type Contact struct {
	Name  string
	Email string
}

type Contacts = []Contact

type ContactData struct {
	Contacts Contacts
}

func newContact(name string, email string) Contact {
	return Contact{
		Name:  name,
		Email: email,
	}
}

func createNewData() ContactData {
	return ContactData{
		Contacts: Contacts{
			newContact("John Doe", "jdoe@jyj.in"),
			newContact("Jane Doe", "jndoe@jyj.in"),
		},
	}
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

	// explicit type declaration of contactList
	var contactList ContactData
	contactList = createNewData()
	// OR
	// implicit type declaration of contactList
	// contactList := createNewData()
	e.Renderer = newTemplate()
	e.GET("/", func(c echo.Context) error {
		// http.StatusOK
		return c.Render(200, "index", contactList)
	})

	e.POST("/contacts", func(c echo.Context) error {
		name := c.FormValue("name")
		email := c.FormValue("email")
		contactList.Contacts = append(contactList.Contacts, newContact(name, email))
		return c.Render(200, "index", contactList)
	})

	e.Logger.Fatal(e.Start(":6969"))
}
