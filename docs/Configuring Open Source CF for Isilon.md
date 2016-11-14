#Enabling Isilon in Open Source CF

####Pre-Requisits
Before enabling Isilon on the Diego Cells with the following guide, ensure you have already performed the following 2 steps.

- Installed and have a running `LibStorage` host connected to your Isilon.

- Installed and have a running `Isilon Service Broker` connected to the above LibStorage host. This may be deployed inside of CloudFoundry as an application, externally through BOSH or stand-alone on a VM.

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
          logging:
            level: debug
          embedded: true  
          server:
            logging:
              level: debug  
            services:  
              isilon:  
                driver: isilon
          integration:
            volume:
              operations:
                mount:
                  path: /var/vcap/data  
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

####Step 4: Re-Deploy Diego Cells

```
bosh deployment my-cf-manifest.yml
bosh deploy
...
bosh deployment my-diego-cells.yml
bosh deploy
```

If this succeeds, congratulations! Libstorage should be running on a VM, Rexray should be running on Diego Cells, and the Isilon Service Broker should be deployed on opensource CloudFoundry. All of this to allow you to use your Isilon as a Persistent Storage solution. If you need any additional help contact us at any of the following!

- Slack Channel:
  - Organization: <http://cloudfoundry.slack.com>
  - Channel: `#persi`
- Contact: [EMC Dojo](mailto:emcdojo@emc.com) [@EMCDojo](https://twitter.com/hashtag/emcdojo)
- Blog: [EMC Dojo Blog](http://dojoblog.emc.com)
