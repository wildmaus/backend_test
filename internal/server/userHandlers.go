package server

import (
	"backend_test/internal/model"
	"backend_test/pkg/utils"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func (s *server) updateUser(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling update user at %s\n", req.URL.Path)
	params, err := utils.ParseUintMass(req, "id", "amount")
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		return
	}
	id, amount := params[0], params[1]
	user, err := s.storage.User().FindOne(req.Context(), id)
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
	if err := s.storage.User().CreateWithTx(req.Context(), &user, &tx, create); err != nil {
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
	user, err := s.storage.User().FindOne(req.Context(), id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(404)
		return
	}
	utils.RenderJson(w, model.UserDto{Id: user.Id, Balance: user.Balance})
	w.WriteHeader(200)
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
	sortPrm, offset, err := parseQuery(req)
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		return
	}
	txs, err := s.storage.Transaction().FindTxByUser(req.Context(), id, sortPrm, offset)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
	utils.RenderJson(w, ResponseTx{Transactions: txs})
}

// get sort params and offset number
func parseQuery(req *http.Request) (string, int32, error) {
	by := req.FormValue("by")
	order := req.FormValue("order")
	page := req.FormValue("page")
	if by == "" {
		by = "date"
	} else if by != "date" && by != "amount" {
		return "", 0, errWrongInputData
	}
	if order == "" {
		order = "DESC"
	} else if strings.ToUpper(order) != "ASC" && strings.ToUpper(order) != "DESC" {
		return "", 0, errWrongInputData
	}
	if page == "" {
		page = "0"
	}
	offset, err := utils.ParseUint(page)
	if err != nil {
		return "", 0, errWrongInputData
	}
	return fmt.Sprintf("%v %v", by, order), offset * 5, nil
}
