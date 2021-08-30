package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func IndexHandler(c echo.Context) error {
	// Please note the the second parameter "home.html" is the template name and should
	// be equal to the value stated in the {{ define }} statement in "view/home.html"
	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"name": "INDEX",
		"msg":  "Hello, World!",
	})
}
