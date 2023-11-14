#!bin/bash

echo "Creating terraform plan file..."
cd ../eks 
terraform plan -var region=ap-northeast-2 -out terraform.plan
EXIST= $(ls | grep 'terraform.plan')
echo "$EXIST"
if [ ${EXIST} -eq ${""} ] ; then
	echo "ERROR: Fail to find terraform.plan"
    exit 1
fi
echo "Applying terraform..."
terraform apply terraform.plan
echo "Removing terraform plan file..."
rm terraform.plan

# EXIST= $(ls | grep 'terraform.plan')
# echo "$EXIST"
# if [ ${EXIST} -eq ${""} ] ; then
# 	echo "ERROR: Fail to remove terraform.plan"
#     exit 1
# fi

REGION=$(terraform output region | sed s/\"//g)
CLUSTER_NAME=$(terraform output cluster_name | sed s/\"//g)
echo "Connecting local kubectl to EKS cluster..."
aws eks --region $REGION update-kubeconfig --name $CLUSTER_NAME

echo "Creating devops namespace..."
kubectl create ns devops-system

echo "Deploying jenkins as Loadbalancer type service object..."
kubectl apply -f ../hack/manifests/jenkins.yml

echo "Issueing ServiceAccount token..."
kubectl create token jenkins -n devops-system