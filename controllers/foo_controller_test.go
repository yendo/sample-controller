package controllers

import (
	"context"
	//	"errors"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	foov1 "github.com/yendo/sample-controller/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Foo controller", func() {

	ctx := context.Background()
	var stopFunc func()

	BeforeEach(func() {
		err := k8sClient.DeleteAllOf(ctx, &foov1.Foo{}, client.InNamespace("default"))
		Expect(err).NotTo(HaveOccurred())
		err = k8sClient.DeleteAllOf(ctx, &appsv1.Deployment{}, client.InNamespace("default"))
		Expect(err).NotTo(HaveOccurred())
		time.Sleep(100 * time.Millisecond)

		mgr, err := ctrl.NewManager(cfg, ctrl.Options{
			Scheme: scheme.Scheme,
		})
		Expect(err).ToNot(HaveOccurred())

		reconciler := FooReconciler{
			Client: k8sClient,
			Scheme: scheme.Scheme,
		}
		err = reconciler.SetupWithManager(mgr)
		Expect(err).NotTo(HaveOccurred())

		ctx, cancel := context.WithCancel(ctx)
		stopFunc = cancel
		go func() {
			err := mgr.Start(ctx)
			if err != nil {
				panic(err)
			}
		}()
		time.Sleep(100 * time.Millisecond)
	})

	AfterEach(func() {
		stopFunc()
		time.Sleep(100 * time.Millisecond)
	})

	It("should create Deployment", func() {
		foo := newFoo()
		err := k8sClient.Create(ctx, foo)
		Expect(err).NotTo(HaveOccurred())

		dep := appsv1.Deployment{}
		Eventually(func() error {
			return k8sClient.Get(ctx, client.ObjectKey{Namespace: "default", Name: "example-foo"}, &dep)
		}).Should(Succeed())
		Expect(dep.Spec.Replicas).Should(Equal(pointer.Int32Ptr(3)))
		Expect(dep.Spec.Template.Spec.Containers[0].Image).Should(Equal("nginx:latest"))
	})
})

func newFoo() *foov1.Foo {
	return &foov1.Foo{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "foo-sample",
			Namespace: "default",
		},
		Spec: foov1.FooSpec{
			DeploymentName: "example-foo",
			Replicas:       pointer.Int32Ptr(3),
		},
	}
}
