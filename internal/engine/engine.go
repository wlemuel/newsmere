package engine

import (
	"fmt"
	"newsmere/internal/db"
	"newsmere/internal/nntp"
	"newsmere/internal/types"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Backend interface {
	GetType() types.BackendType
	GetName() string
	IsRunning() bool

	Run() error
	Stop() error
	Restart() error
}

type Service struct {
	Type types.ServiceType
}

// Engine for managing the whole system.
type Engine struct {
	Backends []Backend
}

// AddService for adding service to engine.
func (e *Engine) AddBackend(b Backend) {
	e.Backends = append(e.Backends, b)
}

func (e *Engine) Run() {
	fmt.Println("Hello, Newsmere!")
	db.GetManager()
	e.startBackends()
	e.startServices()
}

func (e *Engine) Stop() {
	for _, b := range e.Backends {
		b.Stop()
	}
}

func (e *Engine) startBackends() {
	backend, err := nntp.NewClient("tcp", "127.0.0.1:1119")
	if err != nil {
		panic("failed to setup backends")
	}
	backend.Command("mode reader", 2)
	e.AddBackend(backend)
}

func (e *Engine) startServices() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-c:
			e.Stop()
			return
		case <-time.After(60 * time.Second):
			fmt.Println("Hello in a loop")
		}
	}
}
