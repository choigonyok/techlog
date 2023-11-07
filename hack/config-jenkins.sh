#!bin/bash

# 테라폼 output의 jenkins ebs volume id를 json 형식으로 출력
terraform output -json > output.json

# 출력된 파일에서 id를 가져와 k8s configMap manifest에 입력
jq --argjson new_output "$(cat metadata.json)" '.spec.csi.volumeAttributes = $new_output' your-pv.yaml > updated-pv.yaml

# configMap 리소스와 jenkins 오브젝트 생성
kubectl apply -f configmap.yml
kubectl apply -f jenkins.yml

# 잘 생성되었는지 확인