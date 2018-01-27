#!/bin/bash

kubectl delete -f ../templates/infra-clusterrolebinding.yaml
kubectl delete -f ../templates/infra-serviceaccount.yaml
kubectl delete -f ../templates/infra-namespace.yaml
