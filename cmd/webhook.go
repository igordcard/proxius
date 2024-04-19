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
	"encoding/json"
	"net/http"

	corev1 "k8s.io/api/core/v1"
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

	// TODO: figure out configmap name dynamically from proxydef
	for i := range pod.Spec.Containers {
		pod.Spec.Containers[i].EnvFrom = append(pod.Spec.Containers[i].EnvFrom, corev1.EnvFromSource{
			ConfigMapRef: &corev1.ConfigMapEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: "proxydef-config",
				},
			},
		})
	}

	marshaledPod, err := json.Marshal(pod)
	if err != nil {
		log.Info("Failed to encode Pod", "err", err)
		return admission.Errored(http.StatusConflict, err)
	}

	log.Info("Patching Pod with proxy environment", "err", nil)
	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledPod)
}
