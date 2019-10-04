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
    info "Generating mock clients (from config.go)"
    cd "$SOURCE_DIR"azuredevops
    mockgen -source=config.go -destination=mock_config.go -package=azuredevops
}

function generate_mocks() {
    check_gomock
    generate_mock_clients
    info "Mocks generated successfully"
}

generate_mocks
