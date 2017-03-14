package console

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/jmartin82/mmock/definition"
	"github.com/jmartin82/mmock/match"
	"github.com/jmartin82/mmock/statistics"
	"golang.org/x/net/websocket"
)

//Dispatcher is the http console server.
type Dispatcher struct {
	IP       string
	Port     int
	MatchSpy match.Spier
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
	di.clients = []*websocket.Conn{}
	http.Handle("/echo", websocket.Handler(di.echoHandler))
	http.Handle("/js/", http.FileServer(&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo, Prefix: "tmpl"}))
	http.Handle("/css/", http.FileServer(&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo, Prefix: "tmpl"}))

	//verification
	http.HandleFunc("/request/reset", di.requestReset)
	http.HandleFunc("/request/verify", di.requestVerifyHandler)
	http.HandleFunc("/request/all", di.requestAllHandler)
	http.HandleFunc("/request/matched", di.requestMatchedHandler)
	http.HandleFunc("/request/unmatched", di.requestUnMatchedHandler)

	http.HandleFunc("/", di.consoleHandler)

	go di.logFanOut()

	addr := fmt.Sprintf("%s:%d", di.IP, di.Port)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("ListenAndServe: " + err.Error())
	}
}

//CONSOLE
func (di *Dispatcher) consoleHandler(w http.ResponseWriter, r *http.Request) {
	statistics.TrackConsoleRequest()
	tmpl, _ := Asset("tmpl/index.html")
	fmt.Fprintf(w, string(tmpl))
}

func (di *Dispatcher) echoHandler(ws *websocket.Conn) {
	defer func() {
		ws.Close()
	}()

	di.addClient(ws)

	//block
	var message string
	websocket.Message.Receive(ws, &message)
}

//API REQUEST

func (di *Dispatcher) requestVerifyHandler(w http.ResponseWriter, r *http.Request) {
	statistics.TrackVerifyRequest()
	dReq := definition.Request{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&dReq)
	if err != nil {
		log.Println("Invalid request input")
	}
	defer r.Body.Close()
	result := di.MatchSpy.Find(dReq)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)

}

func (di *Dispatcher) requestReset(w http.ResponseWriter, r *http.Request) {
	di.MatchSpy.Reset()
}

func (di *Dispatcher) requestAllHandler(w http.ResponseWriter, r *http.Request) {
	result := di.MatchSpy.GetAll()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)

}

func (di *Dispatcher) requestMatchedHandler(w http.ResponseWriter, r *http.Request) {
	result := di.MatchSpy.GetMatched()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (di *Dispatcher) requestUnMatchedHandler(w http.ResponseWriter, r *http.Request) {
	result := di.MatchSpy.GetUnMatched()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
