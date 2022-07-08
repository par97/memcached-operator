/*
Copyright 2021.

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

	cachev1alpha1 "github.com/example/memcached-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	//mcfgv1 "github.com/openshift/machine-config-operator/pkg/apis/machineconfiguration.openshift.io/v1"
	mcfgv1 "github.com/openshift/machine-config-operator/pkg/apis/machineconfiguration.openshift.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
)

// MemcachedReconciler reconciles a Memcached object
type MemcachedReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const EMPTY = "EMPTY"

//+kubebuilder:rbac:groups=cache.example.com,resources=memcacheds,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cache.example.com,resources=memcacheds/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cache.example.com,resources=memcacheds/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=nodes,verbs=get;list;patch;watch;update
//+kubebuilder:rbac:groups=machineconfiguration.openshift.io,resources=machineconfigs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=machineconfiguration.openshift.io,resources=containerruntimeconfigs,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Memcached object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *MemcachedReconciler) Reconcile(ctx context.Context, req ctrl.Request) (result ctrl.Result, err error) {
	//_ = log.FromContext(ctx)

	// your logic here
	log := ctrllog.Log.WithValues("request", req.NamespacedName)
	log.Info("enter into MemcachedReconciler")

	mem := &cachev1alpha1.Memcached{}
	err = r.Get(ctx, req.NamespacedName, mem)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("memcached CR was not found")
			err = nil
			return
		}
		// Error reading the object - requeue the req.
		return
	}

	value := mem.Spec.Foo
	log.Info("memcached CR spec Foo: " + value)

	nodeList := &corev1.NodeList{}

	var WorkerRoleLabel = "node-role.kubernetes.io/worker"
	//only query nodes with worker role, but this node could also be master/infra
	listOpts := []client.ListOption{
		client.HasLabels([]string{WorkerRoleLabel}),
	}

	if err := r.List(context.TODO(), nodeList, listOpts...); err != nil {
		return ctrl.Result{}, err
	}

	for _, node := range nodeList.Items {
		node.Labels["t1"] = value
		err := r.Patch(ctx, &node, client.Merge)
		if err != nil {
			return ctrl.Result{}, err
		}
		log.Info("node " + node.Name + " label t1=" + value)
	}

	t1 := &mcfgv1.MachineConfig{
		ObjectMeta: metav1.ObjectMeta{Name: value},
		Spec: mcfgv1.MachineConfigSpec{
			OSImageURL: "//:dummy0",
		},
	}
	err = r.Create(ctx, t1)
	if err != nil {
		log.Error(err, "fail to create machine config "+value)
	}
	log.Info("machine config " + value + " is created")

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MemcachedReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cachev1alpha1.Memcached{}).
		Complete(r)
}
