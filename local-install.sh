#!/usr/bin/env bash

set -euo pipefail

. ./commons.sh


PLUGINS_DIR="$HOME/.terraform.d/plugins/"
mkdir -p "$PLUGINS_DIR"

info "Installing provider to $PLUGINS_DIR"
cp "$BUILD_DIR"/* "$PLUGINS_DIR/"