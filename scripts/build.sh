#!/usr/bin/env bash

set -euo pipefail

. $(dirname $0)/commons.sh

skipTest=0
debugBuild=0
doInstall=0

while [[ "$#" -gt 0 ]]; do
  case $1 in
  -s | --SkipTests)
    skipTest=1
    ;;
  -d | --DebugBuild)
    debugBuild=1
    ;;
  -i | --Install)
    doInstall=1
    ;;
  esac
  shift
done

function clean() {
  info "Cleaning $BUILD_DIR"
  rm -rf "$BUILD_DIR"
  mkdir -p "$BUILD_DIR"
}

function compile() {
  NAME=$(cat $PROVIDER_NAME_FILE)
  VERSION=$(cat $PROVIDER_VERSION_FILE)

  BUILD_ARTIFACT="terraform-provider-${NAME}_v${VERSION}"

  info "Attempting to build $BUILD_ARTIFACT"
  (
    cd "$SOURCE_DIR"
    go mod download
    if [ $debugBuild -eq 1 ]; then
      info "Using debug build settings"
      go build -o "$BUILD_DIR/$BUILD_ARTIFACT" -gcflags=all="-N -l"
    else
      go build -o "$BUILD_DIR/$BUILD_ARTIFACT"
    fi
  )
}

function clean_and_build() {
  clean
  if [ $skipTest -ne 1 ]; then
    $(dirname $0)/unittest.sh
  fi
  compile $debugBuild
  info "Build finished successfully"
}

clean_and_build

if [ $doInstall -eq 1 ]; then
  $(dirname $0)/local-install.sh
fi
