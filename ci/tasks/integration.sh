#!/usr/bin/env bash

set -e -x
source cf-persist-service-broker/ci/tasks/util.sh

check_param BROKER_PASSWORD
check_param BROKER_USERNAME
check_param DIEGO_DRIVER_SPEC
check_param INSECURE
check_param LIB_STOR_SERVICE
check_param LIBSTORAGE_URI
check_param PARSED_INSTANCE_ID
check_param PORT
check_param STORAGE_POOL_NAME
check_param TEST_INSTANCE_ID

export GOPATH=$PWD/gocode
export PATH=$PATH:$GOPATH/bin

go get github.com/onsi/ginkgo/ginkgo

mkdir -p gocode/src/github.com/EMC-Dojo
cp -r cf-persist-service-broker gocode/src/github.com/EMC-Dojo

cat > config.json <<EOF
[
  {
    "id": "67cfe587-0ec4-41c8-a7f7-97866d7f2d40",
    "name": "Persistent Storage",
    "description": "Supports EMC ScaleIO & Isilon Storage Arrays for use with CloudFoundry",
    "bindable": true,
    "requires": [
      "volume_mount"
    ],
    "plans": [
        {
          "name": "isilonservice",
          "description": "An isilon service",
          "metadata": {
            "bullets": [
              "Brings you isilon service"
            ],
            "displayName": "isilon"
          }
        },
        {
          "name": "scaleioservice",
          "description": "A scaleio service",
          "metadata": {
            "bullets": [
              "Brings you scaleio service"
            ],
            "displayName": "scaleio"
          }
        }
      ],
     "metadata": {
       "displayName": "Persistent Storage",
       "imageUrl": "imageURL",
       "longDescription": "Dell EMC brings you persistent storage on CloudFoundry",
       "providerDisplayName": "Dell EMC",
       "documentationUrl": "docsURL",
       "supportUrl": "supportURL"
     }
  }
]
EOF

export BROKER_CONFIG_PATH=$PWD/config.json
pushd gocode/src/github.com/EMC-Dojo/cf-persist-service-broker
  ginkgo -r
popd
