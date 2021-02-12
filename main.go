package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	OK 		= "OK"
	BAD 	= "BAD"
	TIMEOUT = "TIMEOUT"
)

type State struct {
	status string
}

func NewState() *State {
	return &State{status: OK}
}

func (s *State) Health(rw http.ResponseWriter, r *http.Request) {
	log.Printf("Received /health request: source=%v status=%v", r.RemoteAddr, s.status)
	switch s.status {
	case OK:
		io.WriteString(rw, "I'm healthy")
	case BAD:
		http.Error(rw, "Internal Error", 500)
	case TIMEOUT:
		time.Sleep(30 * time.Second)
	default:
		io.WriteString(rw, "UNKNOWN")
	}
}

func (s *State) Sabotage(rw http.ResponseWriter, r *http.Request) {
	s.status = BAD
	io.WriteString(rw, "Sabotage ON")
}

func (s *State) Recover(rw http.ResponseWriter, r *http.Request) {
	s.status = OK
	io.WriteString(rw, "Recovered.")
}

func (s *State) Timeout(rw http.ResponseWriter, r *http.Request) {
	s.status = TIMEOUT
	io.WriteString(rw, "Configured to timeout.")
}


func main() {
	httpState := NewState()
	mux := http.NewServeMux()
	mux.HandleFunc("/health", httpState.Health)
	mux.HandleFunc("/sabotage", httpState.Sabotage)
	mux.HandleFunc("/recover", httpState.Recover)
	mux.HandleFunc("/timeout", httpState.Timeout)
	log.Print("Starting http server")
	err := http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), mux)
	log.Fatal(err)
}

