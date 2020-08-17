package main

import (
	"encoding/json"
	"strings"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/api/admission/v1beta1"
	"k8s.io/klog"
)

const (
	podsInitContainerPatch string = `[
		 {"op":"add","path":"/spec/initContainers","value":[{"image":"webhook-added-image","name":"webhook-added-init-container","resources":{}}]}
	]`
)

type PatchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}


func mutatePods(ar v1beta1.AdmissionReview, mirror map[string]string) *v1beta1.AdmissionResponse {
	podResource := metav1.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}
	if ar.Request.Resource != podResource {
		klog.Errorf("expect resource to be %s", podResource)
		return nil
	}

	raw := ar.Request.Object.Raw
	pod := corev1.Pod{}
	deserializer := codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(raw, nil, &pod); err != nil {
		klog.Error(err)
		return toAdmissionResponse(err)
	}
	reviewResponse := v1beta1.AdmissionResponse{}
	reviewResponse.Allowed = true

	length := len(pod.Spec.Containers)
	patchs := make([]PatchOperation, 0, length)
	for _, container := range pod.Spec.Containers {
		for k, v := range mirror {
			if strings.HasPrefix(container.Image, k) {
				newContainer :=  container
				newContainer.Image = strings.Replace(container.Image, k, v, -1)
				patchs = append(patchs, PatchOperation{
					Op:    "replace",
					Path:  "/spec/containers/-",//如果是数组类型，非第一个需要加上“/-”
					Value: newContainer,
				})
			}
		}
	}

	if len(patchs) > 0 {
		patchBytes ,_ := json.Marshal(patchs)

		// reviewResponse.Patch = []byte(podsInitContainerPatch)
		reviewResponse.Patch = patchBytes
		pt := v1beta1.PatchTypeJSONPatch
		reviewResponse.PatchType = &pt
	}

	return &reviewResponse
}
