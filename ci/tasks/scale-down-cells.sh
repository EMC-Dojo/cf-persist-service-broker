#!/usr/bin/env bash

set -e -x
source cf-persist-service-broker/ci/tasks/util.sh

check_param BOSH_DIR
check_param BOSH_PASS
check_param BOSH_USER
check_param CI_DIEGOCELL_IPS
check_param DIEGO_DEPLOYMENT_NAME
check_param SCALEIO_MDM_IPS

#install bosh cli (should add to docker image eventually...)
gem install bosh_cli --no-ri --no-rdoc

#Setup BOSH CLI for us
bosh -n target ${BOSH_DIR}
bosh -n login ${BOSH_USER} ${BOSH_PASS}

manifest="$(bosh download manifest ${DIEGO_DEPLOYMENT_NAME} | perl -0pe 's/(- instances: 2\n  name: CI_cell_z1)(.+?)(- instances:)/- instances:/sg')"
echo "${manifest}" > manifest.yml

bosh deployment manifest.yml
bosh -n deploy

echo "Diego has now been tamed! Meow~"
