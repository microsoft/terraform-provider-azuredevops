#!/usr/bin/env bash

set -euo pipefail

. $(dirname $0)/commons.sh

info "Executing unit tests"
(
    cd "$SOURCE_DIR"
    go test -tags "${*:-all}" -v $(go list ./... | grep -v acceptancetests) || fatal "Build finished in error due to failed tests"
)
