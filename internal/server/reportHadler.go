package server

import (
	"backend_test/pkg/utils"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func (s *server) getReport(w http.ResponseWriter, req *http.Request) {
	type ResponseReport struct {
		Report string `json:"report"`
	}
	log.Printf("handling get report at %s\n", req.URL.Path)
	params, err := utils.ParseUintMass(req, "month", "year")
	if err != nil {
		w.WriteHeader(400)
		return
	}
	month, year := params[0], params[1]
	if month > 12 || month < 1 {
		w.WriteHeader(400)
		return
	}
	report, err := s.storage.Transaction().GetReport(req.Context(), month, year)
	if err != nil {
		log.Print(err)
		w.WriteHeader(500)
		return
	}
	output, err := os.Create(fmt.Sprintf("/download/Report-%v-%v.csv", month, year))
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
	defer output.Close()
	writer := csv.NewWriter(output)
	defer writer.Flush()

	header := []string{"service id", "amount"}
	if err := writer.Write(header); err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
	for _, val := range report {
		var csvRow []string
		csvRow = append(csvRow, fmt.Sprint(val.ServiseId), fmt.Sprint(val.Amount))
		if err := writer.Write(csvRow); err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}
	}
	utils.RenderJson(w, ResponseReport{Report: fmt.Sprintf("http://localhost:8000/download/Report-%v-%v.csv", month, year)})
}

func (s *server) download(w http.ResponseWriter, req *http.Request) {
	filename := mux.Vars(req)["filename"]
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "text/csv")
	http.ServeFile(w, req, fmt.Sprintf("/download/%s", filename))
}
