/*
Copyright 2022.

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

package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	frwildav1 "github.com/philippart-s/voxxed-days-go-operator/api/v1"
)

// NginxOperatorReconciler reconciles a NginxOperator object
type NginxOperatorReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=fr.wilda,resources=nginxoperators,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=fr.wilda,resources=nginxoperators/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=fr.wilda,resources=nginxoperators/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the NginxOperator object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *NginxOperatorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	const helloWorldFinalizer = "fr.wilda/finalizer"

	helloWorld := &frwildav1.NginxOperator{}
	err := r.Get(ctx, req.NamespacedName, helloWorld)
	if err != nil {
		if errors.IsNotFound(err) {
			// CR deleted, nothing to do
			log.Info("No CR found, nothing to do 🧐.")
		} else {
			// Error reading the object - requeue the request.
			log.Error(err, "Failed to get CR NginxOperator")
			return ctrl.Result{}, err
		}
	} else {
		// Add finalizer for this CR
		if !controllerutil.ContainsFinalizer(helloWorld, helloWorldFinalizer) {
			controllerutil.AddFinalizer(helloWorld, helloWorldFinalizer)
			err = r.Update(ctx, helloWorld)
			if err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil
		}

		if helloWorld.GetDeletionTimestamp() != nil {
			// CR marked for deletion ➡️ Goodbye
			log.Info("Goodbye " + helloWorld.Spec.Name + " 😢")
			controllerutil.RemoveFinalizer(helloWorld, helloWorldFinalizer)
			err := r.Update(ctx, helloWorld)
			if err != nil {
				log.Info("Error during deletion")
				return ctrl.Result{}, err
			}
		} else {
			// CR created / updated ➡️ Hello
			log.Info("Hello " + helloWorld.Spec.Name + " 🎉🎉 !!")
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NginxOperatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&frwildav1.NginxOperator{}).
		Complete(r)
}
