#!bin/bash

if [[ -f ../.env ]]; then
  echo ".env file found"
  source ../.env

else
  echo ".env file not found"
  exit 1
fi

# For oauth2-proxy environment variables
kubectl create secret generic oauth2-proxy-secrets \
--from-literal=cookie-secret=$COOKIE_SECRET \
--from-literal=client-id=$CLIENT_ID \
--from-literal=client-secret=$CLIENT_SECRET

# For backend environment variables
kubectl create secret generic backend-secrets \
--from-literal=environment=$ENVIRONMENT \
--from-literal=db-host-read=$DB_HOST_READ \
--from-literal=db-host-write=$DB_HOST_WRITE \
--from-literal=db-driver=$DB_DRIVER \
--from-literal=host=$HOST \
--from-literal=db-user=$DB_USER \
--from-literal=db-password=$DB_PASSWORD \
--from-literal=db-name=$DB_NAME \
--from-literal=github-token=$GITHUB_TOKEN \
--from-literal=aws-access-key=$AWS_ACCESS_KEY \
--from-literal=aws-secret-key=$AWS_SECRET_KEY