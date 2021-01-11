#!/usr/bin/env bash

set -euo pipefail

. $(dirname $0)/commons.sh

SCRIPTS_DIR="$(cd -P "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BUILD_DIR="${SCRIPTS_DIR}/../dist/"
SOURCE_DIR="${SCRIPTS_DIR}/../"
NAME="azuredevops"
BUILD_ARTIFACT="terraform-provider-${NAME}_v${VERSION}"
ARCHIVE_ARTIFACT="terraform-provider-${NAME}_${VERSION}"

OS_ARCH=("freebsd:amd64"
  "freebsd:386"
  "freebsd:arm"
  "freebsd:arm64"
  "windows:amd64"
  "windows:386"
  "linux:amd64"
  "linux:386"
  "linux:arm"
  "linux:arm64"
  "darwin:amd64")


function clean() {
  info "Cleaning $BUILD_DIR"
  rm -rf "$BUILD_DIR"
  mkdir -p "$BUILD_DIR"
}

function release() {
  info "Clean build directory"
  clean

  info "Attempting to build ${BUILD_ARTIFACT}"

  cd "$SOURCE_DIR"
  go mod download
  for os_arch in "${OS_ARCH[@]}" ; do
    OS=${os_arch%%:*}
    ARCH=${os_arch#*:}
    info "GOOS: ${OS}, GOARCH: ${ARCH}"
    (
      env GOOS="${OS}" GOARCH="${ARCH}" CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -X 'github.com/microsoft/terraform-provider-azuredevops/version.ProviderVersion=v${VERSION}'" -o "${BUILD_ARTIFACT}"|
      zip "${ARCHIVE_ARTIFACT}_${OS}_${ARCH}.zip" "${BUILD_ARTIFACT}"
      rm -rf "${BUILD_ARTIFACT}"
    )
  done
  mv *.zip "${BUILD_DIR}"
  cd "${BUILD_DIR}"
  shasum -a 256 *.zip > "${ARCHIVE_ARTIFACT}_SHA256SUMS"
  cp "${ARCHIVE_ARTIFACT}_SHA256SUMS" "${ARCHIVE_ARTIFACT}_SHA256SUMS.sig"
  cat "${ARCHIVE_ARTIFACT}_SHA256SUMS"
  cp ../scripts/dearmor.sh ./
}

release
