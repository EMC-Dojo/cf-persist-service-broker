#!/usr/bin/env bash

set -e -x

source cf-persist-service-broker/ci/tasks/util.sh

# Creates an integer version number from the semantic version format
# May be changed when we decide to fully use semantic versions for releases
export integer_version=`cut -d "." -f1 version-semver/number`
cp -r cf-persist-service-broker promote/cf-persist-service-broker
echo ${integer_version} > promote/integer_version

pushd promote/cf-persist-service-broker/
  echo ${tag_message} >> release_log.txt
  git add release_log.txt
  git config --global user.email emccmd-eng@emc.com
  git config --global user.name EMCCMD-CI
  export annotate_message=":airplane: New final release v${integer_version}"
  git commit -m ${annotate_message} -m "[ci skip]"
popd
