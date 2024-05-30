package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	admissionv1 "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestMutate(t *testing.T) {
	// Create a new instance of the App
	app := &App{}

	// Create a sample admission review request
	admissionReview := &admissionv1.AdmissionReview{
		Request: &admissionv1.AdmissionRequest{
			Object: runtime.RawExtension{Raw: []byte(`{"metadata":{"name":"mypod","labels":{"app":"test"}},"spec":{"containers":[{"name":"mycontainer","image":"busybox","command":["sleep","3600"]}]},"status":{}}`)},
		},
	}

	// Marshal the admission review request
	admissionReviewBytes, err := json.Marshal(admissionReview)
	assert.NoError(t, err)

	// Create a new HTTP request with the admission review request as the body
	req := httptest.NewRequest(http.MethodPost, "/mutate", bytes.NewReader(admissionReviewBytes))
	req.Header.Set("Content-Type", "application/json")

	// Create a new HTTP response recorder
	w := httptest.NewRecorder()

	// Call the Mutate function with the recorder and request
	app.Mutate(w, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Unmarshal the response body into an AdmissionReview object
	var responseAdmissionReview admissionv1.AdmissionReview
	err = json.NewDecoder(w.Body).Decode(&responseAdmissionReview)
	assert.NoError(t, err)

	// Check that the response contains a patch
	assert.NotNil(t, responseAdmissionReview.Response.Patch)
}
