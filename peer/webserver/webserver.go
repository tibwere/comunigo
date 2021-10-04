package webserver

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"text/template"

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

type loginTemplate struct {
	IsError      bool
	ErrorMessage string
}

type indexTemplate struct {
	Username string
	Members  []string
}

type WebServer struct {
	port          uint16
	chatGroupSize uint16
	peerStatus    *peer.Status
	templates     *template.Template
}

func New(exposedPort uint16, size uint16, status *peer.Status) *WebServer {
	return &WebServer{
		port:          exposedPort,
		chatGroupSize: size,
		peerStatus:    status,
		templates:     template.Must(template.ParseGlob("/assets/*.html")),
	}
}

func (ws *WebServer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return ws.templates.ExecuteTemplate(w, name, data)
}

func (ws *WebServer) Startup(wg *sync.WaitGroup) {
	defer wg.Done()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Renderer = ws

	e.Static(RouteRoot, "/assets")
	e.GET(RouteRoot, ws.mainPageHandler)
	e.POST(RouteSing, ws.signNewUserHandler)
	e.POST(RouteSend, ws.sendMessageHandler)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", ws.port)))
}

func (ws *WebServer) getListOfUsername() []string {
	var usernames []string
	for _, member := range ws.peerStatus.Members {
		usernames = append(usernames, member.GetUsername())
	}

	return usernames
}

func (ws *WebServer) mainPageHandler(c echo.Context) error {
	if ws.peerStatus.CurrentUsername == "" {
		return c.Render(http.StatusOK, "login", loginTemplate{
			IsError:      false,
			ErrorMessage: "",
		})
	} else {
		return c.File("/assets/index.html")
	}
}

func (ws *WebServer) signNewUserHandler(c echo.Context) error {
	if ws.peerStatus.CurrentUsername == "" {
		ws.peerStatus.UsernameCh <- c.FormValue("username")

		select {
		case <-ws.peerStatus.DoneCh:
			return c.Render(http.StatusOK, "index", indexTemplate{
				Username: ws.peerStatus.CurrentUsername,
				Members:  ws.getListOfUsername(),
			})

		case <-ws.peerStatus.InvalidCh:
			return c.Render(http.StatusOK, "login", loginTemplate{
				IsError:      true,
				ErrorMessage: "Username already in use!",
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
		return c.HTML(http.StatusOK, "Message sent!")
	}
}
