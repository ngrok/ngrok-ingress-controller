package agenttunneldriver

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	"ngrok.io/ngrok-ingress-controller/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Edge = types.Edge

type AgentTunnelDriver interface {
	EnsureTunnel(ctx context.Context, edge *Edge) error
	RemoveTunnel(ctx context.Context, edge *Edge) error
}

func NewAgentTunnelDriver(c client.Client) AgentTunnelDriver {
	return agentTunnelDriver{
		k8sClient: c,
	}
}

type agentTunnelDriver struct {
	k8sClient client.Client
}

// These just need to issue requests to the agent api to start and stop tunnels
func (atd agentTunnelDriver) EnsureTunnel(ctx context.Context, edge *Edge) error {
	log, err := logr.FromContext(ctx)
	if err != nil {
		return err
	}
	configMapName := "ngrok-ingress-controller-agent-config"
	configMapKey := "agent.yaml"
	config := &v1.ConfigMap{}
	ac := &agentConfig{}
	// Try to find the existing config map
	err = atd.k8sClient.Get(ctx, client.ObjectKey{Name: configMapName, Namespace: atd.k8sClient.Namespace}, config)
	if err == nil {
		if _, ok := config.Data[configMapKey]; ok {
			ac = &agentConfig{
				tunnels: []*tunnel{
					newTunnel("name", "addr"),
				},
			}
		} else {
			panic("Config map is missing the key " + configMapKey + " which shouldn't be possible")
		}
	}

	log.Info(fmt.Sprintf("Agent config is %+v\n", ac)) // TODO: Plumb logger through ctx
	return nil
}

func (atd agentTunnelDriver) RemoveTunnel(ctx context.Context, edge *Edge) error {
	return nil
}

type agentConfig struct {
	tunnels []*tunnel
}

func newTunnel(name string, addr string) *tunnel {
	return &tunnel{
		name:  name,
		addr:  addr,
		proto: "http",
	}
}

type tunnel struct {
	name     string
	addr     string
	proto    string
	hostname string
}
