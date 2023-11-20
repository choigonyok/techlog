#!bin/bash

# For Test
kind create cluster

kubectl create namespace devops-system


# terraform으로 생성한 iam policy/role에 대한 arn 확인
aws iam list-instance-profiles | jq -r '.InstanceProfiles[].Roles[].Arn'

# kube2iam DaemonSet 배포

# EBS CSI Driver installation
echo "Install EBS CSI Driver..."
kubectl apply -k "github.com/kubernetes-sigs/aws-ebs-csi-driver/deploy/kubernetes/overlays/stable/?ref=release-1.24"

# 테라폼 output의 jenkins ebs volume id를 json 형식으로 출력
echo "Enter to terraform dir..."
cd ../terraform

echo "Touch terraform output file..."
terraform output -json > output.json

echo "Enter to manifest dir..."
cd ../hack/manifests

echo "Create EBS_ID variable..."
EBS_ID=$(jq -r '.["jenkins-ebs-id"].value' ../../terraform/output.json)

echo "Editting string..."
sed -e "s/{ebs-volume-id}/$EBS_ID/g" template.yml > pv-jenkins.yml

echo "Apply PersistenceVolume resource..."
kubectl apply -f pv-jenkins.yml

kubectl apply -f jenkins.yml

# 출력된 파일에서 id를 가져와 k8s configMap manifest에 입력
# MASTER_NODE_IP=$(jq -r '.master-ip.value' output.json)
# WORKER_NODE1_IP=$(jq -r '.worker1-ip.value' output.json)
# WORKER_NODE2_IP=$(jq -r '.worker2-ip.value' output.json)

