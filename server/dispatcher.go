package server

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/jmartin82/mmock/definition"
	"github.com/jmartin82/mmock/match"
	"github.com/jmartin82/mmock/proxy"
	"github.com/jmartin82/mmock/scenario"
	"github.com/jmartin82/mmock/statistics"
	"github.com/jmartin82/mmock/translate"
	"github.com/jmartin82/mmock/vars"
)

//Dispatcher is the mock http server
type Dispatcher struct {
	IP         string
	Port       int
	Resolver   Resolver
	Translator translate.MessageTranslator
	Processor  vars.Processor
	Scenario   scenario.Director
	Spier      match.Spier
	Mlog       chan definition.Match
}

func (di Dispatcher) recordMatchData(msg definition.Match) {
	di.Mlog <- msg
}

func (di Dispatcher) randomStatusCode(currentStatus int) int {
	if time.Now().Second()%2 == 0 {
		rand.Seed(time.Now().Unix())
		return rand.Intn(4) + 500
	}
	return currentStatus
}

//ServerHTTP is the mock http server request handler.
//It uses the router to decide the matching mock and translator as adapter between the HTTP impelementation and the mock definition.
func (di *Dispatcher) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	mRequest := di.Translator.BuildRequestDefinitionFromHTTP(req)

	if mRequest.Path == "/favicon.ico" {
		return
	}

	log.Printf("New request: %s %s\n", req.Method, req.URL.String())

	mock, match := di.getMatchingResult(&mRequest)

	//save the match info
	di.Spier.Save(*match)

	//set new scenario
	if mock.Control.Scenario.NewState != "" {
		di.Scenario.SetState(mock.Control.Scenario.Name, mock.Control.Scenario.NewState)
	}

	//translate request
	di.Translator.WriteHTTPResponseFromDefinition(match.Response, w)

	go di.recordMatchData(*match)
}

func (di *Dispatcher) getMatchingResult(request *definition.Request) (*definition.Mock, *definition.Match) {
	response := &definition.Response{}
	result := definition.Result{}
	mock, errs := di.Resolver.Resolve(request)
	if errs == nil {
		result.Found = true
	} else {
		result.Found = false
		result.Errors = errs
	}

	log.Printf("Mock match found: %s. Name : %s\n", strconv.FormatBool(result.Found), mock.Name)

	if result.Found {
		if len(mock.Control.ProxyBaseURL) > 0 {
			pr := proxy.Proxy{URL: mock.Control.ProxyBaseURL}
			response = pr.MakeRequest(mock.Request)
		} else {

			di.Processor.Eval(request, mock)
			if mock.Control.Crazy {
				log.Printf("Running crazy mode")
				mock.Response.StatusCode = di.randomStatusCode(mock.Response.StatusCode)
			}
			if mock.Control.Delay > 0 {
				log.Printf("Adding a delay")
				time.Sleep(time.Duration(mock.Control.Delay) * time.Second)
			}
			response = &mock.Response
		}
		statistics.TrackSuccesfulRequest()
	} else {
		response = &mock.Response
	}

	return mock, &definition.Match{Time: time.Now().Unix(), Request: request, Response: response, Result: result}
}

//Start initialize the HTTP mock server
func (di Dispatcher) Start() {
	addr := fmt.Sprintf("%s:%d", di.IP, di.Port)

	err := http.ListenAndServe(addr, &di)
	if err != nil {
		log.Fatalf("ListenAndServe: " + err.Error())
	}
}
