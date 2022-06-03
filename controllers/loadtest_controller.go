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
	"fmt"

	testsv1 "github.com/nstogner/loadtest-controller/api/v1"
	"github.com/nstogner/loadtest-controller/runner"
	"github.com/prometheus/client_golang/prometheus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
	metricRunsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "loadtest_runs_total",
			Help: "Number of load tests ran",
		},
	)
)

func init() {
	// Register custom metrics with the global prometheus registry
	metrics.Registry.MustRegister(metricRunsTotal)
}

// LoadTestReconciler reconciles a LoadTest object
type LoadTestReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	Concurrency int

	record.EventRecorder
}

//+kubebuilder:rbac:groups=tests.tbd.com,resources=loadtests,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=tests.tbd.com,resources=loadtests/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=tests.tbd.com,resources=loadtests/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the LoadTest object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *LoadTestReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	lg := log.FromContext(ctx)

	lg.Info("Reconciling LoadTest")

	var lt testsv1.LoadTest
	if err := r.Get(ctx, req.NamespacedName, &lt); err != nil {
		return ctrl.Result{}, fmt.Errorf("getting LoadTest: %w", err)
	}

	if lt.Status.Completed {
		lg.Info("Already completed, returning")
		return ctrl.Result{}, nil
	}

	lg.Info("Running load test", "duration", lt.Spec.Duration)
	r.Eventf(&lt, corev1.EventTypeNormal, "LoadTestStarted", "Running load test for %s", lt.Spec.Duration.Duration.String())
	metricRunsTotal.Inc()

	out := runner.Run(ctx, runner.Input{
		Method:    lt.Spec.Method,
		URL:       lt.Spec.Address,
		Duration:  lt.Spec.Duration.Duration,
		ReqPerSec: 60,
	})

	lg.Info("Done running load test", "duration", lt.Spec.Duration)

	lt.Status.Completed = true
	lt.Status.RequestCount = out.RequestCount
	lt.Status.AverageLatency = metav1.Duration{Duration: out.AverageLatency()}

	if err := r.Status().Update(ctx, &lt); err != nil {
		return ctrl.Result{}, fmt.Errorf("updating LoadTest: %w", err)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LoadTestReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&testsv1.LoadTest{}).
		WithOptions(controller.Options{
			MaxConcurrentReconciles: r.Concurrency,
		}).
		Complete(r)
}
