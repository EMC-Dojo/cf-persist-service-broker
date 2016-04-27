#!/usr/bin/env bash

set -e -x

source cf-persist-service-broker/ci/tasks/util.sh

check_param PORT
check_param BROKER_USERNAME
check_param BROKER_PASSWORD
check_param SCALEIO_SERVICE_BROKER_SERVER_URL


mkdir -p gocode/src/github.com/EMC-CMD
export GOPATH=$PWD/gocode
export GOROOT=/usr/local/go
export PATH=$PATH:$GOPATH/bin:$GOROOT/bin

cp -r cf-persist-service-broker gocode/src/github.com/EMC-CMD/cf-persist-service-broker
pushd gocode/src/github.com/EMC-CMD/cf-persist-service-broker

go get github.com/tools/godep
go install github.com/tools/godep
godep restore

go run main.go &

go get github.com/onsi/ginkgo/ginkgo
go install github.com/onsi/ginkgo/ginkgo
ginkgo -r

popd
