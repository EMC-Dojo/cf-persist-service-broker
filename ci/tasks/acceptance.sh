#!/usr/bin/env bash

set -e -x
source cf-persist-service-broker/ci/tasks/util.sh

check_param RACK_ENV
check_param CF_IP
check_param CF_ENDPOINT
check_param CF_USERNAME
check_param CF_PASSWORD
check_param CF_ORG
check_param CF_SPACE
check_param CF_SCALEIO_SERVICE_BROKER

uploaded_data='IAMADATA'
app='scaleio-acceptance-app'
service='scaleio-acceptance-service'
download_url=https://${app}.${CF_ENDPOINT}/data

pushd scaleio-acceptance-app
  cf api https://api.${CF_ENDPOINT} --skip-ssl-validation
  cf auth ${CF_USERNAME} ${CF_PASSWORD}
  cf target -o ${CF_ORG} -s ${CF_SPACE}

  cf push ${app} --no-start
  cf env ${app}
  cf create-service scaleiogo small ${service} -c '{"storage_pool_name": "default"}'
  cf bind-service ${app} ${service}
  cf start ${app}

  curl --insecure -X POST https://${app}.${CF_ENDPOINT}/data -d "${uploaded_data}" -H "Content-Type: text/plain"
  check_persistent ${uploaded_data} ${download_url}

  cf stop ${app}
  cf start ${app}
  check_persistent ${uploaded_data} ${download_url}

  cf restage ${app}
  check_persistent ${uploaded_data} ${download_url}

  cf unbind-service ${app} ${service}
  cf delete ${app} -f
  cf push ${app} --no-start
  cf bind-service ${app} ${service}
  cf start ${app}
  check_persistent ${uploaded_data} ${download_url}

  cf unbind-service ${app} ${service}
  cf delete-service ${service} -f
  cf delete-service-broker ${CF_SCALEIO_SERVICE_BROKER} -f
popd
