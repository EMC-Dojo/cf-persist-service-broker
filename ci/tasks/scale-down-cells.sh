#!/usr/bin/env bash

set -e -x
source cf-persist-service-broker/ci/tasks/util.sh

check_param BOSH_DIRECTOR_PUBLIC_IP
check_param BOSH_PASSWORD
check_param BOSH_USER
check_param DIEGO_DEPLOYMENT_NAME

#install bosh cli (should add to docker image eventually...)
gem install bosh_cli --no-ri --no-rdoc

#Setup BOSH CLI for us
bosh -n target ${BOSH_DIRECTOR_PUBLIC_IP}
bosh -n login ${BOSH_USER} ${BOSH_PASSWORD}

manifest="$(bosh download manifest ${DIEGO_DEPLOYMENT_NAME} | perl -0pe 's/(- instances: 2\n  name: CI_cell_z1)(.+?)(- instances:)/- instances:/sg')"
echo "${manifest}" > manifest.yml

bosh deployment manifest.yml
bosh -n deploy

echo "Diego has now been tamed! Meow~"
