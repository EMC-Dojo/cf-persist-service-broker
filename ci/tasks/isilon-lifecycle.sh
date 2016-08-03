#!/usr/bin/env bash

set -e -x
source cf-persist-service-broker/ci/tasks/util.sh

check_param LIFECYCLE_APP_MEMORY
check_param BROKER_NAME
check_param BROKER_PASSWORD
check_param BROKER_USERNAME
check_param CF_ENDPOINT
check_param CF_ORG
check_param CF_PASSWORD
check_param CF_SERVICE
check_param CF_SPACE
check_param CF_USERNAME
check_param DIEGO_DRIVER_SPEC
check_param EMC_SERVICE_NAME
check_param EMC_SERVICE_UUID
check_param INSECURE
check_param LIB_STOR_SERVICE
check_param LIBSTORAGE_URI
check_param LIFECYCLE_APP_NAME
check_param NUM_DIEGO_CELLS

#Setup CF CLI for us
cf api http://api.$CF_ENDPOINT --skip-ssl-validation
cf auth $CF_USERNAME $CF_PASSWORD
cf target -o $CF_ORG -s $CF_SPACE

get_cf_service |
while read service
  do
  set -e -x
  cf delete-service $service'_TEST_INSTANCE' -f
done;
cf delete-service-broker $BROKER_NAME -f
cf delete $BROKER_NAME -f
cf delete $LIFECYCLE_APP_NAME -f
cf delete-orphaned-routes -f

pushd cf-persist-service-broker
#Push EMC-Persistence broker with '--no-start' to allow setting ENVironment
cf push $BROKER_NAME --no-start

#Set ENVironment for EMC-Persistence broker
cf set-env $BROKER_NAME BROKER_PASSWORD $BROKER_PASSWORD
cf set-env $BROKER_NAME BROKER_USERNAME $BROKER_USERNAME
cf set-env $BROKER_NAME DIEGO_DRIVER_SPEC $DIEGO_DRIVER_SPEC
cf set-env $BROKER_NAME EMC_SERVICE_NAME $BROKER_NAME
cf set-env $BROKER_NAME EMC_SERVICE_UUID $EMC_SERVICE_UUID
cf set-env $BROKER_NAME INSECURE $INSECURE
cf set-env $BROKER_NAME LIB_STOR_SERVICE $LIB_STOR_SERVICE
cf set-env $BROKER_NAME LIBSTORAGE_URI $LIBSTORAGE_URI

#Start EMC-Persistence broker with correct ENVironment
cf start $BROKER_NAME

#Create Service Broker for use with CF & Enable EMC-Persistence Service for use with CF
cf create-service-broker $BROKER_NAME $BROKER_USERNAME $BROKER_PASSWORD http://$BROKER_NAME.$CF_ENDPOINT
cf enable-service-access $BROKER_NAME
popd

pushd lifecycle-app
cf push $LIFECYCLE_APP_NAME --no-start
cf set-env $LIFECYCLE_APP_NAME CF_SERVICE $CF_SERVICE

echo "1" > status.txt

get_cf_service |
while read service
  do
  set -x -e
  cf create-service $BROKER_NAME $service $service'_TEST_INSTANCE'
  cf bind-service $LIFECYCLE_APP_NAME $service'_TEST_INSTANCE'
  cf start $LIFECYCLE_APP_NAME
  curl -X POST -F 'text_box=Concourse BOT was here' http://$LIFECYCLE_APP_NAME.$CF_ENDPOINT | grep -w "Concourse BOT was here"
  cf scale $LIFECYCLE_APP_NAME -i $NUM_DIEGO_CELLS -m $LIFECYCLE_APP_MEMORY -f
    for i in `seq 0 $[$NUM_DIEGO_CELLS*10]`
    do
      set -x -e
      curl_output="$(curl http://$LIFECYCLE_APP_NAME.$CF_ENDPOINT)"
      echo "$curl_output" | grep -w "Concourse BOT was here"
      instance_number="$(echo $curl_output | grep "Instance ID is: " | sed -n -e 's/^.*Instance\ ID\ is:\ //p' | cut -f 1 -d '<')"
      instances[$instance_number]=1
      if [ "${#instances[@]}" == $NUM_DIEGO_CELLS ]
      then
        echo "0" > status.txt
        break
      fi
    done;
  cf stop $LIFECYCLE_APP_NAME
  cf unbind-service $LIFECYCLE_APP_NAME $service'_TEST_INSTANCE'
  cf restage $LIFECYCLE_APP_NAME
  curl http://$LIFECYCLE_APP_NAME.$CF_ENDPOINT | grep -w "can't open file"
done;

if [ "$(cat status.txt)" == 1 ]
then
  echo "Didnt Verify Across All Diego Cells After Multiple Tries!"
  exit 1
else
  echo "Verified All Diego Cells Are Communicating with Isilon Device"
fi
popd
