package server

import (
	"backend_test/internal/storage"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type server struct {
	router  *mux.Router
	storage storage.Storage
	fs      http.Handler
}

func newServer(storage storage.Storage) *server {
	err := os.Mkdir("/download", 0755)
	if err != nil && !os.IsExist(err) {
		log.Println(err)
	}
	s := &server{
		router:  mux.NewRouter(),
		storage: storage,
	}
	s.configureRouter()
	return s
}

func (s *server) configureRouter() {
	s.router.HandleFunc("/", sayhello)

	s.router.HandleFunc("/user/{id:[0-9]+}", s.getBalance).Methods("GET")
	s.router.HandleFunc("/user/{id:[0-9]+}/{amount:[0-9]+}", s.updateUser).Methods("POST", "PUT")
	s.router.HandleFunc("/user/{id:[0-9]+}/tx", s.getUserTx).Queries("sort_by", "{sort_by}", "sort_order", "{sort_order}", "page", "{[0-9]+}").Methods("GET")
	s.router.HandleFunc("/user/{id:[0-9]+}/tx", s.getUserTx).Methods("GET")
	s.router.HandleFunc("/transfer/{fromId:[0-9]+}/{toId:[0-9]+}/{amount:[0-9]+}", s.transfer).Methods("POST")

	s.router.HandleFunc("/tx/{id:[0-9]+}", s.getTx).Methods("GET")

	s.router.HandleFunc("/reserve", s.reserve).Methods("POST", "PUT")
	s.router.HandleFunc("/approve", s.approve).Methods("POST", "PUT")
	s.router.HandleFunc("/cancel", s.cancel).Methods("POST", "PUT")

	s.router.HandleFunc("/report/{month:[0-9]+}/{year:[0-9]+}", s.getReport).Methods("GET")
	s.router.HandleFunc("/download/{filename}", s.download).Methods("GET")
}

func sayhello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Привет!</h1>")
}
