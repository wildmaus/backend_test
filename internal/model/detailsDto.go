package model

type DetailsDto struct {
	Id        int32 `json:"userId"`
	OrderId   int32 `json:"orderId"`
	ServiceId int32 `json:"serviceId"`
	Amount    int32 `json:"amount"`
}
