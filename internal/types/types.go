package types

type BackendType int8

const (
	NNTP BackendType = iota
	IMAP
	RSS
)

type Backend interface {
	GetType() BackendType
	GetName() string
	IsRunning() bool

	Run() error
	Stop() error
	Restart() error
}

type ServiceType int8

const (
	Nntp ServiceType = iota
	Api
	GraphQL
	Web
)

type Service struct {
	Type ServiceType
}
