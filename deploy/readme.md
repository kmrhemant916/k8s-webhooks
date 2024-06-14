CA_BUNDLE=$(kubectl get secret openfaas-webhook-cert-secret -o jsonpath='{.data.ca\.crt}')
