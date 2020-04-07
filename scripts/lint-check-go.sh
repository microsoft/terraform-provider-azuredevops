#!/usr/bin/env bash

set -euo pipefail

. $(dirname $0)/commons.sh

info "Linting Go Files... If this fails, run 'golint ./... | grep -v 'vendor' ' to see errors"
(
    cd "$SOURCE_DIR"

    # Install golint if not there
    go get -u golang.org/x/lint/golint 2>/dev/null

    GOLINT="$(go list -f {{.Target}} golang.org/x/lint/golint)"
    "$GOLINT" -set_exit_status $(go list ./... | grep -v 'vendor')
)
