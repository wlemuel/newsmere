package engine

import (
	"os"
	"testing"
)

func TestJSON(t *testing.T) {
	var (
		engine     = new(Engine)
		testConfig = "testdata/config.json"
	)

	jsonBytes, err := os.ReadFile(testConfig)
	if err != nil {
		t.Fatalf("Error reading config file: %s", testConfig)
	}

	err = engine.UnmarshalJSON(jsonBytes)
	if err != nil {
		t.Fatalf("Error unmarshaling: %v", err)
	}
}

func TestNew(t *testing.T) {
	e := New("testdata/config.json")
	t.Logf("Engine: %v", e)
}

func TestRun(t *testing.T) {
	e := New("testdata/config.json")

	for _, b := range e.Backends {
		err := b.Start()
		if err != nil {
			t.Fatalf("Error start backend (%s): %v", b.Type(), err)
		}
	}

	for _, s := range e.Services {
		err := s.Start()
		if err != nil {
			t.Fatalf("Error start service (%s): %v", s.Type(), err)
		}
	}
}
