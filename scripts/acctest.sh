#!/usr/bin/env bash

set -euo pipefail

. $(dirname $0)/commons.sh

info "Executing acceptance tests"
(
    cd "$SOURCE_DIR"

    # This is similar to the unit test command aside from the following:
    #   - TF_ACC=1 is a flag that will enable the acceptance tests. This flag is
    #     documented here:
    #       https://www.terraform.io/docs/extend/testing/acceptance-tests/index.html#running-acceptance-tests
    #
    #   - A `-run` parameter is used to target *only* tests starting with `TestAcc`. This prefix is
    #     recommended by Hashicorp and is documented here:
    #       https://www.terraform.io/docs/extend/testing/acceptance-tests/index.html#test-files
    TF_ACC=1 go test -timeout 120m -run ^TestAcc -tags "${*:-all}" -v $(go list ./... | grep acceptancetests) || fatal "Build finished in error due to failed tests"
)
