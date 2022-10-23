package server

import (
	"backend_test/internal/model"
	"backend_test/pkg/utils"
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func (s *server) updateUser(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling update user at %s\n", req.URL.Path)
	id, err := utils.ParseUint(mux.Vars(req)["id"])
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		return
	}
	amount, err := utils.ParseUint(mux.Vars(req)["amount"])
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		return
	}
	user, err := s.storage.User().FindOne(context.TODO(), id)
	var create bool
	if err != nil {
		user = model.User{Id: id, Balance: amount, Reserved: 0}
		create = true
	} else {
		if user.Balance+amount < user.Balance {
			w.WriteHeader(400)
			return
		}
		user.Balance += amount
		create = false
	}
	tx := model.TransactionDto{ToId: &id, Amount: amount, Date: time.Now(), Type: 0}
	if err := s.storage.User().CreateWithTx(context.TODO(), &user, &tx, create); err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(204)
}

func (s *server) getBalance(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling get balance at %s\n", req.URL.Path)
	id, err := utils.ParseUint(mux.Vars(req)["id"])
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		return
	}
	user, err := s.storage.User().FindOne(context.TODO(), id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(404)
		return
	}
	w.WriteHeader(200)
	utils.RenderJson(w, model.UserDto{Id: user.Id, Balance: user.Balance})
}

func (s *server) getUserTx(w http.ResponseWriter, req *http.Request) {
	type ResponseTx struct {
		Transactions []model.Transaction `json:"transactions"`
	}
	log.Printf("handling get user's tx at %s\n", req.URL.Path)
	id, err := utils.ParseUint(mux.Vars(req)["id"])
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		return
	}
	by := req.FormValue("by")
	order := req.FormValue("order")
	offset := req.FormValue("offset")
	if by == "" {
		by = "date"
	} else if by != "date" && by != "amount" {
		w.WriteHeader(400)
		return
	}
	if order == "" {
		order = "DESC"
	} else if strings.ToUpper(order) != "ASC" && strings.ToUpper(order) != "DESC" {
		w.WriteHeader(400)
		return
	}
	if offset == "" {
		offset = "0"
	}
	oft, err := utils.ParseUint(offset)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	log.Println(fmt.Sprintf("%v %v", by, order))
	txs, err := s.storage.Transaction().FindTxByUser(context.TODO(), id, fmt.Sprintf("%v %v", by, order), oft*5)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
	utils.RenderJson(w, ResponseTx{Transactions: txs})
}
