#!/usr/bin/env bash

set -e -x
source cf-persist-service-broker/ci/tasks/util.sh

check_param GITHUB_USER
check_param GITHUB_EMAIL
check_param REPO_NAME

# Creates an integer version number from the semantic version format
# May be changed when we decide to fully use semantic versions for releases
export integer_version=`cut -d "." -f1 version-semver/number`
cp -r $REPO_NAME promote/${REPO_NAME}
echo ${integer_version} > promote/integer_version

pushd promote/${REPO_NAME}
  git config --global user.email ${GITHUB_EMAIL}
  git config --global user.name ${GITHUB_USER}
  git config --global push.default simple

  export annotate_message=":airplane: New final release v${integer_version}"
  export log_message="`git log -1 --pretty=oneline`"

  echo "## v${integer_version}" >> CHANGELOG.md
  echo "v${log_message}" >> CHANGELOG.md

  git add CHANGELOG.md
  git commit -m "[ci skip] ${annotate_message}"
popd
