package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/jmartin82/mmock/console"
	"github.com/jmartin82/mmock/definition"
	"github.com/jmartin82/mmock/definition/parser"
	"github.com/jmartin82/mmock/match"
	"github.com/jmartin82/mmock/scenario"
	"github.com/jmartin82/mmock/server"
	"github.com/jmartin82/mmock/statistics"
	"github.com/jmartin82/mmock/translate"
	"github.com/jmartin82/mmock/vars"
	"github.com/jmartin82/mmock/vars/fakedata"
)

//ErrNotFoundPath error from missing or configuration path
var ErrNotFoundPath = errors.New("Configuration path not found")

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

func existsConfigPath(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func getMapping(path string) definition.Mapping {
	path, _ = filepath.Abs(path)
	if !existsConfigPath(path) {
		log.Fatalf(ErrNotFoundPath.Error())
	}

	configMapper := definition.NewConfigMapper()

	configMapper.AddConfigParser(parser.JSONReader{})
	configMapper.AddConfigParser(parser.YAMLReader{})

	mapping := definition.NewConfigMapping(path, configMapper)

	filepath.Walk(path, func(filePath string, fileInfo os.FileInfo, err error) error {
		if !fileInfo.IsDir() {
			URI := strings.TrimPrefix(filePath, path)
			mapping.Load(URI)
		}
		return nil
	})

	return mapping
}

func getRouter(mapping definition.Mapping, checker match.Checker) *server.Router {
	router := server.NewRouter(mapping, checker)
	return router
}

func getVarsProcessor() vars.Processor {

	return vars.Processor{FillerFactory: vars.MockFillerFactory{FakeAdapter: fakedata.FakeAdapter{}}}
}

func startServer(ip string, port int, done chan bool, router server.Resolver, mLog chan definition.Match, scenario scenario.Director, varsProcessor vars.Processor, spier match.Spier) {
	dispatcher := server.Dispatcher{
		IP:         ip,
		Port:       port,
		Resolver:   router,
		Translator: translate.HTTP{},
		Processor:  varsProcessor,
		Scenario:   scenario,
		Spier:      spier,
		Mlog:       mLog,
	}
	dispatcher.Start()
	done <- true
}
func startConsole(ip string, port int, spy match.Spier, scenario scenario.Director, done chan bool, mLog chan definition.Match) {
	dispatcher := console.Dispatcher{
		IP:       ip,
		Port:     port,
		MatchSpy: spy,
		Scenario: scenario,
		Mlog:     mLog}
	dispatcher.Start()
	done <- true
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
	sStatistics := flag.Bool("server-statistics", true, "Mock server sends anonymous statistics")
	cIP := flag.String("console-ip", outIP, "Console server IP")
	cPort := flag.Int("console-port", 8082, "Console server Port")
	console := flag.Bool("console", true, "Console enabled  (true/false)")
	cPath := flag.String("config-path", path, "Mocks definition folder")

	flag.Parse()

	//chanels
	mLog := make(chan definition.Match)
	done := make(chan bool)

	//shared structs
	scenario := scenario.NewMemoryStore()
	checker := match.NewTester(scenario)
	matchStore := match.NewMemoryStore()

	mapping := getMapping(*cPath)
	spy := getMatchSpy(checker, matchStore)
	router := getRouter(mapping, checker)
	varsProcessor := getVarsProcessor()

	if !(*sStatistics) {
		statistics.SetMonitor(statistics.NewNullableMonitor())
		log.Printf("Not sending statistics\n")
	}
	defer statistics.Stop()

	go startServer(*sIP, *sPort, done, router, mLog, scenario, varsProcessor, spy)
	log.Printf("HTTP Server running at %s:%d\n", *sIP, *sPort)
	if *console {
		go startConsole(*cIP, *cPort, spy, scenario, done, mLog)
		log.Printf("Console running at %s:%d\n", *cIP, *cPort)
	}

	<-done

}
