#!/usr/bin/env bash

set -e -x
source cf-persist-service-broker/ci/tasks/util.sh

check_param BROKER_NAME
check_param BROKER_USERNAME
check_param BROKER_PASSWORD
check_param CF_ENDPOINT
check_param CF_USERNAME
check_param CF_PASSWORD
check_param CF_SERVICE
check_param CF_ORG
check_param CF_SPACE
check_param LIBSTORAGE_URI
check_param LIB_STOR_SERVICE
check_param INSECURE
check_param EMC_SERVICE_NAME
check_param EMC_SERVICE_UUID
check_param DIEGO_DRIVER_SPEC
check_param TEST_APP_NAME

cd cf-persist-service-broker

#authentication stuff
cf api http://api.$CF_ENDPOINT --skip-ssl-validation

cf auth $CF_USERNAME $CF_PASSWORD

cf target -o $CF_ORG -s $CF_SPACE

#Push EMC-Persistence broker with '--no-start' to allow setting ENVironment
cf push $BROKER_NAME --no-start

#Set ENVironment for EMC-Persistence broker

cf set-env $BROKER_NAME EMC_SERVICE_UUID $EMC_SERVICE_UUID
cf set-env $BROKER_NAME EMC_SERVICE_NAME $BROKER_NAME
cf set-env $BROKER_NAME LIBSTORAGE_URI $LIBSTORAGE_URI
cf set-env $BROKER_NAME DIEGO_DRIVER_SPEC $DIEGO_DRIVER_SPEC
cf set-env $BROKER_NAME LIB_STOR_SERVICE $LIB_STOR_SERVICE
cf set-env $BROKER_NAME INSECURE $INSECURE
cf set-env $BROKER_NAME BROKER_USERNAME $BROKER_USERNAME
cf set-env $BROKER_NAME BROKER_PASSWORD $BROKER_PASSWORD


#Start EMC-Persistence broker with correct ENVironment
cf start $BROKER_NAME

#Create Service Broker for use with CF
cf create-service-broker $BROKER_NAME $BROKER_USERNAME $BROKER_PASSWORD http://$BROKER_NAME.$CF_ENDPOINT
#Enable EMC-Persistence Service for use with CF
cf enable-service-access $BROKER_NAME

cd ../lifecycle-app
cf push $TEST_APP_NAME --no-start
cf set-env $TEST_APP_NAME CF_SERVICE $CF_SERVICE

get_cf_services |
while read service;
  do
  cf create-service $BROKER_NAME $service $service'_TEST_INSTANCE'
  cf bind-service $TEST_APP_NAME $service'_TEST_INSTANCE'
  cf start $TEST_APP_NAME
  curl -X POST -F 'text_box=Concourse BOT was here' http://$TEST_APP_NAME.$CF_ENDPOINT  | grep -w "Concourse BOT was here"
  cf stop $TEST_APP_NAME
  cf unbind-service $TEST_APP_NAME
  cf restage $TEST_APP_NAME
  curl http://$TEST_APP_NAME.$CF_ENDPOINT | grep -w "can't open file"
done;

get_cf_services |
while read service
  do cf delete-service $service'_TEST_INSTANCE' -f
done;

cf delete-service-broker $BROKER_NAME -f
cf delete $BROKER_NAME -f
cf delete $TEST_APP_NAME -f
