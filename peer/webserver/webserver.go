// Package per la gestione della comunicazione con il frontend
// mediante il webserver go Echo
package webserver

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gitlab.com/tibwere/comunigo/peer"
	"gitlab.com/tibwere/comunigo/utilities"
)

// Route che possono essere servite dal webserver
const (
	RouteSing = "/sign"
	RouteSend = "/send"
	RouteList = "/list"
	RouteInfo = "/info"
	RouteRoot = "/"
)

// In ottica OO, oggetto che rappresenta il webserver
type WebServer struct {
	port          uint16
	chatGroupSize uint16
	tos           utilities.TypeOfService
	peerStatus    *peer.Status
	verbose       bool
}

// "Costruttore" dell'oggetto WebServer
func New(exposedPort uint16, size uint16, tos utilities.TypeOfService, verbose bool, status *peer.Status) *WebServer {
	return &WebServer{
		port:          exposedPort,
		chatGroupSize: size,
		tos:           tos,
		peerStatus:    status,
		verbose:       verbose,
	}
}

// "Metodo della classe WebServer" che inizializza
// l'attività di log su file
func (ws *WebServer) initLogger(e *echo.Echo) {
	logFile, err := os.OpenFile(
		fmt.Sprintf("/logs/peer_%v_ws.log", ws.peerStatus.GetExposedIP()),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0666,
	)
	if err != nil {
		log.Fatal(err)
	}
	e.Logger.SetOutput(logFile)
}

// "Metodo della classe WebServer" responsabile di effettuare il setup
// delle impostazioni di echo e di far partire il web server
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

// "Metodo della classe WebServer" handler della route "/info"
func (ws *WebServer) retrieveInfo(c echo.Context) error {
	if ws.peerStatus.NotYetSigned() {
		return c.NoContent(http.StatusForbidden)
	} else {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"Tos":          ws.tos.ToString(),
			"Username":     ws.peerStatus.GetCurrentUsername(),
			"OtherMembers": ws.peerStatus.GetOtherMembers(),
			"Verbose":      ws.verbose,
		})
	}
}

// "Metodo della classe WebServer" handler della route "/list"
func (ws *WebServer) updateMessageList(c echo.Context) error {
	if ws.peerStatus.NotYetSigned() {
		return c.NoContent(http.StatusForbidden)
	} else {
		switch ws.tos {
		case utilities.TOS_CS_SEQUENCER:
			messages, err := ws.peerStatus.GetMessagesSEQ()
			if err != nil {
				return c.NoContent(http.StatusInternalServerError)
			} else {
				return c.JSON(http.StatusOK, messages)
			}
		case utilities.TOS_P2P_SCALAR:
			messages, err := ws.peerStatus.GetMessagesSC()
			if err != nil {
				return c.NoContent(http.StatusInternalServerError)
			} else {
				return c.JSON(http.StatusOK, messages)
			}
		case utilities.TOS_P2P_VECTORIAL:
			messages, err := ws.peerStatus.GetMessagesVC()
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

// "Metodo della classe WebServer" handler della route "/"
func (ws *WebServer) mainPageHandler(c echo.Context) error {
	if ws.peerStatus.NotYetSigned() {
		return c.File("/assets/login.html")
	} else {
		return c.File("/assets/index.html")
	}
}

// "Metodo della classe WebServer" handler della route "/sign"
func (ws *WebServer) signNewUserHandler(c echo.Context) error {
	if ws.peerStatus.NotYetSigned() {
		ws.peerStatus.PushIntoFrontendBackendChannel(c.FormValue("username"))

		result := <-ws.peerStatus.GetFromFrontendBackendChannel()
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

// "Metodo della classe WebServer" handler della route "/send"
func (ws *WebServer) sendMessageHandler(c echo.Context) error {
	if ws.peerStatus.NotYetSigned() {
		return c.NoContent(http.StatusForbidden)
	} else {
		if delayStr := c.FormValue("delay"); delayStr != "" {
			waitBeforeSend(delayStr)
		}
		ws.peerStatus.PushIntoFrontendBackendChannel(c.FormValue("message"))
		return c.NoContent(http.StatusOK)
	}
}

// "Metodo della classe WebServer" che implementa la logica d'attesa
// nel caso in cui oltre al parametro message venga passato dal frontend
// il parametro delay nel formato "min:max" (tipicamente per motivi di test)
func waitBeforeSend(delayStr string) {
	delayBoundaries := strings.Split(delayStr, ":")
	if len(delayBoundaries) != 2 {
		return
	}

	startDelay, errStart := strconv.Atoi(delayBoundaries[0])
	endDelay, errEnd := strconv.Atoi(delayBoundaries[1])
	if errStart == nil && errEnd == nil {
		delay := rand.Intn((endDelay*1000)-(startDelay*1000)+1) + (startDelay * 1000)
		log.Printf("Delay extracted between %v and %v is %v, so waiting %v millisec before send ...", startDelay*1000, endDelay*1000, delay, delay)
		time.Sleep(time.Duration(delay) * time.Millisecond)
	}
}
