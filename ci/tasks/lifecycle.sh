#!/usr/bin/env bash

set -e -x
source cf-persist-service-broker/ci/tasks/util.sh

check_param RACK_ENV
check_param BROKER_PASSWORD
check_param BROKER_USERNAME
check_param CF_ENDPOINT
check_param CF_IP
check_param CF_ORG
check_param CF_PASSWORD
check_param CF_SCALEIO_SB_APP
check_param CF_SCALEIO_SB_SERVICE
check_param CF_SPACE
check_param CF_USERNAME
check_param LIBSTORAGE_HOST
check_param LIBSTORAGE_STORAGE_DRIVER
check_param SCALEIO_ENDPOINT
check_param SCALEIO_INSECURE
check_param SCALEIO_PASSWORD
check_param SCALEIO_PROTECTION_DOMAIN_ID
check_param SCALEIO_PROTECTION_DOMAIN_NAME
check_param SCALEIO_STORAGE_POOL_NAME
check_param SCALEIO_SYSTEM_ID
check_param SCALEIO_SYSTEM_NAME
check_param SCALEIO_THIN_OR_THICK
check_param SCALEIO_USE_CERTS
check_param SCALEIO_USERNAME
check_param SCALEIO_VERSION

cd cf-persist-service-broker/

bundle install
bundle exec rspec spec --tag type:lifecycle
