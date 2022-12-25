package api

import (
	"fmt"
	"newsmere/internal/domain"

	"github.com/gin-gonic/gin"
)

type apiService struct {
	host string
	port int

	engine *gin.Engine
}

func NewDelivery(s *domain.Service) (domain.ServiceDelivery, error) {
	g := gin.Default()

	service := &apiService{
		host:   s.Host,
		port:   s.Port,
		engine: g,
	}

	return service, nil
}

func (c *apiService) Run() error {
	addr := fmt.Sprintf("%s:%d", c.host, c.port)
	return c.engine.Run(addr)
}
