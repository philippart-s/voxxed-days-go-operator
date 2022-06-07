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
