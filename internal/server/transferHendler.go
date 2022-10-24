package server

import (
	"backend_test/internal/model"
	"backend_test/pkg/utils"
	"log"
	"net/http"
	"time"
)

func (s *server) transfer(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling transfer at %s\n", req.URL.Path)
	params, err := utils.ParseUintMass(req, "fromId", "toId", "amount")
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		return
	}
	fromId, toId, amount := params[0], params[1], params[2]
	from, err := s.storage.User().FindOne(req.Context(), fromId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(404)
		return
	}
	if from.Balance < amount {
		w.WriteHeader(400)
		return
	}
	from.Balance -= amount
	to, err := s.storage.User().FindOne(req.Context(), toId)
	var create bool
	if err != nil {
		create = true
		to = model.User{Id: toId, Balance: amount, Reserved: 0}
	} else {
		if to.Balance+amount < to.Balance {
			w.WriteHeader(400)
			return
		}
		create = false
		to.Balance += amount
	}
	tx := model.TransactionDto{FromId: &fromId, ToId: &toId, Amount: amount, Date: time.Now(), Type: 1}
	if err := s.storage.Transaction().Transfer(req.Context(), &from, &to, &tx, create); err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(204)
}
