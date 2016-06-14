package console

import (
	"fmt"
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/jmartin82/mmock/definition"
	"golang.org/x/net/websocket"
	"html/template"
	"log"
	"net/http"
)

type Dispatcher struct {
	Ip      string
	Port    int
	Mlog    chan definition.Match
	clients []*websocket.Conn
}

func (this *Dispatcher) consoleHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := Asset("tmpl/index.html")
	t, _ := template.New("Console").Parse(string(tmpl))
	t.Execute(w, &this)
}

func (this *Dispatcher) removeClient(i int) {
	copy(this.clients[i:], this.clients[i+1:])
	this.clients[len(this.clients)-1] = nil
	this.clients = this.clients[:len(this.clients)-1]
}

func (this *Dispatcher) addClient(ws *websocket.Conn) {
	this.clients = append(this.clients, ws)
}

func (this *Dispatcher) echoHandler(ws *websocket.Conn) {
	defer func() {
		ws.Close()
	}()

	this.addClient(ws)

	//block
	var message string
	websocket.Message.Receive(ws, &message)
}

func (this *Dispatcher) logFanOut() {
	for match := range this.Mlog {
		for i, c := range this.clients {
			if c != nil {
				if err := websocket.JSON.Send(c, match); err != nil {
					this.removeClient(i)
				}
			}
		}
	}
}

func (this *Dispatcher) Start() {
	this.clients = []*websocket.Conn{}
	http.Handle("/echo", websocket.Handler(this.echoHandler))
	http.Handle("/js/", http.FileServer(&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo, Prefix: "tmpl"}))
	http.Handle("/css/", http.FileServer(&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo, Prefix: "tmpl"}))
	http.HandleFunc("/", this.consoleHandler)

	go this.logFanOut()

	addr := fmt.Sprintf("%s:%d", this.Ip, this.Port)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("ListenAndServe: " + err.Error())
	}
}
