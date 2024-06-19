CA_BUNDLE=$(kubectl get secret scheduler-webhook-cert-secret -o jsonpath='{.data.ca\.crt}')
