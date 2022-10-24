package server

import (
	"backend_test/pkg/utils"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *server) getTx(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling get tx at %s\n", req.URL.Path)
	id, err := utils.ParseUint(mux.Vars(req)["id"])
	if err != nil {
		w.WriteHeader(400)
		return
	}
	tx, err := s.storage.Transaction().FindOne(req.Context(), id)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	utils.RenderJson(w, tx)
}
