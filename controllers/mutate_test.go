package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/kmrhemant916/k8s-webhooks/helpers"
	"github.com/stretchr/testify/assert"
	"gomodules.xyz/jsonpatch/v2"
	admissionv1 "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Mocking the ReadConfig function
var mockConfig = &helpers.Config{
	Service: struct {
		Port string "yaml:\"port\""
	}{
		Port: "8080",
	},
	TargetLabels: []struct {
		Key   string "yaml:\"key\""
		Value string "yaml:\"value\""
	}{
		{Key: "app", Value: "test"},
	},
	Tolerations: []struct {
		Key      string "yaml:\"key\""
		Operator string "yaml:\"operator\""
		Value    string "yaml:\"value\""
		Effect   string "yaml:\"effect\""
	}{
		{Key: "example-key", Operator: "Equal", Value: "example-value", Effect: "NoSchedule"},
	},
	NodeSelector: struct {
		AgentPool string "yaml:\"agentpool\""
	}{
		AgentPool: "example-node",
	},
}

func mockReadConfig(filename string) (*helpers.Config, error) {
	return mockConfig, nil
}

func containsPatchOperation(patch []jsonpatch.Operation, expected jsonpatch.Operation) bool {
	for _, op := range patch {
		if op.Operation == expected.Operation && op.Path == expected.Path && reflect.DeepEqual(op.Value, expected.Value) {
			return true
		}
	}
	return false
}

func TestMutate(t *testing.T) {
	// Create a sample AdmissionReview
	admissionReview := &admissionv1.AdmissionReview{
		Request: &admissionv1.AdmissionRequest{
			UID: "test-uid",
			Object: runtime.RawExtension{
				Raw: []byte(`{
					"apiVersion": "v1",
					"kind": "Pod",
					"metadata": {
						"name": "test-pod",
						"labels": {
							"app": "test"
						}
					},
					"spec": {
						"containers": [
							{
								"name": "test-container",
								"image": "test-image"
							}
						]
					}
				}`),
			},
		},
	}

	body, err := json.Marshal(admissionReview)
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/mutate", bytes.NewReader(body))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	app := &App{ReadConfig: mockReadConfig}
	handler := http.HandlerFunc(app.Mutate)
	handler.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check the response
	respAdmissionReview := &admissionv1.AdmissionReview{}
	err = json.Unmarshal(rr.Body.Bytes(), respAdmissionReview)
	assert.NoError(t, err)

	// Check the patch
	patch := []jsonpatch.Operation{}
	err = json.Unmarshal(respAdmissionReview.Response.Patch, &patch)
	assert.NoError(t, err)

	expectedPatch := []jsonpatch.Operation{
		{
			Operation: "add",
			Path:      "/spec/tolerations",
			Value: []interface{}{
				map[string]interface{}{
					"key":      "example-key",
					"operator": "Equal",
					"value":    "example-value",
					"effect":   "NoSchedule",
				},
			},
		},
		{
			Operation: "add",
			Path:      "/spec/nodeSelector",
			Value: map[string]interface{}{
				"agentpool": "example-node",
			},
		},
	}

	for _, expectedOp := range expectedPatch {
		assert.True(t, containsPatchOperation(patch, expectedOp), "Patch does not contain expected operation: %v", expectedOp)
	}
}
