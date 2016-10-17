#Enabling Isilon in Open Source CF

_This Guide assumes that diego cells and CF BOSH deployments are separate manifests. If this is not the case; combine steps 1, 3, and 4 accordingly to put the proper YAML together in the same manifest._

_Take note that this guide enables the use of the service broker, but does not implement it. For more information on installing the service broker either through the avaiable BOSH releases for [libstorage](https://github.com/EMC-CMD/libstorage-release), [our persistent broker](https://github.com/EMC-CMD/emc-persistence-release), and other options; contact us through victor.fong@emc.com._


####Step 1: Enable Volume Services in CF Manifest

Locate your CF Deployment manifest and add the following property:

```
properties:
  ...
  cc:
    volume_services_enabled: true
```

`volume_services_enabled` allows the experimental "volume services" feature of CF to connect to our service broker.

####Step 2: Create RexRay Release

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

####Step 3: Add Bosh Volume Driver (RexRay)

Edit the Diego manifest with your favorite editor (VIM), and modify the following YAML structures:

```
releases:
- name: rexray-bosh-release  
  version: 'latest'  
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
          volumePath: /rexray  
          nfsHost: #EDIT ME
          dataSubnet: #EDIT ME
          quotas: #EDIT ME
          sharedMounts: true
        linux:  
          volume:  
            fileMode: 0777  
```
__NOTE: VALUES ABOVE MARKED "EDIT ME" WILL NOT WORK, UPDATE THEM__

This should mount Rexray into all Diego Cells on the next deploy. It's probably a good idea to make a copy of this manifest for reference. A re-deployment of Pivotal Elastic Runtime may cause a wipe of these changes, and the valuable fields will have to be recreated.

_The volume path specified above must exist in the Isilon cluster. (i.e. our example would be /ifs/volumes/rexray/)_

####Step 4: Re-Deploy Diego Cells

```
bosh deployment my-cf-manifest.yml
bosh deploy
...
bosh deployment my-diego-cells.yml
bosh deploy
```

If this succeeds, congratulations! Rexray should be running on Diego Cells, and allowing you to use Libstorage for a Persistent Storage solution. If you need help deploying libstorage, check out our [libstorage bosh release](https://github.com/EMC-CMD/libstorage-release)
