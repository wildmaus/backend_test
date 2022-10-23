package server

import (
	"backend_test/internal/storage"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type server struct {
	router  *mux.Router
	storage storage.Storage
}

func newServer(storage storage.Storage) *server {
	s := &server{
		router:  mux.NewRouter(),
		storage: storage,
	}
	s.configureRouter()
	return s
}

func (s *server) configureRouter() {
	s.router.StrictSlash(true)
	s.router.HandleFunc("/", sayhello)

	s.router.HandleFunc("/user/{id:[0-9]+}/", s.getBalance).Methods("GET")
	s.router.HandleFunc("/user/{id:[0-9]+}/{amount:[0-9]+}/", s.updateUser).Methods("POST", "PUT")
	s.router.HandleFunc("/user/{id:[0-9]+}/tx/", s.getUserTx).Methods("GET")

	s.router.HandleFunc("/transfer/{fromId:[0-9]+}/{toId:[0-9]+}/{amount:[0-9]+}/", s.transfer).Methods("POST")

	s.router.HandleFunc("/tx/{id:[0-9]+}/", s.getTx).Methods("GET")

	s.router.HandleFunc("/reserve/", s.reserve).Methods("POST", "PUT")
	s.router.HandleFunc("/approve/", s.approve).Methods("POST", "PUT")
	s.router.HandleFunc("/cancel/", s.cancel).Methods("POST", "PUT")

	s.router.HandleFunc("/report/{month:[0-9]+}/{year:[0-9]+}/", s.getReport).Methods("GET")
}

func sayhello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Привет!</h1>")
}
