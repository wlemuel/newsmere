package domain

type Service struct {
	Type string `json:"type"`
	Host string `json:"host"`
	Port int    `json:"port"`

	Delivery ServiceDelivery
}

type ServiceDelivery interface {
	Run() error
	Stop() error
}
