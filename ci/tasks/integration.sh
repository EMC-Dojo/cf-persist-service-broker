#!/usr/bin/env bash

set -e -x
source cf-persist-service-broker/ci/tasks/util.sh

check_param TEST_INSTANCE_ID
check_param PARSED_INSTANCE_ID
check_param STORAGE_POOL_NAME
check_param BROKER_PASSWORD
check_param BROKER_USERNAME
check_param LIBSTORAGE_URI
check_param LIB_STOR_SERVICE
check_param PORT
check_param EMC_SERVICE_UUID
check_param EMC_SERVICE_NAME
check_param DIEGO_DRIVER_SPEC
check_param INSECURE

export GOPATH=$PWD/gocode
export PATH=$PATH:$GOPATH/bin

go get github.com/onsi/ginkgo/ginkgo

mkdir -p gocode/src/github.com/EMC-Dojo
cp -r cf-persist-service-broker gocode/src/github.com/EMC-Dojo

pushd gocode/src/github.com/EMC-Dojo/cf-persist-service-broker
  ginkgo -r
popd
