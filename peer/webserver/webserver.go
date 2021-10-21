package webserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gitlab.com/tibwere/comunigo/peer"
)

const (
	RouteSing = "/sign"
	RouteSend = "/send"
	RouteList = "/list"
	RouteInfo = "/info"
	RouteRoot = "/"
)

type WebServer struct {
	port          uint16
	chatGroupSize uint16
	tos           string
	peerStatus    *peer.Status
}

func New(exposedPort uint16, size uint16, tos string, status *peer.Status) *WebServer {
	return &WebServer{
		port:          exposedPort,
		chatGroupSize: size,
		tos:           tos,
		peerStatus:    status,
	}
}

func (ws *WebServer) initLogger(e *echo.Echo) {
	logFile, err := os.OpenFile(
		fmt.Sprintf("/logs/peer_%v_ws.log", ws.peerStatus.PublicIP),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0666,
	)
	if err != nil {
		log.Fatal(err)
	}
	e.Logger.SetOutput(logFile)
}

func (ws *WebServer) Startup(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time_rfc3339}]: method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Use(middleware.Recover())
	e.HideBanner = true
	ws.initLogger(e)

	go func() {
		<-ctx.Done()
		log.Println("Webserver shutdown")
		timedCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := e.Shutdown(timedCtx); err != nil {
			log.Printf("Unable to shutdown server (%v", err)
		}
	}()

	e.Static(RouteRoot, "/assets")
	e.GET(RouteRoot, ws.mainPageHandler)
	e.GET(RouteList, ws.updateMessageList)
	e.POST(RouteSing, ws.signNewUserHandler)
	e.POST(RouteSend, ws.sendMessageHandler)
	e.GET(RouteInfo, ws.retrieveInfo)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", ws.port)))
}

// func (ws *WebServer) getListOfOtherUsername() []string {
// 	var usernames []string
// 	for _, member := range ws.peerStatus.OtherMembers {
// 		usernames = append(usernames, member.GetUsername())
// 	}

// 	return usernames
// }

func (ws *WebServer) retrieveInfo(c echo.Context) error {
	if ws.peerStatus.CurrentUsername == "" {
		return c.NoContent(http.StatusForbidden)
	} else {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"Tos":          ws.tos,
			"Username":     ws.peerStatus.CurrentUsername,
			"OtherMembers": ws.peerStatus.OtherMembers,
		})
	}
}

func (ws *WebServer) updateMessageList(c echo.Context) error {

	if ws.peerStatus.CurrentUsername == "" {
		return c.NoContent(http.StatusForbidden)
	} else {
		switch ws.tos {
		case "sequencer":
			messages, err := peer.GetMessagesSEQ(ws.peerStatus.Datastore, ws.peerStatus.CurrentUsername)
			if err != nil {
				return c.NoContent(http.StatusInternalServerError)
			} else {
				return c.JSON(http.StatusOK, messages)
			}
		case "scalar":
			messages, err := peer.GetMessagesSC(ws.peerStatus.Datastore, ws.peerStatus.CurrentUsername)
			if err != nil {
				return c.NoContent(http.StatusInternalServerError)
			} else {
				return c.JSON(http.StatusOK, messages)
			}
		case "vectorial":
			messages, err := peer.GetMessagesVC(ws.peerStatus.Datastore, ws.peerStatus.CurrentUsername)
			if err != nil {
				return c.NoContent(http.StatusInternalServerError)
			} else {
				return c.JSON(http.StatusOK, messages)
			}
		default:
			return c.NoContent(http.StatusInternalServerError)
		}
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
		ws.peerStatus.FrontBackCh <- c.FormValue("username")

		result := <-ws.peerStatus.FrontBackCh
		if result == "SUCCESS" {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"Status": result,
			})
		} else {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"Status":  "ERROR",
				"Message": "Username already in use, please retry with another one!",
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
		ws.peerStatus.FrontBackCh <- c.FormValue("message")
		return c.NoContent(http.StatusOK)
	}
}
