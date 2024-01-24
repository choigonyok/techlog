#!bin/bash

if [[ -f ../.env ]]; then
  echo ".env file found"
  source ../.env

else
  echo ".env file not found"
  exit 1
fi

kubectl create secret generic oauth2-proxy-secrets \
--from-literal=cookie-secret=$COOKIE_SECRET \
--from-literal=client-id=$CLIENT_ID \
--from-literal=client-secret=$CLIENT_SECRET