package server

import (
	"backend_test/internal/model"
	"backend_test/pkg/utils"
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *server) getReport(w http.ResponseWriter, req *http.Request) {
	type ResponseReport struct {
		Report []model.Report `json:"report"`
	}
	log.Printf("handling get report at %s\n", req.URL.Path)
	month, err := utils.ParseUint(mux.Vars(req)["month"])
	if err != nil {
		w.WriteHeader(400)
		return
	}
	year, err := utils.ParseUint(mux.Vars(req)["year"])
	if err != nil {
		w.WriteHeader(400)
		return
	}
	if month > 12 || month < 1 {
		w.WriteHeader(400)
		return
	}
	report, err := s.storage.Transaction().GetReport(context.TODO(), month, year)
	if err != nil {
		log.Print(err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
	utils.RenderJson(w, ResponseReport{Report: report})
}
