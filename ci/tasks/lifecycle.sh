#!/usr/bin/env bash

set -e -x
source cf-persist-service-broker/ci/tasks/util.sh

check_param BROKER_NAME
check_param BROKER_USERNAME
check_param BROKER_PASSWORD
check_param CF_ENDPOINT
check_param CF_USERNAME
check_param CF_PASSWORD
check_param CF_ORG
check_param CF_SPACE
check_param LIBSTORAGE_URI
check_param INSECURE

cd cf-persist-service-broker

#authentication stuff
cf api http://api.$CF_ENDPOINT --skip-ssl-validation

cf auth $CF_USERNAME $CF_PASSWORD

cf target -o $CF_ORG -s $CF_SPACE

#Push EMC-Persistence broker with '--no-start' to allow setting ENVironment
cf push $BROKER_NAME --no-start

#Set ENVironment for EMC-Persistence broker
cf set-env $BROKER_NAME BROKER_USERNAME $BROKER_USERNAME
cf set-env $BROKER_NAME BROKER_PASSWORD $BROKER_PASSWORD
cf set-env $BROKER_NAME LIBSTORAGE_URI $LIBSTORAGE_URI
cf set-env $BROKER_NAME INSECURE $INSECURE

#Start EMC-Persistence broker with correct ENVironment
cf start

#Create Service Broker for use with CF
cf create-service-broker $BROKER_NAME $BROKER_USERNAME $BROKER_PASSWORD http://$BROKER_NAME.$CF_ENDPOINT
#Enable EMC-Persistence Service for use with CF
cf enable-service-access EMC-Persistence

get_cf_services |
while read service;
  do cf create-service EMC-Persistence $service $service'_TEST_INSTANCE'
done;

get_cf_services |
while read service
  do cf delete-service $service'_TEST_INSTANCE' -f
done;

cf delete-service-broker $BROKER_NAME -f
