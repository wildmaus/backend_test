package model

type User struct {
	Id       int32 `json:"id"`
	Balance  int32 `json:"balance"`
	Reserved int32 `json:"reserved"`
}
