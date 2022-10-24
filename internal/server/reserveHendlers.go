package server

import (
	"backend_test/internal/model"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
)

type requestReserve struct {
	Id        int32 `json:"userId"`
	OrderId   int32 `json:"orderId"`
	ServiceId int32 `json:"serviceId"`
	Amount    int32 `json:"amount"`
}

var (
	errWrongInputData = errors.New("wrong json input data")
	errStatus         = errors.New("already approved or canceled")
	errAmountNotEq    = errors.New("request amount not eq reserved")
	errAmountExcced   = errors.New("amount excced balance/reserved")
)

func (s *server) reserve(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling reserve at %s\n", req.URL.Path)
	request := requestReserve{}
	if err := checkInput(&request, req); err != nil {
		log.Println(err)
		w.WriteHeader(400)
		return
	}
	user, err := s.storage.User().FindOne(req.Context(), request.Id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(404)
		return
	}
	if user.Balance < request.Amount || user.Reserved > user.Reserved+request.Amount {
		log.Println(err)
		w.WriteHeader(400)
		return
	}
	user.Balance -= request.Amount
	user.Reserved += request.Amount
	tx := model.TransactionDto{FromId: &request.Id, Amount: request.Amount, Date: time.Now(), Type: 2}
	details := model.Details{OrderId: request.OrderId, ServiceId: request.ServiceId, Status: false}
	if err := s.storage.Details().Reserve(req.Context(), &user, &tx, &details); err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(204)
}

func (s *server) cancel(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling cancel at %s\n", req.URL.Path)
	request := requestReserve{}
	s.checkReserved(&request, w, req, true)
}

func (s *server) approve(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling approve at %s\n", req.URL.Path)
	request := requestReserve{}
	s.checkReserved(&request, w, req, false)
}

// common in approve/cancel
func (s *server) checkReserved(r *requestReserve, w http.ResponseWriter, req *http.Request, cancel bool) {
	if err := checkInput(r, req); err != nil {
		w.WriteHeader(400)
		return
	}
	// find row with details
	details, err := s.storage.Details().FindOne(req.Context(), r.OrderId, r.ServiceId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(404)
		return
	}
	if details.Status != false {
		log.Println(errStatus)
		w.WriteHeader(400)
		return
	}
	// find tx by userId, type and detailsId
	amount, err := s.storage.Transaction().FindByDetails(req.Context(), r.Id, details.Id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(404)
		return
	}
	if amount != r.Amount {
		log.Println(errAmountNotEq)
		w.WriteHeader(400)
		return
	}
	// get user and recalculate reserved
	user, err := s.storage.User().FindOne(req.Context(), r.Id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(404)
		return
	}
	if user.Reserved < r.Amount {
		log.Println(errAmountExcced)
		w.WriteHeader(400)
		return
	}
	user.Reserved -= amount
	var tx model.TransactionDto
	if cancel {
		if user.Balance+amount < user.Balance {
			log.Println(errAmountExcced)
			w.WriteHeader(400)
			return
		}
		user.Balance += amount
		// to = id, from = null
		tx = model.TransactionDto{ToId: &r.Id, Amount: r.Amount, Date: time.Now(), Type: 3, DetailsId: &details.Id}
	} else {
		// from = id, to = null
		tx = model.TransactionDto{FromId: &r.Id, Amount: r.Amount, Date: time.Now(), Type: 4, DetailsId: &details.Id}
	}
	if err := s.storage.Details().SolveReserve(req.Context(), &user, &tx); err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(204)
}

// check json data and fill struct
func checkInput(r *requestReserve, req *http.Request) error {
	if err := json.NewDecoder(req.Body).Decode(r); err != nil {
		return errWrongInputData
	}
	if r.Id < 0 || r.OrderId < 0 || r.ServiceId < 0 || r.Amount < 0 {
		return errWrongInputData
	}
	return nil
}
