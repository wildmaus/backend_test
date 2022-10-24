package model

import "time"

type TransactionDto struct {
	Id        int32     `json:"id"`
	FromId    *int32    `json:"fromId"`
	ToId      *int32    `json:"toId"`
	Amount    int32     `json:"amount"`
	Date      time.Time `json:"date"`
	Type      int32     `json:"type"`
	DetailsId *int32    `json:"detailsId"`
}
