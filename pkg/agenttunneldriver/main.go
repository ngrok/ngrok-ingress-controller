package agenttunneldriver

import (
	"ngrok.io/ngrok-ingress-controller/pkg/types"
)

type Edge = types.Edge

type AgentTunnelDriver interface {
	EnsureTunnel(edge *Edge) error
	RemoveTunnel(edge *Edge) error
}

func NewAgentTunnelDriver() AgentTunnelDriver {
	return agentTunnelDriver{}
}

type agentTunnelDriver struct {
}

// These just need to issue requests to the agent api to start and stop tunnels
func (atd agentTunnelDriver) EnsureTunnel(edge *Edge) error {
	return nil
}

func (atd agentTunnelDriver) RemoveTunnel(edge *Edge) error {
	return nil
}
