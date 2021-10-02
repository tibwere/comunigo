package peer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	POSTUsernameParam = "username"
	POSTMessageParam  = "message"
)

const (
	RouteSing = "/sign"
	RouteSend = "/send"
)

func StartupWebServer(port uint16, size uint16, channels *PeerChannels, usernamePtr *string, wg *sync.WaitGroup) {
	defer wg.Done()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST(RouteSing, func(c echo.Context) error {
		return signNewUserHandler(c, size, usernamePtr, channels.UsernameCh, channels.InvalidCh)
	})

	e.POST(RouteSend, func(c echo.Context) error {
		return sendMessage(c, usernamePtr, channels.RawMessageCh)
	})
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", port)))
}

func comunicateMembers(c echo.Context, members []string) error {
	jsonMembers, err := json.Marshal(members)
	if err != nil {
		return c.NoContent(http.StatusForbidden)
	} else {
		return c.JSON(http.StatusOK, jsonMembers)
	}
}

func notifyInvalidUsername(c echo.Context) error {
	errorMessage := []byte(`{"error":"Username already in use!"}`)
	return c.JSON(http.StatusOK, errorMessage)
}

func signNewUserHandler(c echo.Context, size uint16, usernamePtr *string, usernameCh chan string, invalidCh chan bool) error {
	if *usernamePtr == "" {
		var i uint16
		var members []string

		usernameCh <- c.FormValue(POSTUsernameParam)

		for i = 0; i < size; i++ {
			select {
			case m := <-usernameCh:
				members = append(members, m)
			case <-invalidCh:
				return notifyInvalidUsername(c)
			}
		}

		return comunicateMembers(c, members)

	} else {
		return c.NoContent(http.StatusForbidden)
	}
}

func sendMessage(c echo.Context, usernamePtr *string, messageCh chan string) error {
	if *usernamePtr == "" {
		return c.NoContent(http.StatusForbidden)
	} else {
		messageCh <- c.FormValue(POSTMessageParam)
		return c.HTML(http.StatusOK, "Message sent!")
	}
}
