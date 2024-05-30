package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kmrhemant916/k8s-webhooks/helpers"
	"github.com/wI2L/jsondiff"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
)

const (
	Config = "config/config.yaml"
)

func (app *App) Mutate(w http.ResponseWriter, r *http.Request) {
    admissionReview := &admissionv1.AdmissionReview{}
    err := json.NewDecoder(r.Body).Decode(admissionReview)
    if err != nil {
        helpers.HandleError(w, r, err)
        return
    }
    pod := &corev1.Pod{}
    if err := json.Unmarshal(admissionReview.Request.Object.Raw, pod); err != nil {
        helpers.HandleError(w, r, fmt.Errorf("unmarshal to pod: %v", err))
        return
    }
    originalJSON := admissionReview.Request.Object.Raw
    config, err := helpers.ReadConfig(Config)
    if err != nil {
        panic(err)
    }
    for _, label := range config.TargetLabels {
        if pod.Labels[label.Key] == label.Value {
            toleration := corev1.Toleration{
                Key:      config.Tolerations[0].Key,
                Operator: corev1.TolerationOpEqual,
                Value:    config.Tolerations[0].Value,
                Effect:   corev1.TaintEffectNoSchedule,
            }
            pod.Spec.Tolerations = append(pod.Spec.Tolerations, toleration)
        }
    }
    mutatedJSON, err := json.Marshal(pod)
    if err != nil {
        helpers.HandleError(w, r, fmt.Errorf("marshal pod: %v", err))
        return
    }
	patch, err := jsondiff.Compare(originalJSON, mutatedJSON)
	if err != nil {
        helpers.HandleError(w, r, fmt.Errorf("create JSON patch: %v", err))
        return
	}
    patchb, err := json.Marshal(patch)
    patchType := admissionv1.PatchTypeJSONPatch
    admissionResponse := &admissionv1.AdmissionResponse{
        UID:       admissionReview.Request.UID,
        Allowed:   true,
        Patch:     patchb,
        PatchType: &patchType,
    }
    admissionReview.Response = admissionResponse
    respBytes, err := json.Marshal(admissionReview)
    if err != nil {
        helpers.HandleError(w, r, fmt.Errorf("marshal admission review: %v", err))
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.Write(respBytes)
}