# Ngrok Ingress Controller Architecture and Design

This document is meant to describe the overall architecture and design of this controller. This ranges from high level topics like how its individual controller-runtimes and manager work, to low level design decisions that have been made.

## What is this?

The Ngrok Ingress Controller is an implementation of a k8s [controller](todo). Controllers are meant to watch some k8s resource and react any changes to that resource. A controller will then usually manage some other resource either in k8s or external. In this case, the controller watches for changes to Ingress resources and then manages Ngrok Tunnel resources on agents in the cluster, and Ngrok [edge api](todo) resources.

### More Details

Internally, the ingress controller is actually made up of 2 separate controllers that are managed by a single manager. The first controller is the `TunnelController` which watches for changes to Ingress resources and manages the Ngrok Tunnel resources. The second controller is the `IngressController` which watches for changes to Tunnel resources and manages the Ngrok Edge API resources. The manager handles leader election as well as shared cache clients and such. It makes it so only the ingress controller cares about being the leader.

Each controller implements a reconcile loop that reacts to any k8s events on watched objects.


## Managing Ngrok API Resources

As mentioned, only 1 leader responds to ingress objects and manages api resources. This does a basic mapping of things like hosts, routes, backend, etc into ngrok edges, routes, labeled tunnels, and reserved domains.
The edge id is stored on the ingress object so subsequent reconciles can update the edge.

Changing the hostname creates a brand new edge and has momentary downtime (although the url is changing anyways so.

A single ingress maps to a single edge.

Reserved domains are never deleted.

### Back Propagating Changes

We update things like status, finalizers, edge-id when things change.


## Public Resources
* containers, helm chart, and docs


## FAQ
