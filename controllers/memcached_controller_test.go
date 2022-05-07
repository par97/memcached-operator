package controllers

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/types"

	cachev1alpha1 "github.com/example/memcached-operator/api/v1alpha1"
)

const timeout = time.Second * 30
const interval = time.Millisecond * 250

func CreateFusionObj(obj client.Object, createdObj client.Object) {
	k8sClient.Create(context.TODO(), obj)
	Eventually(func() bool {
		err := k8sClient.Get(context.TODO(), types.NamespacedName{Name: obj.GetName(), Namespace: obj.GetNamespace()}, createdObj)
		return err == nil
	}, timeout, interval).Should(BeTrue())
}

func DeleteFusionObj(objName string, objNs string, obj client.Object) {
	Eventually(func() bool {
		err := k8sClient.Get(context.TODO(), types.NamespacedName{Name: objName, Namespace: objNs}, obj)
		if err == nil {
			k8sClient.Delete(context.TODO(), obj)
		}
		return err != nil
	}, timeout, interval).Should(BeTrue())
}

func InitMemcached() *cachev1alpha1.Memcached {
	return &cachev1alpha1.Memcached{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "memcached1",
			Namespace: "testabc",
		},
		Spec: cachev1alpha1.MemcachedSpec{
			Foo: "bar1",
		},
	}
}

var _ = Describe("Test MemcachedController", func() {

	//mem := &cachev1alpha1.Memcached{}

	Context("When testing memcached", func() {
		// BeforeEach(func() {
		// 	By("Creating a memcached")
		//..CreateFusionObj(InitMemcached(), mem)
		// })

		// AfterEach(func() {
		// 	By("Deleting the memcached")
		// 	DeleteFusionObj("memcached1", "testabc", mem)
		// })

		Context("When mem created", func() {
			// It("mem is created", func() {

			// 	found := &cachev1alpha1.Memcached{}

			// 	By("Finding the created mem")
			// 	Eventually(func() bool {
			// 		err := k8sClient.Get(context.TODO(), types.NamespacedName{Name: "memcached1", Namespace: "testabc"}, found)
			// 		return err == nil
			// 	}, timeout, interval).Should(BeTrue())

			// 	// By("Querying the created fusion bsl parameters")
			// 	// Expect(found.Status.State).Should(Equal("bar1"))
			// })

			It("should be equal", func() {
				Expect("foo").To(Equal("foo"))
			})

		})
	})
})

func Add(a, b int) int {
	return a + b
}
