# CF-Persistence Service Broker

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
LIBSTORAGE_URI=tcp://file_fake_host:9000
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
  - Channel: `#general`
- Contact: [Peter Blum](mailto:peter.blum@emc.com) [@EMCDojo](https://twitter.com/hashtag/emcdojo)
- Blog: [EMC Dojo Blog](dojoblog.emc.com)
