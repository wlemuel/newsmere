package domain

// Backend represents a backend for providing data resource.
type Backend struct {
	Name   string `json:"name"`
	User   string `json:"user,omitempty"`
	Pass   string `json:"pass,omitempty"`
	Server string `json:"server"`
	Port   int    `json:"port,omitempty"`

	Repo BackendRepository
}

type BackendRepository interface {
	Start() error
	Stop() error
	SyncSub(db DBRepository) error
}
