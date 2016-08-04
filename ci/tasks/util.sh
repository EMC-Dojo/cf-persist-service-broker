#!/usr/bin/env bash
set -e -x

check_param() {
  local name=$1
  local value=$(eval echo '$'$name)
  if [ "$value" == 'replace-me' ]; then
    echo "environment variable $name must be set"
    exit 1
  fi
}


get_release_version() {
  version=$(bosh releases | grep $1 | awk '{print $4}' | awk -F "*" '{print $1}')
  if [[ $version =~ ^[0-9]+$ ]]; then
    version=$((version+1))
  else
    version=1
  fi
  echo $version
}

check_persistent() {
  local uploaded_data=$1
  echo $2
  data=$(curl --insecure $2)

  if [ "${uploaded_data}" != "${data}" ]; then
    echo "data is not persist"
    exit 1
  fi
}

get_cf_service() {
  output=$(cf marketplace -s $EMC_SERVICE_NAME)
  if [ $? -eq 1 ]
  then
    exit 1
  fi

  ct=0
  echo "$output" |
  # augment marketplace output to get service names
  while read line;
    do ((ct+=1));
    if ((ct>4))
      then echo $line | awk '{print $1;}';
    fi
  done
}
