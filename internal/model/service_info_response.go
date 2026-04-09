package model

type ServiceInfoResponse struct {
	Service string `json:"service"`
	Status  string `json:"status"`
	Message string `json:"message"`
}
