/*

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
// +kubebuilder:docs-gen:collapse=Apache License

/*
Ideally, we should have one `<kind>_controller_test.go` for each controller scaffolded and called in the `suite_test.go`.
So, let's write our example test for the memcached controller (`memcached_controller_test.go.`)
*/

/*
As usual, we start with the necessary imports. We also define some utility variables.
*/
package controllers

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	memcachedv1 "github.com/example/memcached-operator/api/v1alpha1"
)

// +kubebuilder:docs-gen:collapse=Imports

/*
The first step to writing a simple integration test is to actually create an instance of memcached you can run tests against.
Note that to create a memcached, you’ll need to create a stub memcached struct that contains your memcached’s specifications.

Note that when we create a stub memcached, the memcached also needs stubs of its required downstream objects.
Without the stubbed Job template spec and the Pod template spec below, the Kubernetes API will not be able to
create the memcached.
*/
var _ = Describe("memcached controller:", func() {

	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		memcachedName      = "test-memcached"
		memcachedNamespace = "default"
		FooValue           = "bar3"

		timeout  = time.Second * 10
		duration = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("When creating memcached >", func() {

		It("Should update memcached status >", func() {
			By("By creating a new memcached")
			ctx := context.Background()
			memcached := &memcachedv1.Memcached{
				ObjectMeta: metav1.ObjectMeta{
					Name:      memcachedName,
					Namespace: memcachedNamespace,
				},
				Spec: memcachedv1.MemcachedSpec{
					Foo: FooValue,
				},
			}
			Expect(k8sClient.Create(ctx, memcached)).Should(Succeed())

			memcachedLookupKey := types.NamespacedName{Name: memcachedName, Namespace: memcachedNamespace}
			createdmemcached := &memcachedv1.Memcached{}

			// We'll need to retry getting this newly created memcached, given that creation may not immediately happen.
			By("By checking the new memcached")
			Eventually(func() bool {
				err := k8sClient.Get(ctx, memcachedLookupKey, createdmemcached)
				if err != nil {
					return false
				}
				return err == nil && createdmemcached.Status.State == FooValue

			}, timeout, interval).Should(BeTrue())
			Expect(createdmemcached.Status.State).Should(Equal(FooValue))

		})
		It("Should update memcached status EMPTY >", func() {
			memcachedName1 := "test-memcached1"
			By("By creating a new memcached")
			ctx := context.Background()
			memcached := &memcachedv1.Memcached{
				ObjectMeta: metav1.ObjectMeta{
					Name:      memcachedName1,
					Namespace: memcachedNamespace,
				},
				Spec: memcachedv1.MemcachedSpec{},
			}
			Expect(k8sClient.Create(ctx, memcached)).Should(Succeed())

			memcachedLookupKey := types.NamespacedName{Name: memcachedName1, Namespace: memcachedNamespace}
			createdmemcached := &memcachedv1.Memcached{}

			// We'll need to retry getting this newly created memcached, given that creation may not immediately happen.
			By("By checking the new memcached")
			Eventually(func() bool {
				err := k8sClient.Get(ctx, memcachedLookupKey, createdmemcached)
				if err != nil {
					return false
				}
				return err == nil && createdmemcached.Status.State == EMPTY

			}, timeout, interval).Should(BeTrue())
			Expect(createdmemcached.Status.State).Should(Equal(EMPTY))

		})
	})

})

/*
	After writing all this code, you can run `go test ./...` in your `controllers/` directory again to run your new test!
*/
