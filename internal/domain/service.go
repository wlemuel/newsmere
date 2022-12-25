package domain

type Service struct {
	Type string `json:"type"`
	Host string `json:"host"`
	Port int    `json:"port"`
}

type ServiceDelivery interface {
	Run() error
}