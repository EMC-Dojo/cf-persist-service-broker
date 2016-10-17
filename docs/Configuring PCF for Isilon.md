#Enabling Isilon in PCF

####Context: Why is this Hacky?

Persistence in PCF is only newly available as of version 1.8. When version 1.9 is released (Q4 2016), there will be cleaner avenues to interact with bosh manifests. For now, we must document the horrible, no good, very bad way.

####Step 1: Install Pivotal Elastic Runtime

This is a gigantic step, and is better left to [better docs](https://network.pivotal.io/products/elastic-runtime)

####Step 2: Install Isilon Tile

Install the Isilon Tile, available from the DellEMC Dojo. Contact victor.fong@dell.com for more information on this. (Link coming soon as we opensource it!)

####Step 3: SSH into OpsManager

Using the credentials created through the internal authentication step when OpsManager first booted, ssh into the machine. Ensure that git is installed on this machine (OpsMan 1.8+ seems to include git by default).

####Step 4: Create RexRay Release

Using git, clone the repository at https://github.com/EMC-Dojo/rexray-boshrelease.git like:
```
git clone https://github.com/EMC-Dojo/rexray-boshrelease.git
```
Ensure that you have logged into the bosh director for OpsManager. You can find the Director IP and credentials through the OpsManager UI.

To create and upload the release from the github repository, `cd` into the `rexray-boshrelease` directory, and use the following commands:
```
bosh create release --name rexray-bosh-release
bosh upload release $PWD/dev_releases/rexray-bosh-release/rexray-bosh-release-0+dev.1.yml
```

There should now be a bosh release like :

```
| rexray-bosh-release       | 0+dev.1       | w00th4sh    |
```

you can check that this is on the bosh director with `bosh releases`.

####Step 5: Find CF Manifest

Navigate to the manifest directory with `cd /var/tempest/workspaces/default/deployments/`, and match the deployment manifest with the current deployment of CF.

__Example:__ running `bosh deployments` may give the following block:
```
+-------------------------+-------------------------------+--------------------------------------------------+--------------+
| Name                    | Release(s)                    | Stemcell(s)                                      | Cloud Config |
+-------------------------+-------------------------------+--------------------------------------------------+--------------+
| cf-260f32b99ab848a8d984 | cf-autoscaling/36             | bosh-vsphere-esxi-ubuntu-trusty-go_agent/3262.12 | latest       |
|                         | cf-mysql/26.4                 |                                                  |              |
|                         | cf/239.0.17                   |                                                  |              |
|                         | cflinuxfs2-rootfs/1.26.0      |                                                  |              |
|                         | consul/108                    |                                                  |              |
|                         | diego/0.1485.0                |                                                  |              |
|                         | etcd/60                       |                                                  |              |
|                         | garden-linux/0.342.0          |                                                  |              |
|                         | mysql-backup/1.25.0           |                                                  |              |
|                         | mysql-monitoring/5            |                                                  |              |
|                         | notifications-ui/17           |                                                  |              |
|                         | notifications/24              |                                                  |              |
|                         | pivotal-account/1             |                                                  |              |
|                         | push-apps-manager-release/652 |                                                  |              |
|                         | routing/0.135.0               |                                                  |              |
|                         | service-backup/14             |                                                  |              |
+-------------------------+-------------------------------+--------------------------------------------------+--------------+
```

Given that the name here is `cf-260f32b99ab848a8d984`, the deployment manifest in the directory should be `cf-260f32b99ab848a8d984.yml`.

####Step 6: Edit Manifest

Edit the manifest with your favorite CLI editor (VIM), and modify the following YAML structures:

```
releases:
- name: rexray-bosh-release  
  version: 'latest'  
```
```
jobs:
...
- name: cloud_controller
  ...
  properties:
  ...
    cc:
       volume_services_enabled: true  
...
...
- name: diego_cell
  ...
  ...
  jobs:
  - name: rexray_service  
    release: rexray-bosh-release  
  ...
  properties:
    rexray: |  
        ---  
        rexray:  
          modules:  
            isilon:  
              disabled: false  
              host: tcp://127.0.0.1:9002  
              spec: /var/vcap/data/voldrivers/rexray_isilon.spec  
              http:  
                writetimeout: 900  
                readtimeout: 900  
              type: docker  
              libstorage:  
                service: isilon  
        libstorage:  
          embedded: true  
          server:  
            services:  
              isilon:  
                driver: isilon  
        isilon:  
          endpoint: #EDIT ME
          insecure: true  
          username: #EDIT ME
          password: #EDIT ME
          volumePath: /rexray #READ NOTE2 BELOW
          nfsHost: #EDIT ME
          dataSubnet: #EDIT ME
          quotas: #EDIT ME
          sharedMounts: true
        linux:  
          volume:  
            fileMode: 0777  
```
_NOTE: VALUES ABOVE MARKED "EDIT ME" WILL NOT WORK, UPDATE THEM AND REMOVE THE COMMENT INCLUDING #_

_NOTE2: volumePath refers to the path relative to /ifs/volumes/ on your Isilon. This is the folder where we will create subfolders for your CloudFoundry. ENSURE THIS IS PRE-CREATED AND IS THE CORRECT PATH #_

__

This should deploy the volume driver, Rexray, into all Diego Cells on the next deploy. It's probably a good idea to make a copy of this manifest for reference. A re-deployment of Pivotal Elastic Runtime may cause a wipe of these changes, and the valuable fields will have to be recreated.

####Step 7: Re-Deploy Elastic Runtime (Diego Cells)

Set the bosh deployment to your newly created manifest and deploy.

```
bosh deployment my-cf-manifest-copy.yml
bosh deploy
```

If this succeeds, congratulations! Libstorage should be running on a VM, Rexray should be running on Diego Cells, and the Isilon Service Broker should be deployed on PCF. All of this to allow you to use your Isilon as a Persistent Storage solution. If you need any additional help contact us at any of the following!

- Slack Channel:
  - Organization: <http://cloudfoundry.slack.com>
  - Channel: `#persi`
- Contact: [EMC Dojo](mailto:emcdojo@emc.com) [@EMCDojo](https://twitter.com/hashtag/emcdojo)
- Blog: [EMC Dojo Blog](http://dojoblog.emc.com)
