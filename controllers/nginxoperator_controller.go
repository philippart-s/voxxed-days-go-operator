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
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete

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
	log.Info("⚡️ Event !!! ⚡️")

	nginxCR := &frwildav1.NginxOperator{}
	existingNginxDeployment := &appsv1.Deployment{}
	existingService := &corev1.Service{}

	err := r.Get(ctx, req.NamespacedName, nginxCR)
	if err == nil {
		// Check if the deployment already exists, if not: create a new one.
		err = r.Get(ctx, types.NamespacedName{Name: nginxCR.Name, Namespace: nginxCR.Namespace}, existingNginxDeployment)
		if err != nil && errors.IsNotFound(err) {
			// Define a new deployment
			newNginxDeployment := r.createDeployment(nginxCR)
			log.Info("✨ Creating a new Deployment", "Deployment.Namespace", newNginxDeployment.Namespace, "Deployment.Name", newNginxDeployment.Name)

			err = r.Create(ctx, newNginxDeployment)
			if err != nil {
				log.Error(err, "❌ Failed to create new Deployment", "Deployment.Namespace", newNginxDeployment.Namespace, "Deployment.Name", newNginxDeployment.Name)
				return ctrl.Result{}, err
			}
		} else if err == nil {
			// Deployment exists, check if the Deployment must be updated
			var replicaCount int32 = nginxCR.Spec.ReplicaCount
			if *existingNginxDeployment.Spec.Replicas != replicaCount {
				log.Info("🔁 Number of replicas changes, update the deployment! 🔁")
				existingNginxDeployment.Spec.Replicas = &replicaCount
				err = r.Update(ctx, existingNginxDeployment)
				if err != nil {
					log.Error(err, "❌ Failed to update Deployment", "Deployment.Namespace", existingNginxDeployment.Namespace, "Deployment.Name", existingNginxDeployment.Name)
					return ctrl.Result{}, err
				}
			}
		}

		// Check if the service already exists, if not: create a new one
		err = r.Get(ctx, types.NamespacedName{Name: nginxCR.Name, Namespace: nginxCR.Namespace}, existingService)
		if err != nil && errors.IsNotFound(err) {
			// Create the Service
			newService := r.createService(nginxCR)
			log.Info("✨ Creating a new Service", "Service.Namespace", newService.Namespace, "Service.Name", newService.Name)
			err = r.Create(ctx, newService)
			if err != nil {
				log.Error(err, "❌ Failed to create new Service", "Service.Namespace", newService.Namespace, "Service.Name", newService.Name)
				return ctrl.Result{}, err
			}
		} else if err == nil {
			// Service exists, check if the port have to be updated.
			var port int32 = nginxCR.Spec.Port
			if existingService.Spec.Ports[0].NodePort != port {
				log.Info("🔁 Port number changes, update the service! 🔁")
				existingService.Spec.Ports[0].NodePort = port
				err = r.Update(ctx, existingService)
				if err != nil {
					log.Error(err, "❌ Failed to update Service", "Service.Namespace", existingService.Namespace, "Service.Name", existingService.Name)
					return ctrl.Result{}, err
				}
			}
		} else if err != nil {
			log.Error(err, "Failed to get Service")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// Create a Deployment for the Nginx server.
func (r *NginxOperatorReconciler) createDeployment(nginxCR *frwildav1.NginxOperator) *appsv1.Deployment {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      nginxCR.Name,
			Namespace: nginxCR.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &nginxCR.Spec.ReplicaCount,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "nginx"},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": "nginx"},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: "ovhplatform/hello:1.0",
						Name:  "nginx",
						Ports: []corev1.ContainerPort{{
							ContainerPort: 80,
							Name:          "http",
							Protocol:      "TCP",
						}},
					}},
				},
			},
		},
	}

	// Set nginxCR instance as the owner and controller
	ctrl.SetControllerReference(nginxCR, deployment, r.Scheme)

	return deployment
}

// Create a Service for the Nginx server.
func (r *NginxOperatorReconciler) createService(nginxCR *frwildav1.NginxOperator) *corev1.Service {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      nginxCR.Name,
			Namespace: nginxCR.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": "nginx",
			},
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					NodePort:   nginxCR.Spec.Port,
					Port:       80,
					TargetPort: intstr.FromInt(80),
				},
			},
			Type: corev1.ServiceTypeNodePort,
		},
	}

	// Set nginxCR instance as the owner and controller
	ctrl.SetControllerReference(nginxCR, service, r.Scheme)

	return service
}

// SetupWithManager sets up the controller with the Manager.
func (r *NginxOperatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&frwildav1.NginxOperator{}).
		// Enable service Watching
		Owns(&corev1.Service{}).
		Complete(r)
}
