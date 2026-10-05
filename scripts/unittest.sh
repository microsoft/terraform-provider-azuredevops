#!/usr/bin/env bash

set -euo pipefail

. $(dirname $0)/commons.sh

if [ "$#" -gt 0 ]; then
    fatal "Selecting tests by build tag is no longer supported. Use 'go test -run <pattern>' instead."
fi

info "Executing unit tests"
(
    cd "$SOURCE_DIR"
    go test -v $(go list ./... | grep -v acceptancetests) || fatal "Build finished in error due to failed tests"
)
