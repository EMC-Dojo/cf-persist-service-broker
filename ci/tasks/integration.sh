#!/usr/bin/env bash

set -e -x

source cf-persist-service-broker/ci/tasks/util.sh

check_param BROKER_USERNAME
check_param BROKER_PASSWORD
check_param LIBSTORAGE_HOST_URL
check_param SCALEIO_ENDPOINT
check_param SCALEIO_USERNAME
check_param SCALEIO_PASSWORD
check_param SCALEIO_SYSTEMID
check_param SCALEIO_PROTECTIONDOMAIN_ID
check_param SCALEIO_PROTECTIONDOMAIN_NAME
check_param SCALEIO_STORAGEPOOL_NAME
check_param SCALEIO_THINORTHICK
check_param SCALEIO_VERSION

export GOPATH=$PWD/gocode
export PATH=$PATH:$GOPATH/bin

mkdir -p gocode/src/github.com/EMC-CMD
cp -r cf-persist-service-broker gocode/src/github.com/EMC-CMD/cf-persist-service-broker

pushd gocode/src/github.com/EMC-CMD/cf-persist-service-broker
  godep restore
  go run main.go &
  ginkgo -r
popd
