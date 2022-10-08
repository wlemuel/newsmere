package imap

import "encoding/json"

func New(config json.RawMessage) (*Backend, error) {
	backend := new(Backend)
	err := json.Unmarshal(config, &backend)
	return backend, err
}

func (Backend) Type() string {
	return Type
}
