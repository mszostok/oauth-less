#!/usr/bin/env bash

readonly ROOT_PATH=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

pushd ${ROOT_PATH}/.. > /dev/null

# Exit handler. This function is called anytime an EXIT signal is received.
# This function should never be explicitly called.
function _trap_exit () {
    popd > /dev/null
}
trap _trap_exit EXIT

docker build -t oauth-less:latest .
docker tag oauth-less:latest mszostok/oauth-less:latest
docker push mszostok/oauth-less:latest
