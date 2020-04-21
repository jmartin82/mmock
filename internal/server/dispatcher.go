package server

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/jmartin82/mmock/internal/proxy"
	"github.com/jmartin82/mmock/pkg/match"

	"net/url"
	"regexp"
	"strings"

	"github.com/jmartin82/mmock/internal/statistics"
	"github.com/jmartin82/mmock/pkg/mock"
	"github.com/jmartin82/mmock/pkg/vars"
)

//Dispatcher is the mock http server
type Dispatcher struct {
	IP         string
	Port       int
	PortTLS    int
	ConfigTLS  string
	Resolver   RequestResolver
	Translator mock.MessageTranslator
	Evaluator  vars.Evaluator
	Scenario   match.ScenearioStorer
	Spier      match.TransactionSpier
	Mlog       chan match.Transaction
}

func (di Dispatcher) recordMatchData(msg match.Transaction) {
	di.Mlog <- msg
}

func (di Dispatcher) randomStatusCode(currentStatus int) int {
	if time.Now().Second()%2 == 0 {
		rand.Seed(time.Now().Unix())
		return rand.Intn(4) + 500
	}
	return currentStatus
}

func (di Dispatcher) callWebHook(url string, match *match.Transaction) {
	log.Printf("Running webhook: %s\n", url)
	content, err := json.Marshal(match)
	if err != nil {
		log.Println("Impossible encode the WebHook payload")
		return
	}
	reader := bytes.NewReader(content)
	resp, err := http.Post(url, "application/json", reader)
	if err != nil {
		log.Printf("Impossible send payload to: %s\n", url)
		return
	}
	log.Printf("WebHook response: %d\n", resp.StatusCode)
}

func (di *Dispatcher) handleCallback(mRequest *mock.Request, mock *mock.Definition) {
	url := mock.Callback.Url
	contentType := mock.Callback.ContentType

	if d := mock.Callback.Delay.Duration; d > 0 {
		log.Printf("Delaying callback by: %s\n", d)
		time.Sleep(d)
	}

	log.Printf("Making callback to %s\n", url)
	_, err := http.Post(url, contentType, bytes.NewBufferString(mock.Callback.Body))
	if err != nil {
		log.Printf("Error from callback: %s\n", err)
	} else {
		log.Println("Callback handled successfully")
	}
}

//ServerHTTP is the mock http server request handler.
//It uses the router to decide the matching mock and translator as adapter between the HTTP impelementation and the mock mock.
func (di *Dispatcher) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	mRequest := di.Translator.BuildRequestDefinitionFromHTTP(req)

	if mRequest.Path == "/favicon.ico" {
		return
	}

	log.Printf("New request: %s %s\n", req.Method, req.URL.String())

	mock, transaction := di.getMatchingResult(&mRequest)

	//save the match info
	di.Spier.Save(*transaction)

	//set new scenario
	if mock.Control.Scenario.NewState != "" {
		statistics.TrackScenarioFeature()
		di.Scenario.SetState(mock.Control.Scenario.Name, mock.Control.Scenario.NewState)
	}

	if mock.Control.WebHookURL != "" {
		go di.callWebHook(mock.Control.WebHookURL, transaction)
	}

	//translate request
	di.Translator.WriteHTTPResponseFromDefinition(transaction.Response, w)

	if mock.Callback.Url != "" {
		go di.handleCallback(&mRequest, mock)
	}

	go di.recordMatchData(*transaction)
}

func getProxyResponse(request *mock.Request, definition *mock.Definition) *mock.Response {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	mockProxyBaseURL := definition.Control.ProxyBaseURL
	parsedMockProxyBaseURL, err := url.Parse(mockProxyBaseURL)

	if err != nil {
		return &mock.Response{StatusCode: http.StatusInternalServerError}
	}

	authRegexp := regexp.MustCompile(definition.Request.Path)
	matches := authRegexp.FindStringSubmatch(parsedMockProxyBaseURL.Path)

	if len(matches) > 0 {
		mockProxyBaseURL = strings.Replace(mockProxyBaseURL, matches[0], "", -1)
		mockProxyBaseURL += request.Path
	}

	pr := proxy.Proxy{URL: mockProxyBaseURL, Client: client}

	return pr.MakeRequest(request)
}

func (di *Dispatcher) getMatchingResult(request *mock.Request) (*mock.Definition, *match.Transaction) {
	response := &mock.Response{}
	mock, result := di.Resolver.Resolve(request)

	log.Printf("Definition match found: %s. Name : %s\n", strconv.FormatBool(result.Found), mock.URI)

	if result.Found {
		if len(mock.Control.ProxyBaseURL) > 0 {
			statistics.TrackProxyFeature()
			response = getProxyResponse(request, mock)
		} else {

			di.Evaluator.Eval(request, mock)
			if mock.Control.Crazy {
				log.Println("Running crazy mode")
				mock.Response.StatusCode = di.randomStatusCode(mock.Response.StatusCode)
			}

			if d := mock.Control.Delay.Duration; d > 0 {
				log.Printf("Adding a delay of: %s\n", d)
				time.Sleep(d)
			}

			response = &mock.Response
		}

		statistics.TrackMockRequest()
	} else {
		response = &mock.Response
	}

	match := match.NewTransaction(request, response, result)

	return mock, match

}

//Start initialize the HTTP mock server
func (di Dispatcher) Start() {
	addr := fmt.Sprintf("%s:%d", di.IP, di.Port)
	addrTLS := fmt.Sprintf("%s:%d", di.IP, di.PortTLS)

	errCh := make(chan error)

	go func() {
		errCh <- http.ListenAndServe(addr, &di)
	}()

	go func() {
		err := di.listenAndServeTLS(addrTLS)
		if err != nil {
			log.Println("Impossible start the application.")
			errCh <- err
		}
	}()

	err := <-errCh

	if err != nil {
		log.Fatalf("ListenAndServe: %s", err.Error())
	}

}

func (di Dispatcher) listenAndServeTLS(addrTLS string) error {
	tlsConfig := &tls.Config{}
	pattern := fmt.Sprintf("%s/*.crt", di.ConfigTLS)
	files, err := filepath.Glob(pattern)
	if err != nil || len(files) == 0 {
		log.Println("TLS certificates not found, impossible to start the TLS server.")
		return nil
	}

	for _, crt := range files {
		extension := filepath.Ext(crt)
		name := crt[0 : len(crt)-len(extension)]

		if filepath.Base(crt) == "ca.crt" {
			log.Println("Found ca cert", crt)
			ca, err := ioutil.ReadFile(crt)
			if err != nil {
				return fmt.Errorf("could not load CA Certificate '%s'", crt)
			}

			certPool := x509.NewCertPool()
			if ok := certPool.AppendCertsFromPEM(ca); !ok {
				return fmt.Errorf("could not append ca cert '%s' to CertPool", crt)
			}

			log.Println("Added ca certificate to server ", crt)
			tlsConfig.RootCAs = certPool
			continue
		}

		key := fmt.Sprint(name, ".key")
		log.Printf("Loading X509KeyPair (%s/%s)\n", filepath.Base(crt), filepath.Base(key))
		certificate, err := tls.LoadX509KeyPair(crt, key)
		if err != nil {
			return fmt.Errorf("Invalid certificate: %v", crt)
		}
		tlsConfig.Certificates = append(tlsConfig.Certificates, certificate)
	}
	tlsConfig.BuildNameToCertificate()

	server := http.Server{
		Addr:      addrTLS,
		TLSConfig: tlsConfig,
		Handler:   &di,
	}

	return server.ListenAndServeTLS("", "")
}
