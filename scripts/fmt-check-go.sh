#!/usr/bin/env bash

set -euo pipefail

. $(dirname $0)/commons.sh

info "Formatting Go Files... If this fails, run 'find . -name \"*.go\" | xargs gofmt -s -w' to fix"
(
    cd "$SOURCE_DIR"

    # This runs a go fmt on each file without using the 'go fmt ./...' syntax.
    # This is advantageous because it avoids having to download all of the go
    # dependencies that would have been triggered by using the './...' syntax.
    FILES_WITH_FMT_ISSUES=$(find . -name "*.go" | xargs gofmt -s -l | wc -l)

    # convert to integer...
    FILES_WITH_FMT_ISSUES=$(($FILES_WITH_FMT_ISSUES + 0))

    info "Found $FILES_WITH_FMT_ISSUES with formatting issues..."
    exit $FILES_WITH_FMT_ISSUES
)
