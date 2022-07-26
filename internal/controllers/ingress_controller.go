package controllers

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
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
	// TODO: Figure out the best way to form the edgeName taking into account isolating multiple clusters
	edgeName := strings.Replace(req.NamespacedName.String(), "/", "-", -1)

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
			return ctrl.Result{}, nil
		} else {

			log.Error(err, "unable to fetch ingress")
			return ctrl.Result{}, err
		}
	}

	log.Info(fmt.Sprintf("We did find the ingress %+v \n", ingress))
	log.Info(fmt.Sprintf("TODO: Create the api resources needed for this %s", edgeName))
	return ctrl.Result{}, nil
}

func (t *IngressReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&netv1.Ingress{}).
		Complete(t)
}
