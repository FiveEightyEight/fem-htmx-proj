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

type FormData struct {
	Values map[string]string
	Errors map[string]string
}

type Page struct {
	Data ContactData
	Form FormData
}

func newPage() Page {
	return Page{
		Data: createNewData(),
		Form: newFormData(),
	}
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

func newFormData() FormData {
	return FormData{
		Values: make(map[string]string),
		Errors: make(map[string]string),
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

func (c ContactData) hasEmail(email string) bool {
	for _, contact := range c.Contacts {
		if contact.Email == email {
			return true
		}
	}
	return false
}

func main() {

	e := echo.New()
	e.Use(middleware.Logger())

	// explicit type declaration of contactList
	// var contactList ContactData
	// contactList = createNewData()

	// OR
	// implicit type declaration of contactList
	// contactList := createNewData()
	page := newPage()
	e.Renderer = newTemplate()
	e.GET("/", func(c echo.Context) error {
		// http.StatusOK
		return c.Render(200, "index", page)
	})

	e.POST("/contacts", func(c echo.Context) error {
		name := c.FormValue("name")
		email := c.FormValue("email")

		if page.Data.hasEmail(email) {
			formData := newFormData()
			formData.Values["name"] = name
			formData.Values["email"] = email
			formData.Errors["email"] = "Email already in use"
			return c.Render(422, "form", formData)
		}
		contact := newContact(name, email)
		page.Data.Contacts = append(page.Data.Contacts, contact)

		c.Render(200, "form", newFormData())
		return c.Render(200, "oob-contact", contact)
	})

	e.Logger.Fatal(e.Start(":6969"))
}
