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
	"fmt"

	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	cachev1alpha1 "github.com/example/memcached-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
)

// MemcachedReconciler reconciles a Memcached object
type MemcachedReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const EMPTY = "EMPTY"

var Mgr_client client.Client
var ApiReader client.Reader

//+kubebuilder:rbac:groups=cache.example.com,resources=memcacheds,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cache.example.com,resources=memcacheds/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cache.example.com,resources=memcacheds/finalizers,verbs=update

//+kubebuilder:rbac:groups=cache.example.com,namespace=test,resources=memcacheds,verbs=get;list;
//+kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;

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

	if value != "" {
		mem.Status.State = mem.Spec.Foo
	} else {
		mem.Status.State = EMPTY
	}

	mem1 := &cachev1alpha1.Memcached{}
	namespacedName := types.NamespacedName{
		Namespace: "test",
		Name:      "t1",
	}
	err = r.Get(ctx, namespacedName, mem1)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("memcached CR was not found: t1")
		}
		// Error reading the object - requeue the req.
	} else {
		log.Info("found mem1: " + mem1.Spec.Foo)
	}

	//方法1. 去拿别的ns的pod
	pod1 := &corev1.Pod{}
	namespacedName = types.NamespacedName{
		Namespace: "test",
		Name:      "alpine",
	}
	err = r.Get(ctx, namespacedName, pod1)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("pod alpine in test ns was not found ========1")
		}
		// Error reading the object - requeue the req.
	} else {
		log.Info("pod alpine in test ns is found !!!  ========1")
	}

	//方法2. 去拿别的ns的pod
	pod2 := &corev1.Pod{}
	namespacedName = types.NamespacedName{
		Namespace: "test",
		Name:      "alpine",
	}
	err = Mgr_client.Get(ctx, namespacedName, pod2)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("pod alpine in test ns was not found ========2")
		}
		// Error reading the object - requeue the req.
	} else {
		log.Info("pod alpine in test ns is found !!! ========2")
	}

	//方法3. 去拿别的ns的pod
	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Println(err)
	}
	c, err := client.New(cfg, client.Options{})
	if err != nil {
		fmt.Println(err)
	}
	err = c.Get(ctx, namespacedName, pod2)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("pod alpine in test ns was not found again ========3")
		}
		// Error reading the object - requeue the req.
	} else {
		log.Info("pod alpine in test ns is found!!! ========3 ")
	}

	mem2 := &cachev1alpha1.Memcached{}
	namespacedName2 := types.NamespacedName{
		Namespace: "test",
		Name:      "t1",
	}
	err = c.Get(ctx, namespacedName2, mem2)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("cr t1 in test ns was not found  ========3-a")
		} else {
			fmt.Println(err, "=============== 3-a")
		}

		// Error reading the object - requeue the req.
	} else {
		log.Info("cr t1 in test ns is found!!! ========3-a ")
	}

	err = ApiReader.Get(ctx, namespacedName2, mem2)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("cr t1 in test ns was not found  ========3-b")
		} else {
			fmt.Println(err)
		}

		// Error reading the object - requeue the req.
	} else {
		log.Info("cr t1 in test ns is found!!! ========3-b")
	}

	//方法4. 去拿别的ns的pod
	//这个方法只能在k8s里运行，如果要本地运行，需要换个方法拿config.
	configA, err := rest.InClusterConfig()
	if err != nil {
		fmt.Println("error get InClusterConfig config")
		return ctrl.Result{}, nil
	}
	clientset, err := kubernetes.NewForConfig(configA)
	if err != nil {
		fmt.Println("error NewForConfig")
		return ctrl.Result{}, nil
	}

	pod, err := clientset.CoreV1().Pods("test").Get(context.TODO(), "alpine", metav1.GetOptions{})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("pod alpine in test ns is found!!! ========4", pod.Name)
	}

	return ctrl.Result{}, r.Status().Update(ctx, mem)
}

// SetupWithManager sets up the controller with the Manager.
func (r *MemcachedReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cachev1alpha1.Memcached{}).
		Complete(r)
}
