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

	"github.com/jmartin82/mmock/internal/config"
	"github.com/jmartin82/mmock/internal/config/parser"
	"github.com/jmartin82/mmock/internal/console"
	"github.com/jmartin82/mmock/internal/server"
	"github.com/jmartin82/mmock/internal/statistics"
	"github.com/jmartin82/mmock/pkg/match"
	"github.com/jmartin82/mmock/pkg/match/payload"
	"github.com/jmartin82/mmock/pkg/mock"
	"github.com/jmartin82/mmock/pkg/vars"
	"github.com/jmartin82/mmock/pkg/vars/fake"
)

//VERSION of the application
var VERSION string = "development"

//ErrNotFoundPath error from missing or configuration path
var ErrNotFoundPath = errors.New("Configuration path not found")

//ErrNotFoundDefaultPath if we can't resolve the current path
var ErrNotFoundDefaultPath = errors.New("We can't determinate the current path")

//ErrNotFoundAnyMock when we don't found any valid mock config to load
var ErrNotFoundAnyMock = errors.New("No valid mock config found")

func banner() {
	fmt.Printf("MMock v %s", VERSION)
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

func getTransactionSpy(checker match.Matcher, matchStore match.TransactionStorer) *match.Spy {
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

func getMapping(path string) config.Mapping {
	path, _ = filepath.Abs(path)
	if !existsConfigPath(path) {
		log.Fatalln(ErrNotFoundPath.Error())
	}

	fsMapper := config.NewFileSystemMapper()
	fsMapper.AddParser(parser.JSONReader{})
	fsMapper.AddParser(parser.YAMLReader{})

	fsUpdate := make(chan struct{})

	watcher := config.NewFileWatcher(path, fsUpdate)
	watcher.Bind()

	return config.NewConfigMapping(path, fsMapper, fsUpdate)
}

func getRouter(mapping config.Mapping, checker match.Matcher) *server.Router {
	router := server.NewRouter(mapping, checker)
	return router
}

func getVarsProcessor() *vars.ResponseMessageEvaluator {
	ccg := fake.NewCreditCardGenerator()
	fp := fake.NewFakeDataProvider(ccg)
	ff := vars.NewFillerFactory(fp)
	return vars.NewResponseMessageEvaluator(ff)
}

func startServer(ip string, port, portTLS int, configTLS string, done chan struct{}, router server.RequestResolver, mLog chan match.Transaction, scenario match.ScenearioStorer, varsProcessor vars.Evaluator, spier match.TransactionSpier) {
	dispatcher := server.Dispatcher{
		IP:         ip,
		Port:       port,
		PortTLS:    portTLS,
		ConfigTLS:  configTLS,
		Resolver:   router,
		Translator: mock.HTTP{},
		Evaluator:  varsProcessor,
		Scenario:   scenario,
		Spier:      spier,
		Mlog:       mLog,
	}
	dispatcher.Start()
	done <- struct{}{}
}
func startConsole(ip string, port int, resultsPerPage int, spy match.TransactionSpier, scenario match.ScenearioStorer, mapping config.Mapping, done chan struct{}, mLog chan match.Transaction) {
	dispatcher := console.Dispatcher{
		IP:             ip,
		Port:           port,
		MatchSpy:       spy,
		Scenario:       scenario,
		Mapping:        mapping,
		Mlog:           mLog,
		ResultsPerPage: resultsPerPage,
	}
	dispatcher.Start()
	done <- struct{}{}
}

func main() {
	banner()
	outIP := getOutboundIP()
	path, err := filepath.Abs("./config")
	TLS, err := filepath.Abs("./tls")
	if err != nil {
		panic(ErrNotFoundDefaultPath)
	}

	sIP := flag.String("server-ip", outIP, "Definition server IP")
	sPort := flag.Int("server-port", 8083, "Definition server Port")
	sPortTLS := flag.Int("server-tls-port", 8084, "Definition server TLS Port")
	sStatistics := flag.Bool("server-statistics", true, "Definition server sends anonymous statistics")
	cIP := flag.String("console-ip", outIP, "Console server IP")
	cPort := flag.Int("console-port", 8082, "Console server Port")
	console := flag.Bool("console", true, "Console enabled  (true/false)")
	cPath := flag.String("config-path", path, "Mocks config folder")
	cTLS := flag.String("tls-path", TLS, "TLS config folder (server.crt and server.key should be inside)")
	cStorageCapacity := flag.Int("request-storage-capacity", 100, "Request storage capacity (0 = infinite)")
	cResultsPerPage := flag.Int("results-per-page", 25, "Number of results per page")

	flag.Parse()

	//chanels
	mLog := make(chan match.Transaction)
	done := make(chan struct{})

	//shared structs
	scenario := match.NewInMemoryScenarioStore()
	comparator := payload.NewDefaultComparator()
	tester := match.NewTester(comparator, scenario)
	matchStore := match.NewInMemoryTransactionStore(tester, *cStorageCapacity)
	mapping := getMapping(*cPath)
	spy := getTransactionSpy(tester, matchStore)
	router := getRouter(mapping, tester)
	varsProcessor := getVarsProcessor()

	if *sStatistics {
		fmt.Printf("\n************************************************************************************\n")
		fmt.Printf("* Mmock is collecting anonymous statistics about the usage of the features.        *\n")
		fmt.Printf("* You can disable this behavior adding the following flag -server-statistics=false *\n")
		fmt.Printf("************************************************************************************\n\n")
	} else {
		statistics.SetMonitor(statistics.NewNullableMonitor())
	}

	defer statistics.Stop()

	go startServer(*sIP, *sPort, *sPortTLS, *cTLS, done, router, mLog, scenario, varsProcessor, spy)
	log.Printf("HTTP Server running at http://%s:%d\n", *sIP, *sPort)
	log.Printf("HTTPS Server running at https://%s:%d\n", *sIP, *sPortTLS)
	if *console {
		go startConsole(*cIP, *cPort, *cResultsPerPage, spy, scenario, mapping, done, mLog)
		log.Printf("Console running at http://%s:%d\n", *cIP, *cPort)
	}

	<-done

}
