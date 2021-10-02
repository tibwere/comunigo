package peer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func StartupWebServer(port uint16, size uint16, channels *PeerChannels, usernamePtr *string, wg *sync.WaitGroup) {
	defer wg.Done()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/sign", func(c echo.Context) error {
		return signNewUserHandler(c, size, usernamePtr, channels.UsernameCh)
	})

	e.POST("/send", func(c echo.Context) error {
		return sendMessage(c, usernamePtr, channels.RawMessageCh)
	})
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", port)))
}

func signNewUserHandler(c echo.Context, size uint16, usernamePtr *string, usernameCh chan string) error {
	if *usernamePtr == "" {
		var i uint16
		var members []string

		usernameCh <- c.FormValue("username")

		for i = 0; i < size; i++ {
			members = append(members, <-usernameCh)
		}

		jsonMembers, err := json.Marshal(members)
		if err != nil {
			return c.NoContent(http.StatusForbidden)
		} else {
			return c.JSON(http.StatusOK, jsonMembers)
		}

	} else {
		return c.NoContent(http.StatusForbidden)
	}
}

func sendMessage(c echo.Context, usernamePtr *string, messageCh chan string) error {
	if *usernamePtr == "" {
		return c.NoContent(http.StatusForbidden)
	} else {
		messageCh <- c.FormValue("message")
		return c.HTML(http.StatusOK, "Message sent!")
	}
}
