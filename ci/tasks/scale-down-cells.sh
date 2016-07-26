#!/usr/bin/env bash

set -e -x
source cf-persist-service-broker/ci/tasks/util.sh

check_param DIEGO_DEPLOYMENT_NAME
check_param CI_DIEGOCELL_IPS
check_param SCALEIO_MDM_IPS
check_param BOSH_DIR
check_param BOSH_USER
check_param BOSH_PASS

#install bosh cli (should add to docker image eventually...)
gem install bosh_cli --no-ri --no-rdoc

#Setup BOSH CLI for us
bosh -n target $BOSH_DIR
bosh -n login $BOSH_USER $BOSH_PASS

boshDiegoManifest="$(bosh download manifest $DIEGO_DEPLOYMENT_NAME)"
CI_BoshDiegoManifest="$(echo -e "${boshDiegoManifest}" | sed -e $'s/jobs:/jobs:\\\n- instances: 2\\\n  name: CI_cell_z1\\\n  networks:\\\n  - name: private\\\n    static_ips: ['"$CI_DIEGOCELL_IPS"$']\\\n  properties:\\\n    scaleio:\\\n      mdm:\\\n        ips: ['"$SCALEIO_MDM_IPS"']\\\n    diego:\\\n      rep:\\\n        zone: z1\\\n    metron_agent:\\\n      zone: z1\\\n  vm_type: x-large\\\n  stemcell: trusty-3215\\\n  azs:\\\n  - z1\\\n  templates:\\\n  - name: consul_agent\\\n    release: cf\\\n  - name: rep\\\n    release: diego-release\\\n  - name: garden\\\n    release: garden-linux\\\n  - name: cflinuxfs2-rootfs-setup\\\n    release: cflinuxfs2-rootfs\\\n  - name: metron_agent\\\n    release: cf\\\n  - name: rexray_service\\\n    release: rexray-bosh-release\\\n  - name: setup_sdc\\\n    release: scaleio-sdc-bosh-release\\\n  update:\\\n    max_in_flight: 1\\\n    serial: false/g')"
echo "${boshDiegoManifest}" > baseBOSHDiegoManifest.yml
echo "${CI_BoshDiegoManifest}" > CI_BOSHDiegoManigest.yml

bosh deployment CI_BOSHDiegoManigest.yml
bosh -n deploy

echo "Diego has now evolved! RAWRRRR~"
