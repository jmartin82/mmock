package console

import (
	"fmt"
	"log"
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/jmartin82/mmock/definition"
	"golang.org/x/net/websocket"
)

//Dispatcher is the http console server.
type Dispatcher struct {
	IP      string
	Port    int
	Mlog    chan definition.Match
	clients []*websocket.Conn
}

func (di *Dispatcher) consoleHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := Asset("tmpl/index.html")
	//t, _ := template.New("Console").Parse(string(tmpl))
	//t.Execute(w, &di)
	fmt.Fprintf(w, string(tmpl))
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
	http.HandleFunc("/", di.consoleHandler)

	go di.logFanOut()

	addr := fmt.Sprintf("%s:%d", di.IP, di.Port)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("ListenAndServe: " + err.Error())
	}
}
