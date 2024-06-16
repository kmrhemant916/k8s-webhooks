package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type MockHelpers struct{}

func TestMutate(t *testing.T) {
	app := &App{}

	scheme := runtime.NewScheme()
	corev1.AddToScheme(scheme)

	tests := []struct {
		name                string
		pod                 *corev1.Pod
		expectedStatusCode  int
		expectedTolerations []corev1.Toleration
		expectedNodeSelector map[string]string
	}{
		{
			name: "Matching label should add toleration and node selector",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": "test"},
				},
				Spec: corev1.PodSpec{},
			},
			expectedStatusCode: http.StatusOK,
			expectedTolerations: []corev1.Toleration{
				{Key: "key1", Operator: corev1.TolerationOpEqual, Value: "value1", Effect: corev1.TaintEffectNoSchedule},
			},
			expectedNodeSelector: map[string]string{"role": "worker"},
		},
		{
			name: "Non-matching label should not add toleration or node selector",
			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": "non-matching"},
				},
				Spec: corev1.PodSpec{},
			},
			expectedStatusCode: http.StatusOK,
			expectedTolerations: nil,
			expectedNodeSelector: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			podBytes, err := json.Marshal(tt.pod)
			assert.NoError(t, err)

			admissionReview := admissionv1.AdmissionReview{
				Request: &admissionv1.AdmissionRequest{
					UID: "test-uid",
					Object: runtime.RawExtension{
						Raw: podBytes,
					},
				},
			}

			arBytes, err := json.Marshal(admissionReview)
			assert.NoError(t, err)

			req := httptest.NewRequest("POST", "/mutate", bytes.NewReader(arBytes))
			w := httptest.NewRecorder()

			app.Mutate(w, req)

			resp := w.Result()
			assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)

			if resp.StatusCode == http.StatusOK {
				var ar admissionv1.AdmissionReview
				err := json.NewDecoder(resp.Body).Decode(&ar)
				assert.NoError(t, err)

				var pod corev1.Pod
				err = json.Unmarshal(ar.Response.Patch, &pod)
				assert.NoError(t, err)

				assert.Equal(t, tt.expectedTolerations, pod.Spec.Tolerations)
				assert.Equal(t, tt.expectedNodeSelector, pod.Spec.NodeSelector)
			}
		})
	}
}
