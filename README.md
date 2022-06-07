# voxxed-days-go-operator
Source code for the Go operator talk at Voxxed Days Lux 2022

**⚠️ Au moment de l'écriture de ce tuto, la dernière version de Go est 1.18.1, hors le SDK (1.20.1) n'est pas compatible avec cette version il faut utiliser une version 1.17.x maximum ⚠️**

## 🎉 Init project
 - la branche `01-init-project` contient le résultat de cette étape
 - [installer / mettre](https://sdk.operatorframework.io/docs/installation/) à jour la dernière version du [Operator SDK](https://sdk.operatorframework.io/) (v1.20.1 au moment de l'écriture du readme)
 - créer le répertoire `voxxed-days-go-operator`
 - dans le répertoire `voxxed-days-go-operator`, scaffolding du projet : `operator-sdk init --domain fr.wilda --repo github.com/philippart-s/voxxed-days-go-operator`
 - A ce stage une arborescence complète a été générée, notamment la partie configuration dans `config` et un `Makefile` permettant le lancement des différentes commandes de build
 - vérification que cela compile : `go build`

## 📄 CRD generation
 - la branche `02-crd-generation` contient le résultat de cette étape
 - création de l'API : `operator-sdk create api --version v1 --kind NginxOperator --resource --controller`
 - de nouveau, de nombreux fichiers de générés, notamment le controller `./controllers/nginxoperator_controller.go`
```go
package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

    // ⚠️ Change the github repository name with your own repository name ⚠️
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
	_ = log.FromContext(ctx)

	// TODO(user): your logic here

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NginxOperatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&frwildav1.NginxOperator{}).
		Complete(r)
}
```
  - ensuite on génère la CRD avec la commande `make manifests`:

```yaml
---
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  annotations:
    controller-gen.kubebuilder.io/version: v0.8.0
  creationTimestamp: null
  name: nginxoperators.fr.wilda
spec:
  group: fr.wilda
  names:
    kind: NginxOperator
    listKind: NginxOperatorList
    plural: nginxoperators
    singular: nginxoperator
  scope: Namespaced
  versions:
  - name: v1
    schema:
      openAPIV3Schema:
        description: NginxOperator is the Schema for the nginxoperators API
        properties:
          apiVersion:
            description: 'APIVersion defines the versioned schema of this representation
              of an object. Servers should convert recognized schemas to the latest
              internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources'
            type: string
          kind:
            description: 'Kind is a string value representing the REST resource this
              object represents. Servers may infer this from the endpoint the client
              submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds'
            type: string
          metadata:
            type: object
          spec:
            description: NginxOperatorSpec defines the desired state of NginxOperator
            properties:
              foo:
                description: Foo is an example field of NginxOperator. Edit nginxoperator_types.go
                  to remove/update
                type: string
            type: object
          status:
            description: NginxOperatorStatus defines the observed state of NginxOperator
            type: object
        type: object
    served: true
    storage: true
    subresources:
      status: {}
status:
  acceptedNames:
    kind: ""
    plural: ""
  conditions: []
  storedVersions: []
```
  - puis on peut l'appliquer avec la commande `make install`
  - et vérfier qu'elle a été créée : `kubectl get crd nginxoperators.fr.wilda`
```bash
kubectl get crd nginxoperators.fr.wilda
NAME                      CREATED AT
nginxoperators.fr.wilda   2022-06-07T09:42:46Z
```

## 👋  Hello World
 - la branche `03-hello-world` contient le résultat de cette étape
 - ajouter un champ `name` dans `api/v1/nginxoperator_types.go`:
 ℹ️ les attributs sont définis en utilisant les _JSON tags_ : https://pkg.go.dev/encoding/json#Marshal ℹ️
```go
package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NginxOperatorSpec defines the desired state of NginxOperator
type NginxOperatorSpec struct {
	Name string `json:"name,omitempty"`
}

// NginxOperatorStatus defines the observed state of NginxOperator
type NginxOperatorStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// NginxOperator is the Schema for the nginxoperators API
type NginxOperator struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NginxOperatorSpec   `json:"spec,omitempty"`
	Status NginxOperatorStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// NginxOperatorList contains a list of NginxOperator
type NginxOperatorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NginxOperator `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NginxOperator{}, &NginxOperatorList{})
}
```
  - lancer la commande `make manifests` pour générer le manifest de la CRD
  - mettre à jour la CRD : `make install`
  - vérifier que la CRD a bien été mise à jour:
```bash
$ kubectl get crds nginxoperators.fr.wilda -o json | jq '.spec.versions[0].schema.openAPIV3Schema.properties.spec'

{
  "description": "NginxOperatorSpec defines the desired state of NginxOperator",
  "properties": {
    "name": {
      "description": "Name to say hello",
      "type": "string"
    }
  },
  "type": "object"
}
```
 - modifier le reconciler `controllers/nginxoperator_controller.go`:
```go
package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	frwildav1 "github.com/philippart-s/go-operator-template/api/v1"
)

// NginxOperatorReconciler reconciles a NginxOperator object
type NginxOperatorReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=fr.wilda,resources=nginxoperators,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=fr.wilda,resources=nginxoperators/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=fr.wilda,resources=nginxoperators/finalizers,verbs=update

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
```
  - vérifier que to compile `go build`
  - lancer l'opérateur sur la machine de dev : `make install run`:
```bash
$ make install run

/Users/stef/Talks/operators-for-all-dev/voxxed-days-2022/voxxed-days-go-operator/bin/controller-gen rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases
/Users/stef/Talks/operators-for-all-dev/voxxed-days-2022/voxxed-days-go-operator/bin/kustomize build config/crd | kubectl apply -f -
customresourcedefinition.apiextensions.k8s.io/nginxoperators.fr.wilda configured
/Users/stef/Talks/operators-for-all-dev/voxxed-days-2022/voxxed-days-go-operator/bin/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
go fmt ./...
go vet ./...
go run ./main.go
1.654608708040212e+09   INFO    controller-runtime.metrics      Metrics server is starting to listen    {"addr": ":8080"}
1.6546087080409338e+09  INFO    setup   starting manager
1.6546087080413399e+09  INFO    Starting server {"path": "/metrics", "kind": "metrics", "addr": "[::]:8080"}
1.6546087080414422e+09  INFO    Starting server {"kind": "health probe", "addr": "[::]:8081"}
1.654608708041564e+09   INFO    controller.nginxoperator        Starting EventSource    {"reconciler group": "fr.wilda", "reconciler kind": "NginxOperator", "source": "kind source: *v1.NginxOperator"}
1.654608708041604e+09   INFO    controller.nginxoperator        Starting Controller     {"reconciler group": "fr.wilda", "reconciler kind": "NginxOperator"}
1.654608708143117e+09   INFO    controller.nginxoperator        Starting workers        {"reconciler group": "fr.wilda", "reconciler kind": "NginxOperator", "worker count": 1}
```
  - créer le namespace `test-helloworld-operator`: `kubectl create ns test-helloworld-operator`
  - mettre à jour la CR `config/samples/_v1_nginxoperator.yaml` pour tester:
```yaml
apiVersion: fr.wilda/v1
kind: NginxOperator
metadata:
  name: nginxoperator-sample
spec:
  name: Voxxed Days Lux
```
  - créer la CR dans Kubernetes : `kubectl apply -f ./config/samples/_v1_nginxoperator.yaml -n test-helloworld-operator`
  - la sortie de l'opérateur devrait afficher le message `Hello Voxxed Days Lux 🎉🎉 !!`
  - supprimer la CR : `kubectl delete nginxoperators.fr.wilda/nginxoperator-sample -n test-helloworld-operator`
  - la sortie de l'opérateur devrait afficher le message `Goodbye Voxxed Days Lux 😢`

## 🤖 Nginx operator
 - la branche `04-nginx-operator` contient le résultat de cette étape
 - modifier le fichier `api/v1/nginxoperator_types.go`:
```go
package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NginxOperatorSpec defines the desired state of NginxOperator
type NginxOperatorSpec struct {
	// Number of replicas for the Nginx Pods
	ReplicaCount int32 `json:"replicaCount"`
	// Exposed port for the Nginx server
	Port int32 `json:"port"`
}

// Unchanged code
// ...
```
 - générer la CRD modifiée : `make manifests`
 - deployer la CRD dans Kubernetes : `make install`
 - vérifier que la CRD a bien été mise à jour:
```bash
$ kubectl get crds nginxoperators.fr.wilda -o json | jq '.spec.versions[0].schema.openAPIV3Schema.properties.spec'

{
  "description": "NginxOperatorSpec defines the desired state of NginxOperator",
  "properties": {
    "port": {
      "description": "Exposed port for the Nginx server",
      "format": "int32",
      "type": "integer"
    },
    "replicaCount": {
      "description": "Number of replicas for the Nginx Pods",
      "format": "int32",
      "type": "integer"
    }
  },
  "required": [
    "port",
    "replicaCount"
  ],
  "type": "object"
}
```
 - modifier le reconciler en ajoutant les fonctions permettant la création d'un Deployment et d'un Service dans `controllers/nginxoperator_controller.go`:
```go
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

	return service
}
```
 - modifier le reconciler pour la création du Pod et de son Service dans `controllers/nginxoperator_controller.go`:
```go
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

	frwildav1 "github.com/philippart-s/go-operator-template/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NginxOperatorReconciler reconciles a NginxOperator object
type NginxOperatorReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=fr.wilda,resources=nginxoperators,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=fr.wilda,resources=nginxoperators/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=fr.wilda,resources=nginxoperators/finalizers,verbs=update

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

// ... unmodified code

``` 
 - vérifier que ça compile : `go build`
 - créer le namespace `test-nginx-operator`: `kubectl create ns test-nginx-operator`
 - créer la CR: `config/samples/_v2_nginxoperator.yaml`:
```yaml
apiVersion: fr.wilda/v1
kind: NginxOperator
metadata:
  name: nginxoperator-sample
spec:
  replicaCount: 1
  port: 30080
```
 - lancer l'opérateur en mode dev : `make install run`
 - puis créer la CR sur Kubernetes: `kubectl apply -f ./config/samples/_v2_nginxoperator.yaml -n test-nginx-operator`
 - l'opérateur devrait créer le pod Nginx et son service:
```bash
INFO    controller.nginxoperator        ⚡️ Event !!! ⚡️ {"reconciler group": "fr.wilda", "reconciler kind": "NginxOperator", "name": "nginxoperator-sample", "namespace": "test-nginx-operator"}
INFO    controller.nginxoperator        ✨ Creating a new Deployment    {"reconciler group": "fr.wilda", "reconciler kind": "NginxOperator", "name": "nginxoperator-sample", "namespace": "test-nginx-operator", "Deployment.Namespace": "test-nginx-operator", "Deployment.Name": "nginxoperator-sample"}      
```
      Dans Kubernetes:
```bash
$ kubectl get pod,svc,nginxoperator  -n test-nginx-operator
NAME                                        READY   STATUS    RESTARTS   AGE
pod/nginxoperator-sample-58c4f478ff-phz7t   1/1     Running   0          17s

NAME                           TYPE       CLUSTER-IP    EXTERNAL-IP   PORT(S)        AGE
service/nginxoperator-sample   NodePort   10.3.173.89   <none>        80:30080/TCP   17s

NAME                                          AGE
nginxoperator.fr.wilda/nginxoperator-sample   17s
```
 - tester dans un navigateur ou par un curl l'accès à `http://<node external ip>:30080`, pour récupérer l'IP externe du node : `kubectl cluster-info`

## ✏️ Update and delete CR
 - la branche `05-update-cr` contient le résultat de cette étape
 - changer le port et le nombre de replicas dans la CR `config/samples/_v2_nginxoperator.yaml`:
```yaml
apiVersion: fr.wilda/v1
kind: NginxOperator
metadata:
  name: nginxoperator-sample
spec:
  replicaCount: 2
  port: 30081
```
 - appliquer la CR: `kubectl apply -f ./config/samples/_v2_nginxoperator.yaml -n test-nginx-operator`
 - vérifier que le nombre de pods et le port ont bien changés:
```bash
$ kubectl get pod,svc  -n test-nginx-operator
NAME                                        READY   STATUS    RESTARTS   AGE
pod/nginxoperator-sample-58c4f478ff-l2gwz   1/1     Running   0          10s
pod/nginxoperator-sample-58c4f478ff-phz7t   1/1     Running   0          9m41s

NAME                           TYPE       CLUSTER-IP    EXTERNAL-IP   PORT(S)        AGE
service/nginxoperator-sample   NodePort   10.3.173.89   <none>        80:30081/TCP   9m41s
```
 - tester dans un navigateur ou par un curl l'accès à `http://<node external ip>:30081`
 - supprimer la CR : `kubectl delete nginxoperators.fr.wilda/nginxoperator-sample -n test-nginx-operator`
 - vérifier que rien a été supprimé:
```bash
$ kubectl get pod,svc  -n test-nginx-operator
NAME                                        READY   STATUS    RESTARTS   AGE
pod/nginxoperator-sample-58c4f478ff-l2gwz   1/1     Running   0          2m23s
pod/nginxoperator-sample-58c4f478ff-phz7t   1/1     Running   0          11m

NAME                           TYPE       CLUSTER-IP    EXTERNAL-IP   PORT(S)        AGE
service/nginxoperator-sample   NodePort   10.3.173.89   <none>        80:30081/TCP   11m
```
 - supprimer le Deployment : `kubectl delete deployment nginxoperator-sample -n test-nginx-operator`
 - supprimer le Service : `kubectl delete service nginxoperator-sample  -n test-nginx-operator`

## 👀 Watch CR deletion
 - la branche `06-watch-deletion` contient le résultat de cette étape
 - modifier le reconciler `controllers/nginxoperator_controller.go` pour que le Service et le Deployment soient gérés en suppression:
```go
// unmodified code ...

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

// unmodified code ...
```
 - lancer l'opérateur en mode dev : `make install run`
 - appliquer la CR: `kubectl apply -f ./config/samples/_v2_nginxoperator.yaml -n test-nginx-operator`
 - vérifier que tout est créé:
```bash
$ kubectl get pod,svc  -n test-nginx-operator
NAME                                        READY   STATUS    RESTARTS   AGE
pod/nginxoperator-sample-58c4f478ff-cqm6f   1/1     Running   0          13s
pod/nginxoperator-sample-58c4f478ff-xvs4h   1/1     Running   0          13s

NAME                           TYPE       CLUSTER-IP    EXTERNAL-IP   PORT(S)        AGE
service/nginxoperator-sample   NodePort   10.3.165.48   <none>        80:30081/TCP   13s
```
 - supprimer la CR : `kubectl delete nginxoperators.fr.wilda/nginxoperator-sample -n test-nginx-operator`
 - vérifier que tout est supprimé:
```bash
$ kubectl get pod,svc  -n test-nginx-operator
No resources found in test-nginx-operator namespace.
```

## 👀 Watch service deletion
 - la branche `07-watch-service-deletion` contient le résultat de cette étape
 - utiliser la CR de tests `config/sample/_v2_nginxoperator.yaml` pour créer tous les éléments : `kubectl apply -f ./config/samples/_v2_nginxoperator.yaml -n test-nginx-operator`
 - supprimer le service : `kubectl delete svc/nginxoperator-sample -n test-nginx-operator`
 - constater qu'il n'est pas recréé: `kubectl get svc  -n test-nginx-operator`
```bash
$ kubectl get svc  -n test-nginx-operator

No resources found in test-nginx-operator namespace.
```
 - supprimer la CR : `kubectl delete nginxoperators.fr.wilda/nginxoperator-sample -n test-nginx-operator`
 - modifier le reconciler `controllers/nginxoperator_controller.go` pour qu'il surveille le service:
```go
// unmodified code ...
// SetupWithManager sets up the controller with the Manager.
func (r *NginxOperatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&frwildav1.NginxOperator{}).
		// Enable service Watching
		Owns(&corev1.Service{}).
		Complete(r)
}
```
 - lancer l'opérateur en mode dev : `make install run`
 - appliquer la CR de tests `config/sample/_v2_nginxoperator.yaml` pour créer tous les éléments : `kubectl apply -f ./config/samples/_v2_nginxoperator.yaml -n test-nginx-operator`
 - supprimer le service : `kubectl delete svc/nginxoperator-sample -n test-nginx-operator`
 - l'opérateur le recrée et il devrait réapparaître : 
```bash
INFO    controller.nginxoperator        ⚡️ Event !!! ⚡️ {"reconciler group": "fr.wilda", "reconciler kind": "NginxOperator", "name": "nginxoperator-sample", "namespace": "test-nginx-operator"}
INFO    controller.nginxoperator        ✨ Creating a new Service       {"reconciler group": "fr.wilda", "reconciler kind": "NginxOperator", "name": "nginxoperator-sample", "namespace": "test-nginx-operator", "Service.Namespace": "test-nginx-operator", "Service.Name": "nginxoperator-sample"}
```
```bash
kubectl get pod,svc  -n test-nginx-operator
NAME                                        READY   STATUS    RESTARTS   AGE
pod/nginxoperator-sample-58c4f478ff-ctd8b   1/1     Running   0          54s
pod/nginxoperator-sample-58c4f478ff-rtgqg   1/1     Running   0          54s

NAME                           TYPE       CLUSTER-IP     EXTERNAL-IP   PORT(S)        AGE
service/nginxoperator-sample   NodePort   10.3.245.149   <none>        80:30081/TCP   26s
```
 - supprimer la CR : `kubectl delete nginxoperators.fr.wilda/nginxoperator-sample -n test-nginx-operator`
