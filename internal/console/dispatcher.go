package console

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/jmartin82/mmock/v3/internal/config"
	"github.com/jmartin82/mmock/v3/internal/config/logger"
	"github.com/jmartin82/mmock/v3/internal/statistics"
	"github.com/jmartin82/mmock/v3/pkg/match"
	"github.com/jmartin82/mmock/v3/pkg/mock"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

//go:embed all:ui
var assetFS embed.FS

var log = logger.Log

var pagePattern = regexp.MustCompile(`^[1-9]([0-9]+)?$`)

// ErrInvalidPage the page parameters is invalid
var ErrInvalidPage = errors.New("Invalid page")

type ActionResponse struct {
	Result string `json:"result"`
}

// Dispatcher is the http console server.
type Dispatcher struct {
	IP             string
	Port           int
	ResultsPerPage int
	MatchSpy       match.TransactionSpier
	Scenario       match.ScenearioStorer
	Mapping        config.Mapping
	Mlog           chan match.Transaction
	clients        []*websocket.Conn
}

func (di *Dispatcher) removeClient(i int) {
	copy(di.clients[i:], di.clients[i+1:])
	di.clients[len(di.clients)-1] = nil
	di.clients = di.clients[:len(di.clients)-1]
}

func (di *Dispatcher) addClient(ws *websocket.Conn) {
	di.clients = append(di.clients, ws)
}

func (di *Dispatcher) logFanOut() {
	for match := range di.Mlog {
		for i, c := range di.clients {
			if c != nil {
				if err := websocket.JSON.Send(c, match); err != nil {
					di.removeClient(i)
				}
			}
		}
	}
}

// Start initiates the http console.
func (di *Dispatcher) Start() {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	//WS
	di.clients = []*websocket.Conn{}
	e.GET("/echo", di.webSocketHandler)

	//HTTP
	sub, _ := fs.Sub(assetFS, "ui")
	statics := http.FileServer(http.FS(sub))
	e.GET("/assets/*.js", echo.WrapHandler(statics))
	e.GET("/assets/*.css", echo.WrapHandler(statics))
	e.GET("/favicon.ico", echo.WrapHandler(statics))
	e.GET("/swagger.json", echo.WrapHandler(statics))
	e.GET("/mapping", di.consoleHandler)
	e.GET("/about", di.consoleHandler)
	e.GET("/", di.consoleHandler)

	//verification
	e.GET("/api/request/reset", di.requestResetHandler)
	e.POST("/api/request/verify", di.requestVerifyHandler)
	e.POST("/api/request/reset_match", di.resetMatchHandler)
	e.GET("/api/request/all", di.requestAllHandler)
	e.GET("/api/request/all/:page", di.requestAllPagedHandler)
	e.GET("/api/request/matched", di.requestMatchedHandler)
	e.GET("/api/request/unmatched", di.requestUnMatchedHandler)
	e.GET("/api/scenarios/reset_all", di.scenariosResetHandler)
	e.GET("/api/scenarios", di.scenariosListHandler)
	e.PUT("/api/scenarios/set/:scenario/:state", di.scenariosSetHandler)
	e.PUT("/api/scenarios/pause", di.scenariosPauseHandler)
	e.PUT("/api/scenarios/unpause", di.scenariosUnpauseHandler)

	//mapping
	e.GET("/api/mapping", di.mappingListHandler)
	e.GET("/api/mapping/*", di.mappingGetHandler)
	e.POST("/api/mapping/*", di.mappingCreateHandler)
	e.PUT("/api/mapping/*", di.mappingUpdateHandler)
	e.DELETE("/api/mapping/*", di.mappingDeleteHandler)

	//POST api/mapping (all)

	go di.logFanOut()

	addr := fmt.Sprintf("%s:%d", di.IP, di.Port)
	e.Logger.Fatal(e.Start(addr))
}

// CONSOLE
func (di *Dispatcher) consoleHandler(c echo.Context) error {
	statistics.TrackConsoleRequest()
	tmpl, _ := assetFS.ReadFile("ui/index.html")
	content := string(tmpl)
	return c.HTML(http.StatusOK, content)
}

func (di *Dispatcher) webSocketHandler(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		di.addClient(ws)
		defer ws.Close()
		//block
		var message string
		websocket.Message.Receive(ws, &message)

	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

func (di *Dispatcher) getMappingUri(path string) string {
	root := "/api/mapping/"
	return strings.TrimPrefix(path, root)
}

// API REQUEST
func (di *Dispatcher) mappingListHandler(c echo.Context) (err error) {
	mocks := di.Mapping.List()
	return c.JSON(http.StatusOK, mocks)
}

func (di *Dispatcher) mappingGetHandler(c echo.Context) (err error) {

	URI := di.getMappingUri(c.Request().URL.Path)
	mock := mock.Definition{}
	ok := false
	if mock, ok = di.Mapping.Get(URI); !ok {
		ar := &ActionResponse{
			Result: "not_found",
		}
		return c.JSON(http.StatusNotFound, ar)
	}

	return c.JSON(http.StatusOK, mock)

}

func (di *Dispatcher) mappingDeleteHandler(c echo.Context) (err error) {

	URI := di.getMappingUri(c.Request().URL.Path)
	ok := false
	if _, ok = di.Mapping.Get(URI); !ok {
		ar := &ActionResponse{
			Result: "not_found",
		}
		return c.JSON(http.StatusNotFound, ar)
	}

	if err = di.Mapping.Delete(URI); err != nil {
		return err
	}
	ar := &ActionResponse{
		Result: "deleted",
	}
	return c.JSON(http.StatusOK, ar)

}

func (di *Dispatcher) mappingCreateHandler(c echo.Context) (err error) {

	mock := &mock.Definition{}
	URI := di.getMappingUri(c.Request().URL.Path)

	if _, ok := di.Mapping.Get(URI); ok {
		ar := &ActionResponse{
			Result: "already_exists",
		}
		return c.JSON(http.StatusConflict, ar)
	}

	if err = c.Bind(mock); err != nil {
		ar := &ActionResponse{
			Result: fmt.Sprintf("invalid_mock_definition: %s", err),
		}
		return c.JSON(http.StatusBadRequest, ar)
	}

	err = di.Mapping.Set(URI, *mock)
	if err != nil {
		return
	}

	ar := &ActionResponse{
		Result: "created",
	}
	return c.JSON(http.StatusCreated, ar)

}

func (di *Dispatcher) mappingUpdateHandler(c echo.Context) (err error) {

	mock := &mock.Definition{}
	URI := di.getMappingUri(c.Request().URL.Path)

	if _, ok := di.Mapping.Get(URI); !ok {
		ar := &ActionResponse{
			Result: "not_found",
		}
		return c.JSON(http.StatusNotFound, ar)
	}

	if err = c.Bind(mock); err != nil {
		ar := &ActionResponse{
			Result: fmt.Sprintf("invalid_mock_definition: %s", err),
		}
		return c.JSON(http.StatusBadRequest, ar)
	}

	err = di.Mapping.Set(URI, *mock)
	if err != nil {
		return
	}

	ar := &ActionResponse{
		Result: "updated",
	}
	return c.JSON(http.StatusOK, ar)

}

func (di *Dispatcher) requestVerifyHandler(c echo.Context) error {
	statistics.TrackVerifyRequest()
	dReq := mock.Request{}
	if err := c.Bind(&dReq); err != nil {
		return err
	}
	result := di.MatchSpy.Find(dReq)
	return c.JSON(http.StatusOK, result)
}

func (di *Dispatcher) resetMatchHandler(c echo.Context) error {
	statistics.TrackVerifyRequest()
	dReq := mock.Request{}
	if err := c.Bind(&dReq); err != nil {
		return err
	}

	ar := &ActionResponse{
		Result: "reset match",
	}

	di.MatchSpy.ResetMatch(dReq)
	return c.JSON(http.StatusOK, ar)
}

func (di *Dispatcher) requestResetHandler(c echo.Context) error {
	di.MatchSpy.Reset()
	ar := &ActionResponse{
		Result: "reset",
	}
	return c.JSON(http.StatusOK, ar)
}

func (di *Dispatcher) scenariosResetHandler(c echo.Context) error {
	di.Scenario.ResetAll()
	ar := &ActionResponse{
		Result: "reset",
	}
	return c.JSON(http.StatusOK, ar)
}

func (di *Dispatcher) scenariosListHandler(c echo.Context) error {
	scenarios := di.Scenario.List()
	return c.String(http.StatusOK, scenarios)
}

func (di *Dispatcher) scenariosSetHandler(c echo.Context) error {
	di.Scenario.SetState(c.Param("scenario"), c.Param("state"))
	ar := &ActionResponse{
		Result: "updated",
	}
	return c.JSON(http.StatusOK, ar)
}

func (di *Dispatcher) scenariosPauseHandler(c echo.Context) error {
	di.Scenario.SetPaused(true)
	ar := &ActionResponse{
		Result: "updated",
	}
	return c.JSON(http.StatusOK, ar)
}

func (di *Dispatcher) scenariosUnpauseHandler(c echo.Context) error {
	di.Scenario.SetPaused(false)
	ar := &ActionResponse{
		Result: "updated",
	}
	return c.JSON(http.StatusOK, ar)
}

func (di *Dispatcher) requestAllHandler(c echo.Context) error {

	return c.JSON(http.StatusOK, di.MatchSpy.GetAll())
}

func (di *Dispatcher) requestAllPagedHandler(c echo.Context) error {

	page, err := di.pageParamToInt(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &ActionResponse{
			Result: err.Error(),
		})
	}

	offset := (page - 1) * di.ResultsPerPage

	return c.JSON(http.StatusOK, di.MatchSpy.Get(di.ResultsPerPage, offset))
}

func (di *Dispatcher) requestMatchedHandler(c echo.Context) error {

	return c.JSON(http.StatusOK, di.MatchSpy.GetMatched())
}

func (di *Dispatcher) requestUnMatchedHandler(c echo.Context) error {

	return c.JSON(http.StatusOK, di.MatchSpy.GetUnMatched())
}

func (di *Dispatcher) pageParamToInt(c echo.Context) (int, error) {
	pageParam := c.Param("page")
	if !pagePattern.MatchString(pageParam) {
		return 0, ErrInvalidPage
	}

	page, err := strconv.Atoi(pageParam)
	if err != nil {
		log.Errorf("%v %v", ErrInvalidPage, err)
		return 0, ErrInvalidPage
	}

	return page, nil
}
