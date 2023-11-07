#!bin/bash

# 테라폼 output의 jenkins ebs volume id를 json 형식으로 출력
terraform output -json > output.json

# 출력된 파일에서 id를 가져와 k8s configMap manifest에 입력
jq --argjson new_output "$(cat metadata.json)" '.spec.csi.volumeAttributes = $new_output' your-pv.yaml > updated-pv.yaml

MASTER_NODE_IP=$(jq -r '.master-ip.value' output.json)
WORKER_NODE1_IP=$(jq -r '.worker1-ip.value' output.json)
WORKER_NODE2_IP=$(jq -r '.worker2-ip.value' output.json)

EBS_ID=$(jq -r '.jenkins-ebs-id.value' output.json)

sed -e 's/{ebs-volume-id}/$EBS_ID%/g' template.yml > configmap.yml