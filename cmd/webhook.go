/*
Copyright 2024 Igor DC.

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

package main

import (
	"context"
	"fmt"
	"net/http"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

//+kubebuilder:webhook:path=/mutate-v1-pod,mutating=true,failurePolicy=fail,groups="",resources=pods,verbs=create;update,versions=v1,name=mpod.kb.io,admissionReviewVersions=v1,sideEffects=NoneOnDryRun

type PodMutator struct {
	Client  client.Client
	decoder *admission.Decoder
}

func (a *PodMutator) Handle(ctx context.Context, req admission.Request) admission.Response {
	log := logf.FromContext(ctx)
	pod := &corev1.Pod{}
	if err := a.decoder.Decode(req, pod); err != nil {
		log.Info("Failed to decode Pod", "err", err)
		return admission.Errored(http.StatusBadRequest, err)
	}

	if pod.Annotations == nil {
		pod.Annotations = map[string]string{}
	}
	pod.Annotations["example-mutating-admission-webhook"] = "foo"
	log.Info("Annotated Pod")

	return admission.Allowed("")
}

func (a *PodMutator) Default(ctx context.Context, obj runtime.Object) error {
	log := logf.FromContext(ctx)
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return fmt.Errorf("expected a Pod but got a %T", obj)
	}

	if pod.Annotations == nil {
		pod.Annotations = map[string]string{}
	}
	pod.Annotations["example-mutating-admission-webhook"] = "foo"
	log.Info("Annotated Pod")

	return nil
}
