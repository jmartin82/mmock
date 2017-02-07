package console

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/jmartin82/mmock/definition"
	"github.com/jmartin82/mmock/match"
	"github.com/jmartin82/mmock/translate"
	"golang.org/x/net/websocket"
)

//Dispatcher is the http console server.
type Dispatcher struct {
	IP         string
	Port       int
	Translator translate.MessageTranslator
	Verifier   match.Verifier
	Mlog       chan definition.Match
	clients    []*websocket.Conn
}

func (di *Dispatcher) consoleHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := Asset("tmpl/index.html")
	fmt.Fprintf(w, string(tmpl))
}

func (di *Dispatcher) verifyHandler(w http.ResponseWriter, r *http.Request) {

	dReq := definition.Request{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&dReq)
	if err != nil {

	}
	defer r.Body.Close()
	result := di.Verifier.Verify(dReq)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)

}

func (di *Dispatcher) removeClient(i int) {
	copy(di.clients[i:], di.clients[i+1:])
	di.clients[len(di.clients)-1] = nil
	di.clients = di.clients[:len(di.clients)-1]
}

func (di *Dispatcher) addClient(ws *websocket.Conn) {
	di.clients = append(di.clients, ws)
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
	http.HandleFunc("/verify", di.verifyHandler)
	http.HandleFunc("/", di.consoleHandler)

	go di.logFanOut()

	addr := fmt.Sprintf("%s:%d", di.IP, di.Port)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("ListenAndServe: " + err.Error())
	}
}
