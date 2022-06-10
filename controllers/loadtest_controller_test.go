package controllers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"time"

	v1 "github.com/nstogner/loadtest-controller/api/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LoadTest controller", func() {
	// Define utility constants for object names and testing timeouts/durations and intervals.
	const (
		timeout  = time.Second * 10
		duration = time.Second * 10
		interval = time.Millisecond * 250
	)

	Context("When creating a LoadTest", func() {
		It("Should send requests to a server", func() {
			By("By starting a server to be load tested")

			var serverCalled bool
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				serverCalled = true
			}))
			defer ts.Close()

			By("By create a new LoadTest object")
			lt0 := &v1.LoadTest{
				TypeMeta: metav1.TypeMeta{
					Kind:       "LoadTest",
					APIVersion: "tests.tbd.com/v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "loadtest-test",
					Namespace: "default",
				},
				Spec: v1.LoadTestSpec{
					Method:  "GET",
					Address: ts.URL,
					Duration: metav1.Duration{
						Duration: time.Second / 10,
					},
				},
			}
			Expect(k8sClient.Create(ctx, lt0)).Should(Succeed())

			By("By checking that the server was called")
			Eventually(func() bool {
				return serverCalled
			}, timeout, interval).Should(BeTrue())

			By("By checking that the LoadTest status was updated")
			lt1 := &v1.LoadTest{}
			Eventually(func() error {
				err := k8sClient.Get(ctx, types.NamespacedName{Namespace: lt0.Namespace, Name: lt0.Name}, lt1)
				if err != nil {
					return err
				}

				if !lt1.Status.Completed {
					return errors.New("not completed")
				}

				return nil
			}, timeout, interval).Should(Succeed())

			Expect(lt1.Status.RequestCount).Should(BeNumerically(">", 0))
			Expect(lt1.Status.AverageLatency.Nanoseconds()).Should(BeNumerically(">", 0))

			// Yeah, but did it realllly work???
			//json.NewEncoder(os.Stdout).Encode(lt1)
		})
	})
})
