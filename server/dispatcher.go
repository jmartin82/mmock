package server

import (
	"fmt"
	"github.com/jmartin82/mmock/definition"
	"github.com/jmartin82/mmock/parse"
	"github.com/jmartin82/mmock/route"
	"github.com/jmartin82/mmock/translate"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type Dispatcher struct {
	Ip             string
	Port           int
	Router         route.Router
	Translator     translate.MessageTranslator
	ResponseParser parse.ResponseParser
	Mlog           chan definition.Match
}

func (this Dispatcher) recordMatchData(msg definition.Match) {
	this.Mlog <- msg
}

func (this Dispatcher) randomStatusCode(currentStatus int) int {
	if time.Now().Second()%2 == 0 {
		rand.Seed(time.Now().Unix())
		return rand.Intn(4) + 500
	}
	return currentStatus
}

func (this *Dispatcher) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	mRequest := this.Translator.BuildRequestDefinitionFromHTTP(req)

	if mRequest.Path == "/favicon.ico" {
		return
	}

	log.Printf("New request: %s %s\n", req.Method, req.URL.String())
	result := definition.Result{}
	mock, errs := this.Router.Route(&mRequest)
	if errs == nil {
		result.Found = true
	} else {
		result.Found = false
		result.Errors = errs
	}

	log.Printf("Mock match found: %s\n", strconv.FormatBool(result.Found))

	this.ResponseParser.Parse(&mRequest, &mock.Response)

	if result.Found {
		if mock.Control.Crazy {
			log.Printf("Crazy mode enabled")
			mock.Response.StatusCode = this.randomStatusCode(mock.Response.StatusCode)
		}
		if mock.Control.Delay > 0 {
			time.Sleep(time.Duration(mock.Control.Delay) * time.Second)
		}
	}

	//translate request
	this.Translator.WriteHTTPResponseFromDefinition(&mock.Response, w)

	//log to console
	m := definition.Match{mRequest, mock.Response, result}
	go this.recordMatchData(m)
}

func (this Dispatcher) Start() {
	addr := fmt.Sprintf("%s:%d", this.Ip, this.Port)

	err := http.ListenAndServe(addr, &this)
	if err != nil {
		log.Fatalf("ListenAndServe: " + err.Error())
	}
}
