#!/usr/bin/env bash

set -e -x

source cf-persist-service-broker/ci/tasks/util.sh

check_param BROKER_USERNAME
check_param BROKER_PASSWORD
check_param LIBSTORAGE_HOST
check_param LIBSTORAGE_STORAGE_DRIVER
check_param SCALEIO_ENDPOINT
check_param SCALEIO_USERNAME
check_param SCALEIO_PASSWORD
check_param SCALEIO_SYSTEM_ID
check_param SCALEIO_PROTECTION_DOMAIN_ID
check_param SCALEIO_PROTECTION_DOMAIN_NAME
check_param SCALEIO_STORAGE_POOL_NAME
check_param SCALEIO_THIN_OR_THICK
check_param SCALEIO_VERSION
check_param SCALEIO_INSECURE
check_param SCALEIO_USE_CERTS

export GOPATH=$PWD/gocode
export PATH=$PATH:$GOPATH/bin

go get github.com/onsi/ginkgo/ginkgo  # installs the ginkgo CLI
go get github.com/onsi/gomega         # fetches the matcher library
go get github.com/golang/mock/gomock  # gets the mocking library
mkdir -p gocode/src/github.com/EMC-Dojo
cp -r cf-persist-service-broker gocode/src/github.com/EMC-Dojo/cf-persist-service-broker

pushd gocode/src/github.com/EMC-Dojo/cf-persist-service-broker
  ginkgo -r
popd
