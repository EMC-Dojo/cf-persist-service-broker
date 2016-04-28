#!/usr/bin/env bash

set -e -x

source cf-persist-service-broker/ci/tasks/util.sh

# Creates an integer version number from the semantic version format
# May be changed when we decide to fully use semantic versions for releases
export integer_version=`cut -d "." -f1 version-semver/version`
mkdir promote
echo ${integer_version} > promote/integer_version
echo ":airplane: New release v${integer_version}" > promote/tag_message
