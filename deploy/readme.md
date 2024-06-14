CA_BUNDLE=$(kubectl get secret pod-mutating-webhook-tls -o jsonpath='{.data.ca\.crt}')
