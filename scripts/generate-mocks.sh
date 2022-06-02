#!/usr/bin/env bash

set -euo pipefail

. $(dirname $0)/commons.sh

MOCK_PKG_NAME="azdosdkmocks"


function install_gomock() {
    info "Installing GoMock tools"
    (
        go get github.com/golang/mock/mockgen
    )
}

function check_gomock() {
    info "Checking if GoMock tools are installed"
    (
        if ! [ -x "$(command -v mockgen)" ]; then
            install_gomock
        fi
    )
}

function generate_single_mock_client() {
    info "Generating mock for package: $1"

    # This is a nice trick to get the last element of a string that is delimited by the '/' character
    #   https://stackoverflow.com/questions/3162385/how-to-split-a-string-in-shell-and-get-the-last-field
    PACKAGE_NAME_SIMPLE=$(echo ${1##*/})

    # the generated file needs a unique name, so we can incorporate the package from which it is sourced
    OUTPUT_FILE="${PACKAGE_NAME_SIMPLE}_sdk_mock.go"

    # the prefix of the mock, used to give the generated mock a unique name
    MOCK_PREFIX="$(tr '[:lower:]' '[:upper:]' <<< ${PACKAGE_NAME_SIMPLE:0:1})${PACKAGE_NAME_SIMPLE:1}"
    MOCK_NAME="Mock${MOCK_PREFIX}${2}"

    OUTPUT_DIR="$SOURCE_DIR/$MOCK_PKG_NAME"
    mkdir -p "$OUTPUT_DIR"

    info "Generating mock client: $MOCK_NAME"
    mockgen \
        -package $MOCK_PKG_NAME \
        -destination "$OUTPUT_DIR/$OUTPUT_FILE" \
        -mock_names "$2=$MOCK_NAME" \
        --build_flags=--mod=mod \
        $1 $2

    if [ $? -eq 0 ]
    then
        info "$MOCK_NAME generated successfully"
    else
        info "$MOCK_NAME not generated successfully. This is expected if the package only contains API models and no API SDKs"
    fi
}

function generate_mock_clients() {
    info "Generating mock clients"

    info "Discovering all Azure DevOps Go SDK packages..."

    # Note: the `|| true` and related wrapping with `()` is here to prevent the script from failing
    #       in the case that the status code from `go list all` is non-zero. This can happen, for example,
    #       if the mocks do not exist at the time the script is initially run.
    AZDO_SDK_PACKAGES=$( (go list all || true) | grep 'azure-devops-go-api/azuredevops/')

    info "Found $(echo "$AZDO_SDK_PACKAGES" | wc -l) packages that may need mocking..."
    for PACKAGE in $AZDO_SDK_PACKAGES; do
        # note: all interfaces are currently named `Client`, but this may change so I'm
        #       leaving the parameter hard-coded here.
        generate_single_mock_client "$PACKAGE" "Client" || true
    done
}

function generate_mocks() {
    check_gomock
    generate_mock_clients
    info "Mocks generated successfully"
}

generate_mocks
