package ngrokapidriver

import (
	"context"
	"strings"

	"github.com/ngrok/ngrok-api-go/v4"
	tgb "github.com/ngrok/ngrok-api-go/v4/backends/tunnel_group"
	edge "github.com/ngrok/ngrok-api-go/v4/edges/https"
	edge_route "github.com/ngrok/ngrok-api-go/v4/edges/https_routes"

	"github.com/ngrok/ngrok-api-go/v4/reserved_domains"
)

type NgrokAPIDriver interface {
	FindEdge(ctx context.Context, id string) (*ngrok.HTTPSEdge, error)
	CreateEdge(ctx context.Context, e Edge) (*ngrok.HTTPSEdge, error)
	UpdateEdge(ctx context.Context, e Edge) (*ngrok.HTTPSEdge, error)
	DeleteEdge(ctx context.Context, e Edge) error
	GetReservedDomains(ctx context.Context, edgeID string) ([]ngrok.ReservedDomain, error)
}

type ngrokAPIDriver struct {
	edges           edge.Client
	tgbs            tgb.Client
	routes          edge_route.Client
	reservedDomains reserved_domains.Client
	metadata        string
}

func NewNgrokApiClient(apiKey string) NgrokAPIDriver {
	config := ngrok.NewClientConfig(apiKey)
	return &ngrokAPIDriver{
		edges:           *edge.NewClient(config),
		tgbs:            *tgb.NewClient(config),
		routes:          *edge_route.NewClient(config),
		reservedDomains: *reserved_domains.NewClient(config),
		metadata:        "\"{\"owned-by\":\"ngrok-ingress-controller\"}\"",
	}
}

func (nc ngrokAPIDriver) FindEdge(ctx context.Context, id string) (*ngrok.HTTPSEdge, error) {
	return nc.edges.Get(ctx, id)
}

// Goes through the whole edge object and creates resources for
// * reserved domains
// * tunnel group backends
// * edge routes
// * the edge itself
func (napi ngrokAPIDriver) CreateEdge(ctx context.Context, edgeSummary Edge) (*ngrok.HTTPSEdge, error) {
	// TODO: Support multiple rules and multiple hostports
	domain := strings.Split(edgeSummary.Hostport, ":")[0]
	_, err := napi.reservedDomains.Create(ctx, &ngrok.ReservedDomainCreate{
		Name:        domain,
		Region:      "us", // TODO: Set this from user config
		Description: "Created by ngrok-ingress-controller",
		Metadata:    napi.metadata,
	})
	// Swallow conflicts, just always try to create it
	// TODO: Depending on if we choose to clean up reserved domains or not, we may want to surface this conflict to the user
	if err != nil && !strings.Contains(err.Error(), "ERR_NGROK_413") && !strings.Contains(err.Error(), "ERR_NGROK_7122") {
		return nil, err
	}

	var newEdge *ngrok.HTTPSEdge

	// If the edge ID is already set, try to look it up
	if edgeSummary.Id != "" {
		newEdge, err = napi.edges.Get(ctx, edgeSummary.Id)
		if ngrok.IsNotFound(err) {
			edgeSummary.Id = ""
		} else if err != nil {
			return nil, err
		}
	} else { // Otherwise Make it
		newEdge, err = napi.edges.Create(ctx, &ngrok.HTTPSEdgeCreate{
			Hostports:   &[]string{edgeSummary.Hostport},
			Description: "Created by ngrok-ingress-controller",
			Metadata:    napi.metadata,
		})
		if err != nil {
			return nil, err
		}
	}

	// TODO: the object should have labels tied to routes and loop over the routes to create
	// a backend for each one, even if they are the same
	backend, err := napi.tgbs.Create(ctx, &ngrok.TunnelGroupBackendCreate{
		Labels:      edgeSummary.Labels,
		Description: "Created by ngrok-ingress-controller",
		Metadata:    napi.metadata,
	})
	if err != nil {
		return nil, err
	}

// loop through new edge routs
// construct a key from each remote route
// compare to our list
	// for each one thats not in our list delete
	// for each thats in our list but not remote, create it 
	// if it is in our list, 

	map["/+prefix+nginx"] = someInMemoryRoute
	map["/stats+prefix+haproxy"] = someInMemoryRoute

	// Loop through each route
	// compare each of our annotations to the route settings (like compression)
	// For each one if different, set it locally (do this for all annotations at the same time) and mark it diry
	// Update each route that was changed locally via the api 


	for _, r := range newEdge.Routes {
		externalRouteKey = r.Path + r.Prefix + r.Service
	}

// UPDATE STUFF
	for _, route := range edgeSummary.Routes {
		key := route.Match + route.MatchType + edgeSummary.Labels["k8s.ngrok.com/k8s-backend-name"] // TODO: Add in the namespace and port
		backendKey := edgeSummary.Labels["k8s.ngrok.com/k8s-backend-name"]  // TODO: Add in the namespace and port
		


		// CREATE STUFF
		r, err := napi.routes.Create(ctx, &ngrok.HTTPSEdgeRouteCreate{
			EdgeID:      newEdge.ID,
			MatchType:   route.MatchType,
			Match:       route.Match,
			Description: "Created by ngrok-ingress-controller",
			Metadata:    napi.metadata,
			Compression: "on", // TODO: make this based on annotation
			Backend: &ngrok.EndpointBackendMutate{
				BackendID: backend.ID,
			},
		})

		if r.Compression != edgeSummary.Annotations["k8s.ngrok.com/compression"] {
			napi.Routes.Update
		if err != nil {
			return nil, err
		}

		napi.
	}

	return newEdge, nil
}

// TODO: Implement this
func (nc ngrokAPIDriver) UpdateEdge(ctx context.Context, edgeSummary Edge) (*ngrok.HTTPSEdge, error) {
	// existinEdge := findEdge(edgeSummary.Id)
	// if existingEdge.HostPorts != edgeSummary.HostPorts {
	/// update the reserved domain stuff
	// update the edge's hostports
	// }
	// Check if metadata is different and if so override it

	return nil, nil
}

func (nc ngrokAPIDriver) DeleteEdge(ctx context.Context, e Edge) error {
	edge, err := nc.edges.Get(ctx, e.Id)
	if err != nil {
		return err
	}
	for _, route := range edge.Routes {
		err := nc.routes.Delete(ctx, &ngrok.EdgeRouteItem{EdgeID: e.Id, ID: route.ID})
		if err != nil {
			return err
		}
	}

	// TODO: I could delete the reserved endpoint, but it might make sense to just leave it reserved. Keeping this for now
	err = nc.edges.Delete(ctx, e.Id)
	if err != nil {
		if !ngrok.IsNotFound(err) {
			return err
		} else {
		}
	}
	return nil
}

func (nc ngrokAPIDriver) GetReservedDomains(ctx context.Context, edgeID string) ([]ngrok.ReservedDomain, error) {
	edge, err := nc.FindEdge(ctx, edgeID)
	if err != nil {
		return nil, err
	}
	hostPortDomains := []string{}
	for _, hostport := range *edge.Hostports {
		hostPortDomains = append(hostPortDomains, strings.Split(hostport, ":")[0])
	}

	domainsItr := nc.reservedDomains.List(nil)
	var matchingReservedDomains []ngrok.ReservedDomain
	// Loop while there are more domains and check if they match any of the hostPortDomains. If so add it to the reservedDomains
	for {
		if !domainsItr.Next(ctx) {
			err := domainsItr.Err()
			if err != nil {
				return nil, err
			}
			break
		}
		domain := domainsItr.Item()

		for _, hostPortDomain := range hostPortDomains {
			if domain.Domain == hostPortDomain {
				matchingReservedDomains = append(matchingReservedDomains, *domain)
			}
		}
	}

	return matchingReservedDomains, nil
}
