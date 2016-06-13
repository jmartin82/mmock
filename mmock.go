package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"mmock/console"
	"mmock/definition"
	"mmock/match"
	"mmock/parse"
	"mmock/parse/fakedata"
	"mmock/route"
	"mmock/server"
	"mmock/translate"
	"net"
	"path/filepath"
	"strings"
)

var ErrNotFoundDefaultPath = errors.New("We can't determinate the current path")
var ErrNotFoundAnyMock = errors.New("No valid mock definition found")

func banner() {
	fmt.Println("MMock v 0.0.1")
	fmt.Println("")

	fmt.Print(
		`		.---. .---. 
               :     : o   :    me want request!
           _..-:   o :     :-.._    /
       .-''  '  ` + "`" + `---' ` + "`" + `---' "   ` + "`" + `` + "`" + `-.    
     .'   "   '  "  .    "  . '  "  ` + "`" + `.  
    :   '.---.,,.,...,.,.,.,..---.  ' ;
    ` + "`" + `. " ` + "`" + `.                     .' " .'
     ` + "`" + `.  '` + "`" + `.                   .' ' .'
      ` + "`" + `.    ` + "`" + `-._           _.-' "  .'  .----.
        ` + "`" + `. "    '"--...--"'  . ' .'  .'  o   ` + "`" + `.
        .'` + "`" + `-._'    " .     " _.-'` + "`" + `. :       o  :
      .'      ` + "`" + `` + "`" + `` + "`" + `--.....--'''    ' ` + "`" + `:_ o       :
    .'    "     '         "     "   ; ` + "`" + `.;";";";'
   ;         '       "       '     . ; .' ; ; ;
  ;     '         '       '   "    .'      .-'
  '  "     "   '      "           "    _.-' 
 `)
	fmt.Println("")
}

// Get preferred outbound ip of this machine
func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "127.0.0.1"
	}
	defer conn.Close()

	log.Println("Getting external IP")
	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")

	return localAddr[0:idx]
}

func getRouter(mocks []definition.Mock, dUpdates chan []definition.Mock) *route.RequestRouter {
	log.Printf("Loding router with %d definitions\n", len(mocks))
	return &route.RequestRouter{Mocks: mocks, Matcher: match.MockMatch{}, DUpdates: dUpdates}
}

func startServer(ip string, port int, done chan bool, router route.Router, mLog chan definition.Match) {
	filler := parse.FakeDataParse{fakedata.FakeAdapter{}}
	dispatcher := server.Dispatcher{ip, port, router, translate.HTTPTranslator{}, filler, mLog}
	dispatcher.Start()
	done <- true
}
func startConsole(ip string, port int, done chan bool, mLog chan definition.Match) {
	dispatcher := console.Dispatcher{Ip: ip, Port: port, Mlog: mLog}
	dispatcher.Start()
	done <- true
}

func main() {
	banner()
	outIp := getOutboundIP()
	path, err := filepath.Abs("./config")
	if err != nil {
		panic(ErrNotFoundDefaultPath)
	}

	sIp := flag.String("server-ip", outIp, "Mock server IP")
	sPort := flag.Int("server-port", 8082, "Mock Server Port")
	cIp := flag.String("console-ip", outIp, "Console Server IP")
	cPort := flag.Int("cconsole-port", 8083, "Console server Port")
	cPath := flag.String("config-path", path, "Mocks definition folder")
	console := flag.Bool("console", true, "Console enabled  (true/false)")
	flag.Parse()

	//chanels
	mLog := make(chan definition.Match)
	dUpdates := make(chan []definition.Mock)
	done := make(chan bool)

	path, _ = filepath.Abs(*cPath)
	log.Printf("Reading Mock definition from: %s\n", path)
	definitionReader := definition.FileDefinition{path, dUpdates}
	mocks := definitionReader.ReadMocksDefinition()
	if len(mocks) == 0 {
		log.Fatalln(ErrNotFoundAnyMock.Error())
	}
	definitionReader.WatchDir()

	router := getRouter(mocks, dUpdates)
	router.MockChangeWatch()

	go startServer(*cIp, *cPort, done, router, mLog)
	log.Printf("HTTP Server running at %s:%d\n", *cIp, *cPort)
	if *console {
		go startConsole(*sIp, *sPort, done, mLog)
		log.Printf("Console running at %s:%d\n", *sIp, *sPort)
	}

	<-done

}
