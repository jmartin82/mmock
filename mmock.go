package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"path/filepath"
	"strings"

	"github.com/jmartin82/mmock/console"
	"github.com/jmartin82/mmock/definition"
	"github.com/jmartin82/mmock/match"
	"github.com/jmartin82/mmock/route"
	"github.com/jmartin82/mmock/scenario"
	"github.com/jmartin82/mmock/server"
	"github.com/jmartin82/mmock/translate"
	"github.com/jmartin82/mmock/vars"
	"github.com/jmartin82/mmock/vars/fakedata"
)

//ErrNotFoundDefaultPath if we can't resolve the current path
var ErrNotFoundDefaultPath = errors.New("We can't determinate the current path")

//ErrNotFoundAnyMock when we don't found any valid mock definition to load
var ErrNotFoundAnyMock = errors.New("No valid mock definition found")

func banner() {
	fmt.Println("MMock v 2.0.2")
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

func getMatchSpy(checker match.Checker, matchStore match.Store) match.Spier {
	return match.NewSpy(checker, matchStore)
}

func getRouter(mocks []definition.Mock, checker match.Checker, dUpdates chan []definition.Mock) *route.RequestRouter {
	log.Printf("Loding router with %d definitions\n", len(mocks))

	router := route.NewRouter(mocks, checker, dUpdates)
	router.MockChangeWatch()
	return router
}

func getVarsProcessor() vars.VarsProcessor {

	return vars.VarsProcessor{FillerFactory: vars.MockFillerFactory{FakeAdapter: fakedata.FakeAdapter{}}}
}

func startServer(ip string, port int, done chan bool, router route.Router, mLog chan definition.Match, scenario scenario.ScenarioManager, varsProcessor vars.VarsProcessor, matchStore match.Store) {
	dispatcher := server.Dispatcher{
		IP:            ip,
		Port:          port,
		Router:        router,
		Translator:    translate.HTTPTranslator{},
		VarsProcessor: varsProcessor,

		Scenario: scenario,
		Store:    matchStore,
		Mlog:     mLog,
	}
	dispatcher.Start()
	done <- true
}
func startConsole(ip string, port int, spy match.Spier, done chan bool, mLog chan definition.Match) {
	dispatcher := console.Dispatcher{
		IP:       ip,
		Port:     port,
		MatchSpy: spy,
		Mlog:     mLog}
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

	sIP := flag.String("server-ip", outIP, "Mock server IP")
	sPort := flag.Int("server-port", 8083, "Mock server Port")
	cIP := flag.String("console-ip", outIP, "Console server IP")
	cPort := flag.Int("console-port", 8082, "Console server Port")
	console := flag.Bool("console", true, "Console enabled  (true/false)")
	cPath := flag.String("config-path", path, "Mocks definition folder")

	flag.Parse()
	path, _ = filepath.Abs(*cPath)

	//chanels
	mLog := make(chan definition.Match)
	dUpdates := make(chan []definition.Mock)
	done := make(chan bool)

	//shared structs
	scenario := scenario.NewInMemoryScenario()
	checker := match.NewTester(scenario)
	matchStore := match.NewMemoryStore()

	mocks := getMocks(path, dUpdates)
	spy := getMatchSpy(checker, matchStore)
	router := getRouter(mocks, checker, dUpdates)
	varsProcessor := getVarsProcessor()

	go startServer(*sIP, *sPort, done, router, mLog, scenario, varsProcessor, matchStore)
	log.Printf("HTTP Server running at %s:%d\n", *sIP, *sPort)
	if *console {
		go startConsole(*cIP, *cPort, spy, done, mLog)
		log.Printf("Console running at %s:%d\n", *cIP, *cPort)
	}

	<-done

}
