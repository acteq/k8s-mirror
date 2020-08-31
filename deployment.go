package main

import (
	"encoding/json"
	"strings"
	"strconv"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/api/admission/v1beta1"
	"k8s.io/klog"
)


func mutateDeployments(ar v1beta1.AdmissionReview, mirror map[string]string) *v1beta1.AdmissionResponse {
	resource := metav1.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	if ar.Request.Resource != resource {
		klog.Errorf("expect resource to be %s", resource)
		return nil
	}

	raw := ar.Request.Object.Raw
	deployment := appsv1.Deployment{}
	deserializer := codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(raw, nil, &deployment); err != nil {
		klog.Error(err)
		return toAdmissionResponse(err)
	}
	reviewResponse := v1beta1.AdmissionResponse{}
	reviewResponse.Allowed = true

	length := len(deployment.Spec.Template.Spec.Containers)
	patchs := make([]PatchOperation, 0, length)
	for i, container := range deployment.Spec.Template.Spec.Containers {
		for k, v := range mirror {
			originImage := container.Image
			if strings.HasPrefix(originImage, k) {
				newImage := strings.Replace(originImage, k, v, -1)
				patchs = append(patchs, PatchOperation{
					Op:    "replace",
					Path:  "/spec/template/spec/containers/" + strconv.Itoa(i) +"/image",
					Value: newImage,
				})
				klog.Infof(newImage)
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

