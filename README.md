# CF-Persistence Service Broker [![Documentation Status](https://readthedocs.org/projects/cf-persist-service-broker/badge/?version=latest)](http://cf-persist-service-broker.readthedocs.io/en/latest/?badge=latest) 
This service broker provides the necessarily binding to orchestrate Cloud Controller and Diego Volume Manager. By using this service broker, Cloud Foundry applications can gain access to persistence service as if using a local filesystem. The orchestration uses [RexRay](https://github.com/emccode/rexray) and [libstorage](https://github.com/emccode/libstorage) internally. In theory, the service broker should be able to orchestrate all types of persistence provided by RexRay when configured correctly. We have tested orchestrating EMC ScaleIO and are currently in process of testing EMC Isilon and VMAX. Please check out our [documentations](http://cf-persist-service-broker.readthedocs.io/en/latest/)  and [blog](http://dojoblog.emc.com) for more updates.

### Setup
When starting the Service Broker you need to set a username and password. This is done by setting the environment variables `BROKER_USERNAME` and `BROKER_PASSWORD`.

This is an example of how we do this on CF.

 ```
 cf set-env persist BROKER_USERNAME broker
 cf set-env persist BROKER_PASSWORD broker
 ```

### Configuration

CF Persistence Service broker requires configuration to enable communication with [libstorage](https://github.com/emccode/libstorage). It can be configured using environment variables or a configuration file. Specify a configuration file using `-config path/to/file.yml`.  Here is an example configuration file:

```
libstorage:
  host: tcp://file_fake_host:9000
  storage:
    driver: scaleio

scaleio:
  endpoint:             https://file_fake_endpoint/api
  insecure:             true
  useCerts:             false
  userName:             file_fake_user
  password:             file_fake_password
  systemID:             file_fake_sys_id
  systemName:           file_fake_sys_name
  protectionDomainID:   file_fake_protection_domain_id
  protectionDomainName: file_fake_protection_domain_name
  storagePoolName:      file_fake_storage_pool_name
  thinOrThick:          ThinProvisioned
  version:              2.0
```

Alternatively, you can specify values using environment variables.  The variables should be named as follows:
```
LIBSTORAGE_HOST=tcp://file_fake_host:9000
LIBSTORAGE_STORAGE_DRIVER=scaleio
SCALEIO_ENDPOINT=https://file_fake_endpoint/api
SCALEIO_INSECURE=true
SCALEIO_USE_CERTS=false
SCALEIO_USERNAME=username
SCALEIO_PASSWORD=password
SCALEIO_SYSTEM_ID=systemId
SCALEIO_SYSTEM_NAME=systemName
SCALEIO_PROTECTION_DOMAIN_ID=scaleIOProtectionDomainID
SCALEIO_PROTECTION_DOMAIN_NAME=scaleIOProtectionDomainName
SCALEIO_STORAGE_POOL_NAME=scaleIOStoragePoolName
SCALEIO_THIN_OR_THICK=scaleIOThinOrThick
SCALEIO_VERSION=2.0
```


### About
More documentation coming soon.

### Contact
- Slack Channel:
  - Organization: <http://cloudfoundry.slack.com>
  - Channel: `#persi`
- Contact: [EMC Dojo](mailto:emcdojo@emc.com) [@EMCDojo](https://twitter.com/hashtag/emcdojo)
- Blog: [EMC Dojo Blog](dojoblog.emc.com)
