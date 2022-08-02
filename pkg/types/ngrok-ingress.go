package types

// The goal here is to have 1 big object representing all of the possible
// data represented in a service/ingress/gateway/crd based implementation.
// The goal is for the k8s aware code to handle converting each controllers'
// specific resources and converting everything to a common format.
// The goal being that we can modify/add controllers and evolve our ngrok
// api independently
//
// These are in no way complete, and likely not right. Just a concept.
// Edge may be a bad name. I'm open to suggestions here.
type Edge struct {
	Description string
	MetaData    MetaData
	// hostports?
	ReservedDomain ReservedDomain
	Routes         []Route
	Modules        []Module
}

type ReservedDomain struct {
	Name string
}

type TunnelBackend struct {
	Hostname string
	Port     int
}

type Route struct {
	Path          string
	TunnelBackend TunnelBackend
}

type Module struct {
	Name string
}

type MetaData = map[string]string
