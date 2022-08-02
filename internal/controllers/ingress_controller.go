package controllers

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"ngrok.io/ngrok-ingress-controller/pkg/agenttunneldriver"
	"ngrok.io/ngrok-ingress-controller/pkg/ngrokapidriver"
	"ngrok.io/ngrok-ingress-controller/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// This implements the Reconciler for the controller-runtime
// https://pkg.go.dev/sigs.k8s.io/controller-runtime#section-readme
type IngressReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
// +kubebuilder:rbac:groups="networking.k8s.io",resources=ingresses,verbs=get;list;watch;update
// +kubebuilder:rbac:groups="networking.k8s.io",resources=ingressclasses,verbs=get;list;watch

// This reconcile function is called by the controller-runtime manager.
// It is invoked whenever there is an event that occurs for a resource
// being watched (in our case, ingress objects). If you tail the controller
// logs and delete, update, edit ingress objects, you see the events come in.
func (t *IngressReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := t.Log.WithValues("ingress", req.NamespacedName)

	// Fetch the ingress class. Later on, we'll filter based on this.
	// I believe this client provided by the controller-runtime uses a cache
	// It might be better to do in the ManagerSetup with a filter later on
	ingressClasses := &netv1.IngressClassList{}
	if err := t.List(ctx, ingressClasses); err != nil {
		log.Error(err, "unable to list ingress classes")
		return ctrl.Result{}, err
	}
	log.Info(fmt.Sprintf("found %s ingress classes", ingressClasses))

	ingress := &netv1.Ingress{}
	if err := t.Get(ctx, req.NamespacedName, ingress); err != nil {
		if client.IgnoreNotFound(err) == nil {
			log.Info("ingress not found, must have been deleted")
			// TODO: We can't construct a full edge object anymore now that the ingress
			// object has been deleted. Unless we can infer every backend resource just based
			// on the initial naming scheme derived from the namespaced name, we'll have to
			// modify the ingress objects with a finalizer (see details in readme)
			return ctrl.Result{}, nil
		} else {

			log.Error(err, "unable to fetch ingress")
			return ctrl.Result{}, err
		}
	}

	log.Info(fmt.Sprintf("We did find the ingress. Lets ensure its created and up to date in the backend %+v \n", ingress))
	edge := newEdge(ingress)
	if err := agenttunneldriver.NewAgentTunnelDriver().EnsureTunnel(edge); err != nil {
		return ctrl.Result{}, err
	}
	if err := ngrokapidriver.NewApiDriver().EnsureEdge(edge); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (t *IngressReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&netv1.Ingress{}).
		Complete(t)
}

func newEdge(ingress *netv1.Ingress) *types.Edge {
	e := &types.Edge{
		Description: "Some meaningful description to make sense in the dashboard",
		MetaData: types.MetaData{
			"is owned by ingress controller": "a name that should be determined elsewhere",
			"some user provided":             "meta data for whatever they want",
		},
	}

	for _, r := range ingress.Spec.Rules {
		for _, p := range r.HTTP.Paths {
			e.Routes = append(e.Routes, types.Route{
				Path: p.Path,
				TunnelBackend: types.TunnelBackend{
					Hostname: p.Backend.Service.Name,
					Port:     int(p.Backend.Service.Port.Number),
				},
			})
		}

		if r.Host != "" {
			// modify edge to have a reserved domain
		}
	}

	if ingress.ObjectMeta.Annotations["ngrok.ingress.kubernetes.io/https-compression"] == "true" {
		// modify edge with https compression module
	}

	return e
}
