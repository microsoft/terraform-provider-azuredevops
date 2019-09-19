#!/usr/bin/env bash

set -euo pipefail

. ./commons.sh

function clean() {
    rm -rf "$BUILD_DIR"
    mkdir -p "$BUILD_DIR"
}

function test() {
    info "Executing unit tests"
    (
        cd "$SOURCE_DIR"
        go test ./... || fatal "Build finished in error due to failed unit tests"
    )
}

function compile() {
    NAME=$(cat $PROVIDER_NAME_FILE)
    VERSION=$(cat $PROVIDER_VERSION_FILE)

    BUILD_ARTIFACT="terraform-provider-${NAME}_v${VERSION}"

    info "Attempting to build $BUILD_ARTIFACT"
    (
        ROOT=$(pwd)
        cd "$SOURCE_DIR"
	go mod download 
        go build -o "$ROOT/$BUILD_DIR/$BUILD_ARTIFACT"
    )
}

function clean_and_build() {
    clean
    # test - enable in https://github.com/nmiodice/terraform-azure-devops-hack/issues/31
    compile
    info "Build finished successfully"
}

clean_and_build
