#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$DIR" || exit 1

VERSION=`cat version`

gb build \
  -ldflags "-X 'gosaic/environment.Version=$VERSION' -X 'gosaic/environment.BuildTime=$(date)' -X 'gosaic/environment.BuildUser=$(whoami)' -X 'gosaic/environment.BuildHash=$(git rev-parse HEAD)'" \
  all
