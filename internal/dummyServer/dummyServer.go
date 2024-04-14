package dummyServer

import (
    "fmt"
    "net/http"
    "context"
    "sync"
    "github.com/jmartin82/mmock/v3/internal/config/logger"
)

var log = logger.Log

type DummyServer struct {
   Port int
   Srv *http.Server
}


func (ds DummyServer) hello(w http.ResponseWriter, req *http.Request) {

    fmt.Fprintf(w, "hello\n")
}

func (ds DummyServer) headers(w http.ResponseWriter, req *http.Request) {

    for name, headers := range req.Header {
        for _, h := range headers {
            fmt.Fprintf(w, "%v: %v\n", name, h)
        }
    }
}

func Start(wg *sync.WaitGroup, port int) DummyServer {
  ds := DummyServer{Port: port, Srv: &http.Server{Addr: fmt.Sprintf(":%v", port)}}
  log.Debugf("DUMSSTART dr after assign: %v", ds)

    http.HandleFunc("/hello", ds.hello)
    http.HandleFunc("/headers", ds.headers)

    go func() {
      wg.Done() 

      log.Debugf("DUMSSTART after Done: %v", ds)
      if err := ds.Srv.ListenAndServe(); err != http.ErrServerClosed {
	// unexpected error. port in use?
	log.Fatalf("ListenAndServe(): %v", err)
      }
    }()

     
    return ds
   }

func (ds DummyServer) Stop(){
  log.Infof("DummyServer.Stop ds.Srv: %v", ds.Srv)
  if err := ds.Srv.Shutdown(context.TODO()); err != nil {
        panic(err) // failure/timeout shutting down the server gracefully
   }
}
