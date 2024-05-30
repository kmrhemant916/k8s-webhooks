package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"

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
        log.Printf("Error decoding admission review: %v", err)
        helpers.HandleError(w, r, err)
        return
    }
    pod := &corev1.Pod{}
    if err := json.Unmarshal(admissionReview.Request.Object.Raw, pod); err != nil {
        log.Printf("Error unmarshalling pod: %v", err)
        helpers.HandleError(w, r, fmt.Errorf("unmarshal to pod: %v", err))
        return
    }
    log.Printf("Original Pod: %+v", pod)
    originalJSON := admissionReview.Request.Object.Raw
    config, err := helpers.ReadConfig(Config)
    if err != nil {
        log.Printf("Error reading config: %v", err)
        panic(err)
    }
    for _, label := range config.TargetLabels {
        if value, exists := pod.Labels[label.Key]; exists && value == label.Value {
            log.Printf("Label %s=%s matched", label.Key, value)
            toleration := corev1.Toleration{
                Key:      config.Tolerations[0].Key,
                Operator: corev1.TolerationOpEqual,
                Value:    config.Tolerations[0].Value,
                Effect:   corev1.TaintEffectNoSchedule,
            }
            pod.Spec.Tolerations = append(pod.Spec.Tolerations, toleration)
            break
        }
    }
    if pod.Spec.NodeSelector == nil {
        pod.Spec.NodeSelector = make(map[string]string)
    }
    selectorValues := reflect.ValueOf(config.NodeSelector)
    for i := 0; i < selectorValues.NumField(); i++ {
        key := selectorValues.Type().Field(i).Tag.Get("yaml")
        value := selectorValues.Field(i).String()
        pod.Spec.NodeSelector[key] = value
    }
    mutatedJSON, err := json.Marshal(pod)
    if err != nil {
        log.Printf("Error marshalling pod: %v", err)
        helpers.HandleError(w, r, fmt.Errorf("marshal pod: %v", err))
        return
    }
    log.Printf("Mutated Pod JSON: %s", string(mutatedJSON))
	patch, err := jsondiff.Compare(originalJSON, mutatedJSON)
	if err != nil {
        log.Printf("Error creating JSON patch: %v", err)
        helpers.HandleError(w, r, fmt.Errorf("create JSON patch: %v", err))
        return
	}
    patchb, err := json.Marshal(patch)
    log.Printf("JSON Patch: %s", string(patchb))
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
        log.Printf("Error marshalling admission review response: %v", err)
        helpers.HandleError(w, r, fmt.Errorf("marshal admission review: %v", err))
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.Write(respBytes)
}