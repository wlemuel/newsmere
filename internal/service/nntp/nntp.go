package nntp

import (
	"encoding/json"
	"fmt"
	"net"
	"newsmere/internal/types"
)

// MessageID provides convenient access to the article's Message ID.
func (a *Article) MessageID() string {
	return a.Header.Get("Message-Id")
}

func New(config json.RawMessage, operator Operator) (*Service, error) {
	service := new(Service)
	err := json.Unmarshal(config, &service)
	service.operator = operator
	return service, err
}

func (Service) Type() string {
	return Type
}

func (s *Service) Start() error {
	host := s.Host
	if host == "" {
		host = "127.0.0.1"
	}

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", host, s.Port))
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}

	fmt.Printf("[Service] %s listen at: %s:%d\n", s.Type(), host, s.Port)

	s.server = NewServer(s.operator)

	for {
		c, err := listener.AcceptTCP()
		if err != nil {
			return err
		}
		go s.server.Process(c)
	}
}

func (s *Service) Stop() error {
	return nil
}

func (s *Service) Restart() error {
	return nil
}

func (s *Service) Status() string {
	return types.StatusUp
}
