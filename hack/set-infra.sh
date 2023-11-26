#!bin/bash

echo "Create terraform plan file..."
cd ../eks 
terraform plan -out terraform.plan

echo "Apply terraform..."
terraform apply terraform.plan
echo "Removing terraform plan file..."
rm terraform.plan

echo "Connect local kubectl..."
REGION=$(terraform output region | sed s/\"//g)
CLUSTER_NAME=$(terraform output cluster_name | sed s/\"//g)
echo "Connecting local kubectl to EKS cluster..."
aws eks --region $REGION update-kubeconfig --name $CLUSTER_NAME

echo "Create namespaces..."
kubectl create ns devops-system
kubectl create ns argocd

echo "Create nginx ingress controller..."
kubectl apply -f ../hack/manifests/ingress-controller.yml
HTTP_ARN=$(terraform output target_group_http_arn | sed 's/\//\\\//g')
# HTTPS_ARN=$(terraform output target_group_https_arn | sed 's/\//\\\//g')
sed -i '' "s/targetGroupARN: HTTP_ARN/targetGroupARN: $HTTP_ARN/" ../hack/manifests/target-group-binding.yml
# sed -i '' "s/targetGroupARN: HTTPS_ARN/targetGroupARN: $HTTPS_ARN/" ../hack/manifests/target-group-binding.yml
kubectl apply -f ../hack/manifests/target-group-binding.yml
sed -i '' 's/targetGroupARN: .*/targetGroupARN: HTTP_ARN/g' ../hack/manifests/target-group-binding.yml
# sed -i '' "1,/---/s/targetGroupARN: .*/targetGroupARN: HTTP_ARN/" ../hack/manifests/target-group-binding.yml

echo "Deploy test ingress and nginx..."
kubectl apply -f ../hack/manifests/ingress.yml
kubectl apply -f ../hack/manifests/nginx-test.yml

echo "Deploy K8S metric server..."
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/high-availability-1.21+.yaml

# echo "Deploy jenkins as Loadbalancer type service object..."
# FILE_SYSTEM_ID=$(terraform output efs_id | sed s/\"//g )
# sed -i '' "s/fileSystemId: FILE_SYSTEM_ID/fileSystemId: $FILE_SYSTEM_ID/" ../hack/manifests/jenkins.yml
# kubectl apply -f ../hack/manifests/jenkins.yml
# sed -i '' "s/fileSystemId: .*/fileSystemId: FILE_SYSTEM_ID/" ../hack/manifests/jenkins.yml

# echo "Issueing ServiceAccount token..."
# kubectl create token jenkins -n devops-system

# echo "Deploying ConfigMap for Kaniko..."
# kubectl create configmap config.json --from-file=../hack/config.json -n devops-system

# echo "Deploying ArgoCD..."
# kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
# # kubectl patch svc argocd-server -n argocd -p '{"spec": {"type": "LoadBalancer"}}'
# kubectl patch svc argocd-server -n argocd -p '{"spec": {"type": "NodePort"}}'
# kubectl patch svc argocd-server -n argocd -p '{"spec": {"ports": [{"name": "https", "nodePort": 30000, "port": 443, "targetPort": 8080}]}}'
# PASSWORD=$(kubectl get secret argocd-initial-admin-secret -n argocd -o json | jq -r '.data.password' | base64 -d)
# echo "PASSWORD: $PASSWORD"
# PASSWORD=""

# kubectl patch deployment argocd-redis -n argocd -p \
#   '{"spec":{"template":{"spec":{"volumes":[{"name":"argocd-redis","persistentVolumeClaim":{"claimName":"argocd"}}]}}}}'

# # kubectl patch deployment argocd-redis -n argocd -p \
# #   '{"spec":{"template":{"spec":{"containers": {"volumeMounts": {"name":"argocd-redis","mountPath":"/datas"}}}}}}'