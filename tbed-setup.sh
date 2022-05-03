#!/bin/bash

kubeconfig_file="$HOME/public/config-tbed"
gopath="$HOME/public/k8s-performance-tests/clusterloader2/pkg"

echo "Setting up environment for tbed cluster. "

export KUBECONFIG=$kubeconfig_file
echo "KUBECONFIG=$KUBECONFIG"

export GOPATH=$gopath
echo "GOPATH=$GOPATH" 

echo ""

echo "Test cluster connection: "

sleep 1
echo "K8s cluster nodes: "
kubectl get nodes -o wide

echo ""

sleep 1
echo "cluster-loader namespace: "
kubectl get all -o wide -n cluster-loader

