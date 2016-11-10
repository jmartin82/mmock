package server

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/jmartin82/mmock/definition"
	"github.com/jmartin82/mmock/parse"
	"github.com/jmartin82/mmock/persist"
	"github.com/jmartin82/mmock/proxy"
	"github.com/jmartin82/mmock/route"
	"github.com/jmartin82/mmock/translate"
)

//Dispatcher is the mock http server
type Dispatcher struct {
	IP             string
	Port           int
	Router         route.Router
	Translator     translate.MessageTranslator
	ResponseParser parse.ResponseParser
	BodyPersister  persist.BodyPersister
	Mlog           chan definition.Match
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
	var response definition.Response
	mRequest := di.Translator.BuildRequestDefinitionFromHTTP(req)

	if mRequest.Path == "/favicon.ico" {
		return
	}

	log.Printf("New request: %s %s\n", req.Method, req.URL.String())
	result := definition.Result{}
	mock, errs := di.Router.Route(&mRequest)
	if errs == nil {
		result.Found = true
	} else {
		result.Found = false
		result.Errors = errs
	}

	log.Printf("Mock match found: %s\n", strconv.FormatBool(result.Found))

	if result.Found {
		if len(mock.Control.ProxyBaseURL) > 0 {
			pr := proxy.Proxy{URL: mock.Control.ProxyBaseURL}
			response = pr.MakeRequest(mock.Request)
		} else {
			di.ResponseParser.Parse(&mRequest, &mock.Response)
			di.BodyPersister.Persist(&mock.Persist, &mRequest, &mock.Response)
			if mock.Control.Crazy {
				log.Printf("Running crazy mode")
				mock.Response.StatusCode = di.randomStatusCode(mock.Response.StatusCode)
			}
			if mock.Control.Delay > 0 {
				log.Printf("Adding a delay")
				time.Sleep(time.Duration(mock.Control.Delay) * time.Second)
			}
			response = mock.Response
		}

	}

	//translate request
	di.Translator.WriteHTTPResponseFromDefinition(&response, w)

	//log to console
	m := definition.Match{Request: mRequest, Response: response, Result: result}
	go di.recordMatchData(m)
}

//Start initialize the HTTP mock server
func (di Dispatcher) Start() {
	addr := fmt.Sprintf("%s:%d", di.IP, di.Port)

	err := http.ListenAndServe(addr, &di)
	if err != nil {
		log.Fatalf("ListenAndServe: " + err.Error())
	}
}
