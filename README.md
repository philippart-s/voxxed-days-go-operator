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