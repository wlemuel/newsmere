package engine

import (
	"newsmere/internal/domain"
	_nntpRepo "newsmere/internal/repository/nntp"
	"newsmere/internal/repository/sqlite"
	_restApi "newsmere/internal/service/api"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var conf domain.Config

func init() {
	viper.SetConfigFile("config.json")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	viper.Unmarshal(&conf)

	log.SetOutput(os.Stdout)
	if conf.Debug {
		log.Println("Service RUN on DEBUG mode")
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}
}

type Engine struct {
	storage  domain.DBRepository
	backends []domain.Backend
	services []domain.Service
}

func New() Engine {
	return Engine{}
}

func (e *Engine) Run() {
	e.ensureStorage()
	e.startBackends()
	e.startServices()
}

func (e *Engine) ensureStorage() {
	if e.storage == nil {
		e.storage = sqlite.NewSqliteRepo()
	}
}

func (e *Engine) startBackends() {
	for _, b := range conf.Backends {
		if b.Type == "nntp" {
			repo, err := _nntpRepo.NewRepo(&b)
			if err != nil {
				log.Warnln(err)
				continue
			}
			b.Repo = repo
			e.backends = append(e.backends, b)
		}
	}
}

func (e *Engine) stopBackends() {
	for _, b := range e.backends {
		b.Repo.Stop()
	}
}

func (e *Engine) startServices() {
	for _, s := range conf.Services {
		if s.Type == "api" {
			delivery, err := _restApi.NewDelivery(&s)
			if err != nil {
				log.Warnln(err)
				continue
			}
			s.Delivery = delivery
			e.services = append(e.services, s)

			delivery.Run()
		}
	}
}

func (e *Engine) stopServices() {
	for _, s := range e.services {
		s.Delivery.Stop()
	}
}
