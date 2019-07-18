#!/bin/bash

set -e

pushd "$(dirname "$0")/../docker/base" > /dev/null
./build.sh
popd > /dev/null

pushd "$(dirname "$0")/../docker/app" > /dev/null
./build.sh
popd > /dev/null