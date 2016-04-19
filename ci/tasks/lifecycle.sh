#!/usr/bin/env bash

set -e -x

source cf-persist-service-broker/ci/tasks/util.sh

check_param RACK_ENV
check_param CF_IP
check_param CF_ENDPOINT
check_param CF_USERNAME
check_param CF_PASSWORD
check_param CF_ORG
check_param CF_SPACE
check_param SCALEIO_ENDPOINT
check_param SCALEIO_USERNAME
check_param SCALEIO_PASSWORD
check_param BROKER_USERNAME
check_param BROKER_PASSWORD

curl -L "https://cli.run.pivotal.io/stable?release=linux64-binary&source=github" | tar -zx
export PATH=${PWD}:$PATH

cd cf-persist-service-broker/

echo "${CF_IP} api.${CF_ENDPOINT}" >> /etc/hosts
echo "${CF_IP} login.${CF_ENDPOINT}" >> /etc/hosts
echo "${CF_IP} cf-persist-service-broker-lifecycle.${CF_ENDPOINT}" >> /etc/hosts

bundle install
bundle exec rspec spec --tag type:lifecycle
