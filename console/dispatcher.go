package console

import (
	"fmt"
	"net/http"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/jmartin82/mmock/definition"
	"github.com/jmartin82/mmock/match"
	"github.com/jmartin82/mmock/scenario"
	"github.com/jmartin82/mmock/statistics"
	"github.com/labstack/echo"
	"golang.org/x/net/websocket"
)

type ActionResponse struct {
	Result string `json:"result"`
}

//Dispatcher is the http console server.
type Dispatcher struct {
	IP       string
	Port     int
	MatchSpy match.Spier
	Scenario scenario.Director
	Mapping  definition.Mapping
	Mlog     chan definition.Match
	clients  []*websocket.Conn
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

//Start initiates the http console.
func (di *Dispatcher) Start() {
	e := echo.New()
	//WS
	di.clients = []*websocket.Conn{}
	e.GET("/echo", di.webSocketHandler)

	//HTTP
	statics := http.FileServer(&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo, Prefix: "tmpl"})
	e.GET("/js/*", echo.WrapHandler(statics))
	e.GET("/css/*", echo.WrapHandler(statics))
	e.GET("/", di.consoleHandler)

	//verification
	e.GET("/__admin/request/reset", di.requestResetHandler)
	e.GET("/__admin/request/verify", di.requestVerifyHandler)
	e.GET("/__admin/request/all", di.requestAllHandler)
	e.GET("/__admin/request/matched", di.requestMatchedHandler)
	e.GET("/__admin/request/unmatched", di.requestUnMatchedHandler)
	e.GET("/__admin/scenarios/reset_all", di.scenariosResetHandler)

	//mapping
	e.GET("/__admin/mapping", di.mappingListHandler)
	e.GET("/__admin/mapping/:uri", di.mappingGetHandler)
	e.POST("/__admin/mapping/:uri", di.mappingCreateHandler)
	e.PUT("/__admin/mapping/:uri", di.mappingUpdateHandler)
	e.DELETE("/__admin/mapping/:uri", di.mappingDeleteHandler)

	//POST __admin/mapping (all)

	go di.logFanOut()

	addr := fmt.Sprintf("%s:%d", di.IP, di.Port)
	e.Logger.Fatal(e.Start(addr))
}

//CONSOLE
func (di *Dispatcher) consoleHandler(c echo.Context) error {
	statistics.TrackConsoleRequest()
	tmpl, _ := Asset("tmpl/index.html")
	return c.HTML(http.StatusOK, string(tmpl))
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

//API REQUEST
func (di *Dispatcher) mappingListHandler(c echo.Context) (err error) {
	mocks := di.Mapping.List()
	return c.JSON(http.StatusOK, mocks)
}

func (di *Dispatcher) mappingGetHandler(c echo.Context) (err error) {

	URI := c.Param("uri")
	mock := definition.Mock{}
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

	URI := c.Param("uri")
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

	mock := &definition.Mock{}
	URI := c.Param("uri")

	if _, ok := di.Mapping.Get(URI); ok {
		ar := &ActionResponse{
			Result: "already_exists",
		}
		return c.JSON(http.StatusConflict, ar)
	}

	if err = c.Bind(mock); err != nil {
		return
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

	mock := &definition.Mock{}
	URI := c.Param("uri")

	if _, ok := di.Mapping.Get(URI); !ok {
		ar := &ActionResponse{
			Result: "not_found",
		}
		return c.JSON(http.StatusNotFound, ar)
	}

	if err = c.Bind(mock); err != nil {
		return
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
	dReq := definition.Request{}
	if err := c.Bind(&dReq); err != nil {
		return err
	}
	result := di.MatchSpy.Find(dReq)
	return c.JSON(http.StatusOK, result)
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

func (di *Dispatcher) requestAllHandler(c echo.Context) error {

	return c.JSON(http.StatusOK, di.MatchSpy.GetAll())
}

func (di *Dispatcher) requestMatchedHandler(c echo.Context) error {

	return c.JSON(http.StatusOK, di.MatchSpy.GetMatched())
}

func (di *Dispatcher) requestUnMatchedHandler(c echo.Context) error {

	return c.JSON(http.StatusOK, di.MatchSpy.GetUnMatched())
}
