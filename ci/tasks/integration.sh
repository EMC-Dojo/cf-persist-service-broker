#!/usr/bin/env bash

set -e -x

source cf-persist-service-broker/ci/tasks/util.sh

check_param TEST_INSTANCE_ID
check_param PARSED_INSTANCE_ID
check_param TEST_SIZE
check_param STORAGE_POOL_NAME
check_param BROKER_PASSWORD
check_param BROKER_USERNAME
check_param BROKER_PORT
check_param LIBSTORAGE_URI
check_param LIBSTORAGE_DRIVER_TYPE
check_param INSECURE

export GOPATH=$PWD/gocode
export PATH=$PATH:$GOPATH/bin

godep restore
mkdir -p gocode/src/github.com/EMC-Dojo
cp -r cf-persist-service-broker gocode/src/github.com/EMC-Dojo

pushd gocode/src/github.com/EMC-Dojo/cf-persist-service-broker
  ginkgo -r
popd
