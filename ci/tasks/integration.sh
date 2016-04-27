#!/usr/bin/env bash

set -e -x

source cf-persist-service-broker/ci/tasks/util.sh

check_param PORT
check_param BROKER_USERNAME
check_param BROKER_PASSWORD
check_param SCALEIO_SERVICE_BROKER_SERVER_URL

export GOPATH=$PWD/gocode
export PATH=$PATH:$GOPATH/bin

mkdir -p gocode/src/github.com/EMC-CMD
cp -r cf-persist-service-broker gocode/src/github.com/EMC-CMD/cf-persist-service-broker

pushd gocode/src/github.com/EMC-CMD/cf-persist-service-broker
godep restore
go run main.go &
ginkgo -r
popd
