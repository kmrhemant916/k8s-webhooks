#!/bin/bash
TAG=$(date +"%H%M%S")
docker build -t iad.ocir.io/idt7ybnr03cb/scheduler-webhook:${TAG} .
echo "Pushing - iad.ocir.io/idt7ybnr03cb/scheduler-webhook:${TAG}"
docker push iad.ocir.io/idt7ybnr03cb/scheduler-webhook:${TAG}