#!/usr/bin/env bash

set -e -x
source cf-persist-service-broker/ci/tasks/util.sh

check_param BOSH_DIRECTOR_PUBLIC_IP
check_param BOSH_PASSWORD
check_param BOSH_USER
check_param DIEGO_CELL_IPS
check_param DIEGO_DEPLOYMENT_NAME
check_param SCALEIO_MDM_IPS
#install bosh cli (should add to docker image eventually...)
gem install bosh_cli --no-ri --no-rdoc

#Setup BOSH CLI for us
bosh -n target ${BOSH_DIRECTOR_PUBLIC_IP}
bosh -n login ${BOSH_USER} ${BOSH_PASSWORD}



manifest="$(bosh download manifest ${DIEGO_DEPLOYMENT_NAME} | perl -0pe 's/(- instances: 2\n  name: CI_cell_z1)(.+?)(- instances:)/- instances:/sg')"
echo "${manifest}" > manifest.yml

bosh deployment manifest.yml
bosh -n deploy

pushd rexray-boshrelease
  rexray_release_version=$(get_release_version ci-rexray-boshrelease)
  bosh -n delete release ci-rexray-boshrelease || true
  bosh -n create release --force --name ci-rexray-boshrelease --version ${rexray_release_version}
  bosh -n upload release
popd

pushd scaleio-sdc-boshrelease
  sdc_release_version=$(get_release_version ci-scaleio-sdc-boshrelease)
  bosh -n delete release ci-scaleio-sdc-boshrelease || true
  bosh -n create release --force --name ci-scaleio-sdc-boshrelease --version ${sdc_release_version}
  bosh -n upload release
popd

manifest="$(bosh download manifest ${DIEGO_DEPLOYMENT_NAME} | sed -e $'s/jobs:/jobs:\\\n- instances: 2\\\n  name: CI_cell_z1\\\n  networks:\\\n  - name: private\\\n    static_ips: ['"${DIEGO_CELL_IPS}"$']\\\n  properties:\\\n    scaleio:\\\n      mdm:\\\n        ips: ['"${SCALEIO_MDM_IPS}"$']\\\n    diego:\\\n      rep:\\\n        zone: z1\\\n    metron_agent:\\\n      zone: z1\\\n  vm_type: x-large\\\n  stemcell: trusty-3215\\\n  azs:\\\n  - z1\\\n  templates:\\\n  - name: consul_agent\\\n    release: cf\\\n  - name: rep\\\n    release: diego-release\\\n  - name: garden\\\n    release: garden-linux\\\n  - name: cflinuxfs2-rootfs-setup\\\n    release: cflinuxfs2-rootfs\\\n  - name: metron_agent\\\n    release: cf\\\n  - name: rexray_service\\\n    release: ci-rexray-boshrelease\\\n  - name: setup_sdc\\\n    release: ci-scaleio-sdc-boshrelease\\\n  update:\\\n    max_in_flight: 1\\\n    serial: false/g')"
echo "${manifest}" | perl -0pe "s/- name: ci-rexray-boshrelease\n  version: [0-9a-z\+\.]+/- name: ci-rexray-boshrelease\n  version: ${rexray_release_version}/sg" > manifest.yml
manifest="$(cat manifest.yml)"
echo "${manifest}" | perl -0pe "s/- name: ci-scaleio-sdc-boshrelease\n  version: [0-9a-z\+\.]+/- name: ci-scaleio-sdc-boshrelease\n  version: ${sdc_release_version}/sg" > manifest.yml

bosh -n deploy
echo "Diego has now evolved! RAWRRRR~"
