#!/usr/bin/env bash

cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )"

GIT=`command -v git 2>/dev/null`
if [ -z "$GIT" ]; then
  >&2 echo "git not found in PATH"
  exit 1
fi

$GIT rev-parse --git-dir >/dev/null 2>&1
if [ "$?" -ne "0" ]; then
  >&2 echo "not a git directory"
  exit 1
fi

ORIGIN=${ORIGIN:-origin}
CURRENT=`$GIT rev-parse --abbrev-ref HEAD 2>/dev/null`

if [ -z "$CURRENT" ]; then
  >&2 echo "unable to get current branch"
  exit 1
fi

gitstatus=`$GIT status 2> /dev/null | tail -n1`

if [[ $gitstatus != *"working directory clean"* ]]; then
  >&2 echo "please clean local changes on working copy before proceeding"
  exit 1
fi

$GIT pull --ff-only $ORIGIN $CURRENT >/dev/null 2>&1
if [ "$?" -ne "0" ]; then
  >&2 echo "failed to fast-forward pull"
  exit 1
fi

version=`cat version`
if [ -z "$version" ]; then
  >&2 echo "unable to determine version"
  exit 1
fi

if $GIT config -l | grep -q 'user.signingkey='; then
  flag="-s"
else
  flag="-a"
fi

$GIT tag $flag -m "Tag for release version $version" $version
if [ "$?" -ne "0" ]; then
  >&2 echo "failed to create git tag"
  exit 1
fi

$GIT push $ORIGIN $version
