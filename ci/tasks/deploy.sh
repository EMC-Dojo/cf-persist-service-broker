#!/usr/bin/env bash

set -e -x

source cf-persist-service-broker/ci/tasks/util.sh

check_param CF_IP
check_param CF_ENDPOINT
check_param CF_USERNAME
check_param CF_PASSWORD
check_param CF_ORG
check_param CF_SPACE
check_param BROKER_USERNAME
check_param BROKER_PASSWORD
check_param SCALEIO_ENDPOINT
check_param SCALEIO_USERNAME
check_param SCALEIO_PASSWORD

curl -L "https://cli.run.pivotal.io/stable?release=linux64-binary&source=github" | tar -zx
export PATH=${PWD}:$PATH

cd cf-persist-service-broker/

echo "${CF_IP} api.${CF_ENDPOINT}" >> /etc/hosts
echo "${CF_IP} login.${CF_ENDPOINT}" >> /etc/hosts
echo "${CF_IP} cf-persist-service-broker-acceptance.${CF_ENDPOINT}" >> /etc/hosts

cf api "https://api.${CF_ENDPOINT}" --skip-ssl-validation
cf auth $CF_USERNAME "$CF_PASSWORD"
cf target -o $CF_ORG -s $CF_SPACE
cf push cf-persist-service-broker-acceptance --no-start -b go_buildpack
cf set-env cf-persist-service-broker-acceptance BROKER_USERNAME $BROKER_USERNAME
cf set-env cf-persist-service-broker-acceptance BROKER_PASSWORD $BROKER_PASSWORD
cf set-env cf-persist-service-broker-acceptance SCALEIO_ENDPOINT $SCALEIO_ENDPOINT
cf set-env cf-persist-service-broker-acceptance SCALEIO_USERNAME $SCALEIO_USERNAME
cf set-env cf-persist-service-broker-acceptance SCALEIO_PASSWORD $SCALEIO_PASSWORD
cf start cf-persist-service-broker-acceptance
