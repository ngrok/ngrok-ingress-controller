package ngrokapidriver

import (
	"ngrok.io/ngrok-ingress-controller/pkg/types"
)

type Edge = types.Edge

type ApiDriver interface {
	EnsureEdge(edge *Edge) error
	RemoveEdge(edge *Edge) error
}

func NewApiDriver() ApiDriver {
	return apiDriver{}
}

type apiDriver struct {
}

// For both of these, we want to basically walk through the edge
// data and create/update/delete backend resources each pieces of data represents
func (ad apiDriver) EnsureEdge(edge *Edge) error {
	return nil
}

func (ad apiDriver) RemoveEdge(edge *Edge) error {
	return nil
}
