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

SCALEIO_CLIENT_CONFIGURATION_FILEPATH=${PWD}/scaleio_client_config.yml
cat > ${SCALEIO_CLIENT_CONFIGURATION_FILEPATH} <<EOF
libstorage:
  host: ${LIBSTORAGE_HOST_URL}
  storage:
    driver: scaleio
scaleio:
  endpoint:             ${SCALEIO_ENDPOINT}
  insecure:             true
  useCerts:             false
  userName:             ${SCALEIO_USERNAME}
  password:             ${SCALEIO_PASSWORD}
  systemID:             ${SCALEIO_SYSTEMID}
  protectionDomainID:   ${SCALEIO_PROTECTIONDOMAIN_ID}
  protectionDomainName: ${SCALEIO_PROTECTIONDOMAIN_NAME}
  storagePoolName:      ${SCALEIO_STORAGEPOOL_NAME}
  thinOrThick:          ${SCALEIO_THINORTHICK}
  version:              ${SCALEIO_VERSION}
EOF

pushd gocode/src/github.com/EMC-CMD/cf-persist-service-broker
  godep restore
  go run main.go &
  ginkgo -r
popd
