#!/usr/bin/env bash

set -euo pipefail


BUILD_DIR="./bin"
SOURCE_DIR="./src"
PROVIDER_NAME_FILE="./PROVIDER_NAME.txt"
PROVIDER_VERSION_FILE="./PROVIDER_VERSION.txt"


function log() {
    LEVEL="$1"
    shift
    echo "[$LEVEL] $@"
}

function info() {
    log "INFO" $@
}

function fatal() {
    log "FATAL" $@
    exit 1
}