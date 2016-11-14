package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"path/filepath"
	"strings"

	"github.com/jmartin82/mmock/amqp"
	"github.com/jmartin82/mmock/console"
	"github.com/jmartin82/mmock/definition"
	"github.com/jmartin82/mmock/match"
	"github.com/jmartin82/mmock/parse"
	"github.com/jmartin82/mmock/parse/fakedata"
	"github.com/jmartin82/mmock/persist"
	"github.com/jmartin82/mmock/route"
	"github.com/jmartin82/mmock/server"
	"github.com/jmartin82/mmock/translate"
)

//ErrNotFoundDefaultPath if we can't resolve the current path
var ErrNotFoundDefaultPath = errors.New("We can't determinate the current path")

//ErrNotFoundAnyMock when we don't found any valid mock definition to load
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
	router := route.NewRouter(mocks, match.MockMatch{}, dUpdates)
	router.MockChangeWatch()
	return router
}

func startServer(ip string, port int, done chan bool, router route.Router, mLog chan definition.Match, persistPath string) {
	filler := parse.FakeDataParse{Fake: fakedata.FakeAdapter{}}

	var persister persist.BodyPersister

	if strings.Index(persistPath, "mongodb://") == 0 {
		persister = persist.NewMongoBodyPersister(persistPath, filler)
	} else {
		persister = persist.NewFileBodyPersister(persistPath, filler)
	}

	sender := amqp.NewRabbitMQSender(filler)
	dispatcher := server.Dispatcher{IP: ip,
		Port:           port,
		Router:         router,
		Translator:     translate.HTTPTranslator{},
		ResponseParser: filler,
		BodyPersister:  persister,
		Mlog:           mLog,
		AMQPSender:     sender,
	}
	dispatcher.Start()
	done <- true
}
func startConsole(ip string, port int, done chan bool, mLog chan definition.Match) {
	dispatcher := console.Dispatcher{IP: ip, Port: port, Mlog: mLog}
	dispatcher.Start()
	done <- true
}

func getMocks(path string, updateCh chan []definition.Mock) []definition.Mock {
	log.Printf("Reading Mock definition from: %s\n", path)

	definitionReader := definition.NewFileDefinition(path, updateCh)

	definitionReader.AddConfigReader(definition.JSONReader{})
	definitionReader.AddConfigReader(definition.YAMLReader{})

	mocks := definitionReader.ReadMocksDefinition()
	if len(mocks) == 0 {
		log.Fatalln(ErrNotFoundAnyMock.Error())
	}
	definitionReader.WatchDir()
	return mocks
}

func main() {
	banner()
	outIP := getOutboundIP()
	path, err := filepath.Abs("./config")
	if err != nil {
		panic(ErrNotFoundDefaultPath)
	}

	persistPath, _ := filepath.Abs("./data")
	//persistPath := "mongodb://localhost/mmock"

	sIP := flag.String("server-ip", outIP, "Mock server IP")
	sPort := flag.Int("server-port", 8083, "Mock Server Port")
	cIP := flag.String("console-ip", outIP, "Console Server IP")
	cPort := flag.Int("cconsole-port", 8082, "Console server Port")
	cPath := flag.String("config-path", path, "Mocks definition folder")
	cPersistPath := flag.String("config-persist-path", persistPath, "Path to the folder where requests can be persisted or connection string to mongo database starting with mongodb:// and having database at the end /DatabaseName")
	console := flag.Bool("console", true, "Console enabled  (true/false)")
	flag.Parse()

	//chanels
	mLog := make(chan definition.Match)
	dUpdates := make(chan []definition.Mock)
	done := make(chan bool)

	path, _ = filepath.Abs(*cPath)

	if strings.Index(persistPath, "mongodb://") != 0 {
		persistPath, _ = filepath.Abs(*cPersistPath)
	}

	mocks := getMocks(path, dUpdates)
	router := getRouter(mocks, dUpdates)

	go startServer(*sIP, *sPort, done, router, mLog, persistPath)
	log.Printf("HTTP Server running at %s:%d\n", *sIP, *sPort)
	if *console {
		go startConsole(*cIP, *cPort, done, mLog)
		log.Printf("Console running at %s:%d\n", *cIP, *cPort)
	}

	<-done

}
