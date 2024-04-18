/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/igordcard/proxius/api/v1alpha1"
	proxyv1alpha1 "github.com/igordcard/proxius/api/v1alpha1"
)

// ProxyDefReconciler reconciles a ProxyDef object
type ProxyDefReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=proxy.igordc.com,resources=proxydefs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=proxy.igordc.com,resources=proxydefs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=proxy.igordc.com,resources=proxydefs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ProxyDef object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *ProxyDefReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	proxydef := &v1alpha1.ProxyDef{}
	err := r.Get(ctx, req.NamespacedName, proxydef)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// If the custom resource is not found then, it usually means that it was deleted or not created
			// In this way, we will stop the reconciliation
			log.Info("proxydef resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get proxydef")
		return ctrl.Result{}, err
	}

	// Let's just set the status as Unknown when no status are available
	if proxydef.Status.Conditions == nil || len(proxydef.Status.Conditions) == 0 {
		meta.SetStatusCondition(&proxydef.Status.Conditions, metav1.Condition{Type: typeSyncingProxyDef, Status: metav1.ConditionUnknown, Reason: "Reconciling", Message: "Starting reconciliation"})
		if err = r.Status().Update(ctx, proxydef); err != nil {
			log.Error(err, "Failed to update proxydef status")
			return ctrl.Result{}, err
		}

		// Let's re-fetch the proxydef Custom Resource after update the status
		// so that we have the latest state of the resource on the cluster and we will avoid
		// raising the issue "the object has been modified, please apply your changes to
		// the latest version and try again" which would re-trigger the reconciliation
		// if we try to update it again in the following operations
		if err := r.Get(ctx, req.NamespacedName, proxydef); err != nil {
			log.Error(err, "Failed to re-fetch proxydef")
			return ctrl.Result{}, err
		}
	}

	// The following are a few possible return options for a Reconciler:
	// With the error:
	// 		return ctrl.Result{}, err
	// Without an error:
	// 		return ctrl.Result{Requeue: true}, nil
	// Therefore, to stop the Reconcile, use:
	// 		return ctrl.Result{}, nil
	// Reconcile again after X time:
	//  	return ctrl.Result{RequeueAfter: 5 * time.Minute}, nil

	return ctrl.Result{}, nil
}

// Definitions to manage status conditions
const (
	// typeReadyProxyDef represents that the ProxyDef has already generated the respective ConfigMap
	typeReadyProxyDef = "Ready"
	// typeSyncingProxyDef represents that the ProxyDef is in the progress of generating a ConfigMap
	typeSyncingProxyDef = "Syncing"
	// typeDegradedProxyDef represents that a failure has occurred preventing a ConfigMap from being generated
	typeDegradedProxyDef = "Degraded"
)

// SetupWithManager sets up the controller with the Manager.
func (r *ProxyDefReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&proxyv1alpha1.ProxyDef{}).
		Complete(r)
}
