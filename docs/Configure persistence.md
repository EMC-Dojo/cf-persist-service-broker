# Allow Cloud Foundry app to use EMC persistent storage


## 1. Prerequisite
- Cloud Foundry Installation running with Diego. For more information about configuring CF or Diego, please see [https://docs.cloudfoundry.org/](https://docs.cloudfoundry.org/)


## 2. Install Rexray and ScaleIO SDC in Diego Cell
#### What is Rexray?

Rexray combines all of EMC storage solutions under one single interface. To read more about rexray, please go to the docs at [Rexray](http://rexray.readthedocs.io/en/v0.4.0-docs/).

 
#### What is ScaleIO SDC?

ScaleIO SDC is a component of ScaleIO that allow a VM to map to a ScaleIO volume. After, users can mount this volume to their system and use it.


#### Configure ScaleIO SDC and RexRay in Diego Manifest

In order to use ScaleIO with Cloud Foundry, we need ScaleIO SDC and Rexray to live on all of the Diego cells. In the diego manifest, [SDC Bosh Release](https://github.com/EMC-Dojo/ScaleIO-SDC-Bosh-Release) and [RexRay Bosh Release](https://github.com/EMC-Dojo/rexray-boshrelease) jobs should be collocated with the Diego jobs like below:


```
- instances: 1
  name: cell
  networks:
  - name: private
  properties:
    diego:
      rep:
        zone: z1
    metron_agent:
      zone: z1
  resource_pool: cell_z1
  templates:
  - name: setup_sdc
    release: scaleio-sdc-bosh-release
  - name: rexray_service
    release: rexray-bosh-release
  - name: consul_agent
    release: cf-release
  - name: rep
    release: diego-release
  - name: garden
    release: garden-linux
  - name: cflinuxfs2-rootfs-setup
    release: cflinuxfs2-rootfs
  - name: metron_agent
    release: cf-release
  update:
    max_in_flight: 1
    serial: false
```

## 3. Create CF ScaleIO Service
#####1. Deploy CF Persist Service Broker

Clone EMC CF Persist Service Broker to your workspace.
   
```
git clone https://github.com/EMC-Dojo/cf-persist-service-broker.git
```

Host the service broker to Cloud Foundry or third-party hosting. Below is the commands to push the service broker to Cloud Foundry

```
cf push #{service_broker_app_name} --no-start
cf set-env #{service_broker_app_name} BROKER_PASSWORD #{broker_password}
cf set-env #{service_broker_app_name} BROKER_USERNAME #{broker_username}
cf set-env #{service_broker_app_name} LIBSTORAGE_HOST #{libstorage_host}
cf set-env #{service_broker_app_name} LIBSTORAGE_STORAGE_DRIVER #{libstorage_storage_driver}
cf set-env #{service_broker_app_name} SCALEIO_ENDPOINT #{scaleio_endpoint}
cf set-env #{service_broker_app_name} SCALEIO_INSECURE #{scaleio_insecure}
cf set-env #{service_broker_app_name} SCALEIO_PASSWORD #{scaleio_password}
cf set-env #{service_broker_app_name} SCALEIO_PROTECTION_DOMAIN_ID #{scaleio_protection_domain_id}
cf set-env #{service_broker_app_name} SCALEIO_PROTECTION_DOMAIN_NAME #{scaleio_protection_domain_name}
cf set-env #{service_broker_app_name} SCALEIO_STORAGE_POOL_NAME #{scaleio_storage_pool_name}
cf set-env #{service_broker_app_name} SCALEIO_SYSTEM_ID #{scaleio_system_id}
cf set-env #{service_broker_app_name} SCALEIO_SYSTEM_NAME #{scaleio_system_name}
cf set-env #{service_broker_app_name} SCALEIO_THIN_OR_THICK #{scaleio_thin_or_thick}
cf set-env #{service_broker_app_name} SCALEIO_USE_CERTS #{scaleio_use_certs}
cf set-env #{service_broker_app_name} SCALEIO_USERNAME #{scaleio_username}
cf set-env #{service_broker_app_name} SCALEIO_VERSION #{scaleio_version}
cf start #{service_broker_app_name}
```

#####2. Register Service Broker

After hosting the app, we can register it as a service broker in Cloud Foundry. To register, run the following commands:

```
cf create-service-broker scaleiogo #{custom_broker_username} #{custom_broker_password} #{service_broker_app_url} cf enable-service-access scaleiogo
```

You would be able to see it in `cf marketplace`

#####3. Create Service Instance

Now you can create a scaleio service instance

```
cf create-service scaleiogo small #{custom_instance_name} -c '{"storage_pool_name": "#{your_storage_pool_name}"}'
```
	
#####4. Bind Service and Unbind Service to Your App
Push your app requiring persitent to Cloud Foundry. Let's call it `persistence_app`. Then bind it to scaleio_service_instance

```
cf bind-service persistence_app #{custom_instance_name}
cf restage persistence_app
```

If you want to unbind the service, run the following

```
cf unbind-service persistence_app #{custom_instance_name}
cf restage persistence_app
```

##### Important Note

```
The bind-service and unbind-service will not map/unmap or mount/unmount volume until you do an app restage command
```
	
	
		

	
