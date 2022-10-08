package interfaces

type Backend interface {
	Type() string

	Start() error
	Stop() error
	Restart() error
	Status() string
}

type Service interface {
	Type() string

	Start() error
	Stop() error
	Restart() error
	Status() string
}
