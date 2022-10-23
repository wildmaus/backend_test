package model

type Details struct {
	Id        int32 `json:"id"`
	OrderId   int32 `json:"orderId"`
	ServiceId int32 `json:"serviceId"`
	Status    bool  `json:"status"`
}
