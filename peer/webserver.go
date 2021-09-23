package peer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var usernameReceived bool

func StartupWebServer(port uint16, size uint16, usernameCh chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/sign", func(c echo.Context) error {
		return signNewUserHandler(c, size, usernameCh)
	})
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", port)))
}

func signNewUserHandler(c echo.Context, size uint16, usernameCh chan string) error {
	if !usernameReceived {
		usernameReceived = true
		usernameCh <- c.FormValue("username")
		var i uint16
		var members []string
		for i = 0; i < size; i++ {
			members = append(members, <-usernameCh)
		}

		jsonMembers, err := json.Marshal(members)
		fmt.Println(jsonMembers)
		if err != nil {
			return c.NoContent(http.StatusForbidden)
		}

		return c.JSON(http.StatusOK, jsonMembers)
	} else {
		return c.NoContent(http.StatusForbidden)
	}
}
