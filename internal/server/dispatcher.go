package server

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jmartin82/mmock/v3/internal/config/logger"
	"github.com/jmartin82/mmock/v3/internal/proxy"
	"github.com/jmartin82/mmock/v3/internal/statistics"
	"github.com/jmartin82/mmock/v3/pkg/match"
	"github.com/jmartin82/mmock/v3/pkg/mock"
	"github.com/jmartin82/mmock/v3/pkg/vars"
)

var log = logger.Log

// Dispatcher is the mock http server
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
	log.Infof("Running webhook: %s\n", url)
	content, err := json.Marshal(match)
	if err != nil {
		log.Error("Impossible encode the WebHook payload")
		return
	}
	reader := bytes.NewReader(content)
	resp, err := http.Post(url, "application/json", reader)
	if err != nil {
		log.Errorf("Impossible send payload to: %s\n", url)
		return
	}
	log.Infof("WebHook response: %d\n", resp.StatusCode)
}

// ServerHTTP is the mock http server request handler.
// It uses the router to decide the matching mock and translator as adapter between the HTTP impelementation and the mock mock.
func (di *Dispatcher) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	mRequest := di.Translator.BuildRequestDefinitionFromHTTP(req)

	if mRequest.Path == "/favicon.ico" {
		return
	}

	log.Infof("%s %s %s\n", terminal.Info("New request:"), req.Method, req.URL.String())

	mock, transaction := di.getMatchingResult(&mRequest)

	//save the match info
	di.Spier.Save(*transaction)

	//set new scenario
	if mock.Control.Scenario.NewState != "" {
		statistics.TrackScenarioFeature()

		di.Scenario.SetState(
		  mock.Control.Scenario.Name,
		  mock.Control.Scenario.NewState)

		if len(mock.Control.Scenario.Values) != 0 {
		   di.Scenario.SetStateValues(
		     mock.Control.Scenario.Name,
		     mock.Control.Scenario.Values)
		 }
	}

	if mock.Control.WebHookURL != "" {
		go di.callWebHook(mock.Control.WebHookURL, transaction)
	}

	//translate request
	di.Translator.WriteHTTPResponseFromDefinition(transaction.Response, w)

	if mock.Callback.Url != "" {
		go func() {
			_, err := HandleCallback(mock.Callback)
			if err != nil {
				log.Errorf("Error from HandleCallback: %s", err)
			} else {
				log.Info("Callback made successfully")
			}
		}()
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
	var response *mock.Response

	mock, result := di.Resolver.Resolve(request)
	colorTerminalOutput(result, mock)

	if result.Found {
		if len(mock.Control.ProxyBaseURL) > 0 {
			statistics.TrackProxyFeature()
			response = getProxyResponse(request, mock)
		} else {
			di.Evaluator.Eval(request, mock)
			response = &mock.Response
		}

		if mock.Control.Crazy {
			log.Info("Running crazy mode")
			response.StatusCode = di.randomStatusCode(response.StatusCode)
		}
		if d := mock.Control.Delay.Duration; d > 0 {
			log.Infof("Adding a delay of: %s\n", d)
			time.Sleep(d)
		}

		statistics.TrackMockRequest()
	} else {
		response = &mock.Response
	}

	match := match.NewTransaction(request, response, result)

	return mock, match

}

func colorTerminalOutput(result *match.Result, mock *mock.Definition) {
	msg := fmt.Sprintf("%s %s\n", terminal.Error("Definition match found:"), strconv.FormatBool(result.Found))
	if result.Found {
		msg = fmt.Sprintf("%s %s %s %s\n", terminal.Success("Definition match found:"), strconv.FormatBool(result.Found), terminal.Success("Name:"), mock.URI)
	}
	log.Info(msg)
}

// Start initialize the HTTP mock server
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
			log.Errorf("Impossible start the application.")
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
		log.Info("TLS certificates not found, impossible to start the TLS server.")
		return nil
	}

	for _, crt := range files {
		extension := filepath.Ext(crt)
		name := crt[0 : len(crt)-len(extension)]

		if filepath.Base(crt) == "ca.crt" {
			log.Info("Found ca cert", crt)
			ca, err := ioutil.ReadFile(crt)
			if err != nil {
				return fmt.Errorf("could not load CA Certificate '%s'", crt)
			}

			certPool := x509.NewCertPool()
			if ok := certPool.AppendCertsFromPEM(ca); !ok {
				return fmt.Errorf("could not append ca cert '%s' to CertPool", crt)
			}

			log.Infof("Added ca certificate to server ", crt)
			tlsConfig.RootCAs = certPool
			continue
		}

		key := fmt.Sprint(name, ".key")
		log.Infof("Loading X509KeyPair (%s/%s)\n", filepath.Base(crt), filepath.Base(key))
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
