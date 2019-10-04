#!/usr/bin/env bash

set -euo pipefail

. $(dirname $0)/commons.sh

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
        go build -o "$BUILD_DIR/$BUILD_ARTIFACT"
    )
}

function clean_and_build() {
    clean
    $(dirname $0)/unittest.sh
    compile
    info "Build finished successfully"
}

clean_and_build
