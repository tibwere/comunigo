package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	port = "80"
)

type TemplateRegistry struct {
	templates *template.Template
}

func (t *TemplateRegistry) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Renderer = &TemplateRegistry{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}

	e.File("/", "views/index.html")

	e.POST("/login", func(c echo.Context) error {
		username := c.FormValue("username")
		return c.HTML(http.StatusOK, fmt.Sprintf("<h1>Welcome %s</h1>", username))
	})

	e.Logger.Fatal(e.Start(":" + port))
}
