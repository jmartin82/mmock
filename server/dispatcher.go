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
	"github.com/jmartin82/mmock/route"
	"github.com/jmartin82/mmock/translate"
)

//Dispatcher is the mock http server
type Dispatcher struct {
	Ip             string
	Port           int
	Router         route.Router
	Translator     translate.MessageTranslator
	ResponseParser parse.ResponseParser
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

	di.ResponseParser.Parse(&mRequest, &mock.Response)

	if result.Found {
		if mock.Control.Crazy {
			log.Printf("Crazy mode enabled")
			mock.Response.StatusCode = di.randomStatusCode(mock.Response.StatusCode)
		}
		if mock.Control.Delay > 0 {
			time.Sleep(time.Duration(mock.Control.Delay) * time.Second)
		}
	}

	//translate request
	di.Translator.WriteHTTPResponseFromDefinition(&mock.Response, w)

	//log to console
	m := definition.Match{Request: mRequest, Response: mock.Response, Result: result}
	go di.recordMatchData(m)
}

//Start initialize the HTTP mock server
func (di Dispatcher) Start() {
	addr := fmt.Sprintf("%s:%d", di.Ip, di.Port)

	err := http.ListenAndServe(addr, &di)
	if err != nil {
		log.Fatalf("ListenAndServe: " + err.Error())
	}
}
