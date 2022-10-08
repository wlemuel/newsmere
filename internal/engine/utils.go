package engine

import (
	"encoding/json"
	"fmt"
	nntp_bk "newsmere/internal/backend/nntp"
	"newsmere/internal/operator"
	nntp_sv "newsmere/internal/service/nntp"
)

func backendDecode(typeName string, config json.RawMessage) (Backend, error) {
	switch typeName {
	case nntp_bk.Type:
		return nntp_bk.New(config)
	default:
		return nil, fmt.Errorf(errUnknownBackendType, typeName)
	}
}

func serviceDecode(typeName string, config json.RawMessage) (Service, error) {
	switch typeName {
	case nntp_sv.Type:
		return nntp_sv.New(config, operator.New())
	default:
		return nil, fmt.Errorf(errUnknownServiceType, typeName)
	}
}
