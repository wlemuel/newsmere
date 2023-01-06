package api

import (
	"fmt"
	"newsmere/internal/domain"

	"github.com/gin-gonic/gin"
)

type apiService struct {
	host string
	port int

	core *gin.Engine
}

func NewDelivery(s *domain.Service) (domain.ServiceDelivery, error) {
	g := gin.Default()

	service := &apiService{
		host: s.Host,
		port: s.Port,
		core: g,
	}
	service.setup()

	return service, nil
}

func (c *apiService) Run() error {
	addr := fmt.Sprintf("%s:%d", c.host, c.port)
	return c.core.Run(addr)
}

func (c *apiService) Stop() error {
	return nil
}

func (c *apiService) setup() error {
	c.core.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	return nil
}
