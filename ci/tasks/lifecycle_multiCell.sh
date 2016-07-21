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
check_param NUM_DIEGO_CELLS
check_param APP_MEMORY
check_param DIEGO_DEPLOYMENT_NAME
check_param CI_DIEGOCELL_IPS
check_param SCALEIO_MDM_IPS
check_param BOSH_DIR
check_param BOSH_USER
check_param BOSH_PASS

cd cf-persist-service-broker

#install bosh cli (should add to docker image eventually...)
gem install bosh_cli --no-ri --no-rdoc

#Setup BOSH CLI for us
bosh -n target $BOSH_DIR
bosh -n login $BOSH_USER $BOSH_PASS

#Setup CF CLI for us
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

echo "1" > status.txt
boshDiegoManifest="$(bosh download manifest $DIEGO_DEPLOYMENT_NAME)"
CI_BoshDiegoManifest="$(echo -e "${boshDiegoManifest}" | sed -e $'s/jobs:/jobs:\\\n- instances: 2\\\n  name: CI_cell_z1\\\n  networks:\\\n  - name: private\\\n    static_ips: ['"$CI_DIEGOCELL_IPS"$']\\\n  properties:\\\n    scaleio:\\\n      mdm:\\\n        ips: ['"$SCALEIO_MDM_IPS"$']\\\n    diego:\\\n      rep:\\\n        zone: z1\\\n    metron_agent:\\\n      zone: z1\\\n  vm_type: x-large\\\n  stemcell: trusty-3215\\\n  azs:\\\n  - z1\\\n  templates:\\\n  - name: consul_agent\\\n    release: cf\\\n  - name: rep\\\n    release: diego-release\\\n  - name: garden\\\n    release: garden-linux\\\n  - name: cflinuxfs2-rootfs-setup\\\n    release: cflinuxfs2-rootfs\\\n  - name: metron_agent\\\n    release: cf\\\n  - name: rexray_service\\\n    release: rexray-bosh-release\\\n  - name: setup_sdc\\\n    release: scaleio-sdc-bosh-release\\\n  update:\\\n    max_in_flight: 1\\\n    serial: false/g')"
echo "${boshDiegoManifest}" > baseBOSHDiegoManifest.yml
echo "${CI_BoshDiegoManifest}" > CI_BOSHDiegoManigest.yml
bosh deployment CI_BOSHDiegoManigest.yml
bosh -n deploy

get_cf_service |
while read service
  do
  set -x -e
  cf create-service $BROKER_NAME $service $service'_TEST_INSTANCE'
  cf bind-service $TEST_APP_NAME $service'_TEST_INSTANCE'
  cf start $TEST_APP_NAME
  curl -X POST -F 'text_box=Concourse BOT was here' http://$TEST_APP_NAME.$CF_ENDPOINT | grep -w "Concourse BOT was here"

  cf scale $TEST_APP_NAME -i $NUM_DIEGO_CELLS -m $APP_MEMORY -f
    for i in `seq 0 $[$NUM_DIEGO_CELLS*10]`
    do
      set -x -e
      curl_output="$(curl http://$TEST_APP_NAME.$CF_ENDPOINT)"
      echo "$curl_output" | grep -w "Concourse BOT was here"
      instance_number="$(echo $curl_output | grep "Instance ID is: " | sed -n -e 's/^.*Instance\ ID\ is:\ //p' | cut -f 1 -d '<')"
      instances[$instance_number]=1
      if [ "${#instances[@]}" == $NUM_DIEGO_CELLS ]
      then
        echo "0" > status.txt
        break
      fi
    done;

  cf stop $TEST_APP_NAME
  cf unbind-service $TEST_APP_NAME $service'_TEST_INSTANCE'
  cf restage $TEST_APP_NAME
  curl http://$TEST_APP_NAME.$CF_ENDPOINT | grep -w "can't open file"
done;

get_cf_service |
while read service
  do
  set -e -x
  cf delete-service $service'_TEST_INSTANCE' -f
done;

cf delete-service-broker $BROKER_NAME -f
cf delete $BROKER_NAME -f
cf delete $TEST_APP_NAME -f

bosh deployment baseBOSHDiegoManifest.yml
bosh -n deploy

if [ "$(cat status.txt)" == 1 ]
then
  echo "Didnt Verify Across All Diego Cells After Multiple Tries!"
  exit 1
else
  echo "Verified All Diego Cells Are Communicating with Isilon Device"
fi
