#!/bin/bash
TAG=$(date +"%H%M%S")
docker build -t iad.ocir.io/idt7ybnr03cb/hemant:${TAG} .
echo "Pushing - iad.ocir.io/idt7ybnr03cb/hemant:${TAG}"
docker push iad.ocir.io/idt7ybnr03cb/hemant:${TAG}