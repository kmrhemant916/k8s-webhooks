#!/bin/bash
TAG=$1
docker build -t iad.ocir.io/idt7ybnr03cb/hemant:${TAG} .
docker push iad.ocir.io/idt7ybnr03cb/hemant:${TAG}