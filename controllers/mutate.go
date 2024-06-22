package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/kmrhemant916/k8s-webhooks/helpers"
	"gomodules.xyz/jsonpatch/v2"
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
    podJSON, err := json.MarshalIndent(pod, "", "  ")
    if err != nil {
        log.Printf("Error marshalling pod to JSON: %v", err)
        helpers.HandleError(w, r, fmt.Errorf("marshal pod to JSON: %v", err))
        return
    }
    log.Printf("Original Pod: %s", string(podJSON))
    originalJSON := admissionReview.Request.Object.Raw
    config, err := app.ReadConfig(Config)
    if err != nil {
        log.Printf("Error reading config: %v", err)
        panic(err)
    }
    for _, label := range config.TargetLabels {
        if value, exists := pod.Labels[label.Key]; exists && value == label.Value {
            log.Printf("Label %s=%s matched", label.Key, value)
            if config.Patch.Tolerations.Enable {
                toleration := corev1.Toleration {
                    Key:      config.Patch.Tolerations.Value[0].Key,
                    Operator: corev1.TolerationOpEqual,
                    Value:    config.Patch.Tolerations.Value[0].Value,
                    Effect:   corev1.TaintEffectNoSchedule,
                }
                pod.Spec.Tolerations = append(pod.Spec.Tolerations, toleration)
            }
            if config.Patch.NodeSelector.Enable {
                if pod.Spec.NodeSelector == nil {
                    pod.Spec.NodeSelector = make(map[string]string)
                }
                selectorValues := reflect.ValueOf(config.Patch.NodeSelector.Value)
                for i := 0; i < selectorValues.NumField(); i++ {
                    key := selectorValues.Type().Field(i).Tag.Get("yaml")
                    value := selectorValues.Field(i).String()
                    pod.Spec.NodeSelector[key] = value
                }
            }
			if config.Patch.ImagePullSecrets.Enable {
				for _, secret := range config.Patch.ImagePullSecrets.Value {
					imagePullSecret := corev1.LocalObjectReference{
						Name: secret.Name,
					}
					pod.Spec.ImagePullSecrets = append(pod.Spec.ImagePullSecrets, imagePullSecret)
				}
			}          
            break
        }
    }
    mutatedJSON, err := json.MarshalIndent(pod, "", "  ")
    if err != nil {
        log.Printf("Error marshalling pod to JSON: %v", err)
        helpers.HandleError(w, r, fmt.Errorf("marshal pod to JSON: %v", err))
        return
    }
    log.Printf("Mutated Pod JSON: %s", string(mutatedJSON))
    patch, err := jsonpatch.CreatePatch(originalJSON, mutatedJSON)
    if err != nil {
        log.Printf("Error creating JSON patch: %v", err)
        helpers.HandleError(w, r, fmt.Errorf("create JSON patch: %v", err))
        return
    }
    patchBytes, err := json.MarshalIndent(patch, "", "  ")
    if err != nil {
        log.Printf("Error marshalling patch: %v", err)
        helpers.HandleError(w, r, fmt.Errorf("marshal patch: %v", err))
        return
    }
    log.Printf("JSON Patch: %s", string(patchBytes))
    patchType := admissionv1.PatchTypeJSONPatch
    admissionResponse := &admissionv1.AdmissionResponse{
        UID:       admissionReview.Request.UID,
        Allowed:   true,
        Patch:     patchBytes,
        PatchType: &patchType,
    }
    admissionReview.Response = admissionResponse
    respBytes, err := json.MarshalIndent(admissionReview, "", "  ")
    if err != nil {
        log.Printf("Error marshalling admission review response: %v", err)
        helpers.HandleError(w, r, fmt.Errorf("marshal admission review: %v", err))
        return
    }
    log.Printf("Admission Review Response: %s", string(respBytes))
    w.Header().Set("Content-Type", "application/json")
    w.Write(respBytes)
}