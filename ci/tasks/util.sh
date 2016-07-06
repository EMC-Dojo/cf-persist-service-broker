#!/usr/bin/env bash

check_param() {
  local name=$1
  local value=$(eval echo '$'$name)
  if [ "$value" == 'replace-me' ]; then
    echo "environment variable $name must be set"
    exit 1
  fi
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

get_cf_services() {
  local output=("yes")
  ct=0
  cf marketplace -s EMC-Persistence |
  # augment marketplace output to get service names
  while read line;
    do ((ct+=1));
    if ((ct>4))
      then echo $line | awk '{print $1;}';
    fi
  done
}
