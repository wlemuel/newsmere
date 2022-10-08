package engine

import (
	"encoding/json"
	"log"
	"os"
)

// Engine for managing the whole system.
type Engine struct {
	Backends []Backend `json:"backends"`
	Services []Service `json:"services"`
}

func New(configFile string) Engine {
	if configFile == "" {
		configFile = "config.json"
	}

	configBytes, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}

	var e Engine
	err = json.Unmarshal(configBytes, &e)
	if err != nil {
		log.Fatal(err)
	}
	return e
}

func (e *Engine) UnmarshalJSON(b []byte) error {
	type engine2 *Engine
	_ = json.Unmarshal(b, engine2(e))

	// clean the state
	e.Backends = []Backend{}
	e.Services = []Service{}

	raw := struct {
		Backends []json.RawMessage `json:"backends"`
		Services []json.RawMessage `json:"services"`
	}{}

	err := json.Unmarshal(b, &raw)
	if err != nil {
		return err
	}

	configTypes := struct {
		Backends []struct {
			Type string `json:"type"`
		}
		Services []struct {
			Type string `json:"type"`
		}
	}{}
	err = json.Unmarshal(b, &configTypes)
	if err != nil {
		return err
	}

	for i, b := range configTypes.Backends {
		backend, err := backendDecode(b.Type, raw.Backends[i])
		if err != nil {
			return err
		}
		e.Backends = append(e.Backends, backend)
	}

	for i, s := range configTypes.Services {
		service, err := serviceDecode(s.Type, raw.Services[i])
		if err != nil {
			return err
		}
		e.Services = append(e.Services, service)
	}
	return nil
}

func (e *Engine) Run() error {
	for _, b := range e.Backends {
		err := b.Start()
		if err != nil {
			return err
		}
	}

	for _, s := range e.Services {
		err := s.Start()
		if err != nil {
			return err
		}
	}

	return nil
}
