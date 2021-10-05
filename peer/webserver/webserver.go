package webserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gitlab.com/tibwere/comunigo/peer"
)

const (
	RouteSing = "/sign"
	RouteSend = "/send"
	RouteList = "/list"
	RouteRoot = "/"
)

type WebServer struct {
	port          uint16
	chatGroupSize uint16
	peerStatus    *peer.Status
}

func New(exposedPort uint16, size uint16, status *peer.Status) *WebServer {
	return &WebServer{
		port:          exposedPort,
		chatGroupSize: size,
		peerStatus:    status,
	}
}

func (ws *WebServer) Startup(wg *sync.WaitGroup) {
	defer wg.Done()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static(RouteRoot, "/assets")
	e.GET(RouteRoot, ws.mainPageHandler)
	e.POST(RouteList, ws.updateMessageList)
	e.POST(RouteSing, ws.signNewUserHandler)
	e.POST(RouteSend, ws.sendMessageHandler)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", ws.port)))
}

func sendJSONString(c echo.Context, data map[string]interface{}) error {
	jsondata, err := json.Marshal(data)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	} else {
		return c.JSON(http.StatusOK, string(jsondata))
	}
}

func (ws *WebServer) getListOfOtherUsername() []string {
	var usernames []string
	for _, member := range ws.peerStatus.Members {
		if member.GetUsername() != ws.peerStatus.CurrentUsername {
			usernames = append(usernames, member.GetUsername())
		}
	}

	return usernames
}

func (ws *WebServer) updateMessageList(c echo.Context) error {

	if ws.peerStatus.CurrentUsername == "" {
		return c.NoContent(http.StatusForbidden)
	} else {
		nextIndex, err := strconv.ParseUint(c.FormValue("next"), 10, 16)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		messages, err := peer.GetMessages(ws.peerStatus.Datastore, ws.peerStatus.CurrentUsername, nextIndex)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		return sendJSONString(c, map[string]interface{}{
			"MessageList": messages,
		})
	}
}

func (ws *WebServer) mainPageHandler(c echo.Context) error {
	if ws.peerStatus.CurrentUsername == "" {
		return c.File("/assets/login.html")
	} else {
		return c.File("/assets/index.html")
	}
}

func (ws *WebServer) signNewUserHandler(c echo.Context) error {
	if ws.peerStatus.CurrentUsername == "" {
		ws.peerStatus.UsernameCh <- c.FormValue("username")

		select {
		case <-ws.peerStatus.DoneCh:
			return sendJSONString(c, map[string]interface{}{
				"Username": ws.peerStatus.CurrentUsername,
				"Members":  ws.getListOfOtherUsername(),
			})

		case <-ws.peerStatus.InvalidCh:
			return sendJSONString(c, map[string]interface{}{
				"IsError":      true,
				"ErrorMessage": "Username already in use, please retry with another one!",
			})
		}
	} else {
		return c.NoContent(http.StatusForbidden)
	}
}

func (ws *WebServer) sendMessageHandler(c echo.Context) error {
	if ws.peerStatus.CurrentUsername == "" {
		return c.NoContent(http.StatusForbidden)
	} else {
		ws.peerStatus.RawMessageCh <- c.FormValue("message")
		return c.NoContent(http.StatusOK)
	}
}
