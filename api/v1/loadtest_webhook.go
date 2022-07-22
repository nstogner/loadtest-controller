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

package v1

import (
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var loadtestlog = logf.Log.WithName("loadtest-resource")

func (r *LoadTest) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-tests-tbd-com-v1-loadtest,mutating=true,failurePolicy=fail,sideEffects=None,groups=tests.tbd.com,resources=loadtests,verbs=create;update,versions=v1,name=mloadtest.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &LoadTest{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *LoadTest) Default() {
	loadtestlog.Info("default", "name", r.Name)

	//if r.Spec.Method == "" {
	//	r.Spec.Method = "GET"
	//}
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-tests-tbd-com-v1-loadtest,mutating=false,failurePolicy=fail,sideEffects=None,groups=tests.tbd.com,resources=loadtests,verbs=delete;create;update,versions=v1,name=vloadtest.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &LoadTest{}

func (r *LoadTest) validate() error {
	var allErrs field.ErrorList

	if r.Spec.Duration.Duration > time.Hour {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec").Child("duration"), r.Spec.Duration.Duration, "must be no more than 1 hour"))
	}

	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		schema.GroupKind{Group: "tests.tbd.com", Kind: "LoadTest"},
		r.Name, allErrs)
}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *LoadTest) ValidateCreate() error {
	loadtestlog.Info("validate create", "name", r.Name)

	return r.validate()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *LoadTest) ValidateUpdate(old runtime.Object) error {
	loadtestlog.Info("validate update", "name", r.Name)

	return r.validate()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *LoadTest) ValidateDelete() error {
	loadtestlog.Info("validate delete", "name", r.Name)

	return nil
}
