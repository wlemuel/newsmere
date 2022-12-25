package nntp

import (
	"encoding/json"
	"fmt"
	"net"
	"net/textproto"
	"newsmere/internal/domain"
)

// nntpClient is an NNTP client.
type nntpClient struct {
	name         string
	conn         *textproto.Conn
	netconn      net.Conn
	tls          bool
	Banner       string
	capabilities []string
}

// NewBackend will create a backend contain nntpClient
func NewBackend(config json.RawMessage) (*domain.Backend, error) {
	backend := new(domain.Backend)
	err := json.Unmarshal(config, &backend)
	if err != nil {
		return nil, err
	}

	client, err := NewRepo(backend)
	if err != nil {
		return nil, err
	}

	backend.Repo = client
	return backend, nil
}

func NewRepo(b *domain.Backend) (domain.BackendRepository, error) {
	addr := fmt.Sprintf("%s:%d", b.Server, b.Port)
	if b.Port == 0 {
		addr = b.Server + ":119"
	}
	client, err := newClient("tcp", addr)
	if err != nil {
		return nil, err
	}

	client.name = b.Name

	if b.User != "" {
		_, err := client.Authenticate(b.User, b.Pass)
		if err != nil {
			return nil, err
		}
	}

	return client, nil
}

func (c *nntpClient) Start() error {
	return nil
}

func (c *nntpClient) Stop() error {
	return c.Close()
}

func (c *nntpClient) SyncSub(db domain.DBRepository) error {
	count, err := db.CountSubsByTypeRepo("nntp", c.name)
	if err != nil {
		return err
	}

	if count == 0 {
		var subs []domain.Subscription
		groups, err := c.List("")
		if err != nil {
			return err
		}

		for _, g := range groups {
			subs = append(subs, domain.Subscription{
				Name: g.Name,
				High: int(g.High),
				Low:  int(g.Low),
				Type: "nntp",
				Repo: c.name,
			})
		}

		err = db.StoreSubs(subs)
	}
	return nil
}
