package cqrs

import (
	"errors"
	"fmt"
)

// TryMapPayload tries to map the payload to the given type. It returns an error if the payload is nil or if the payload is
// not of the given type.
func TryMapPayload[T Payload](payload Payload) (*T, error) {
	if payload == nil {
		return nil, errors.New("payload is nil")
	}

	v, ok := payload.(T)
	if !ok {
		return nil, fmt.Errorf("payload is not of type %T", v)
	}

	return &v, nil
}
