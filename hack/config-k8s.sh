#!bin/bash

# configMap 리소스 생성
kubectl apply -f configmap.yml

# devops-system namespace 생성
kubectl create namespace devops-system

# jenkins-master 배포
kubectl create namespace devops-systemkubectl apply -f jenkins.yml

# 잘 생성되었는지 확인