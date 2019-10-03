#!/usr/bin/env bash

set -euo pipefail

. $(dirname $0)/commons.sh

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

function generate_mock_clients() {
    info "Generating mock clients (from azdo_config.go)"
    cd "$SOURCE_DIR"
    mockgen -source=azdo_config.go -destination=mock_azdo_client.go -package=main
}

function generate_mocks() {
    check_gomock
    generate_mock_clients
    info "Mocks generated successfully"
}

generate_mocks
