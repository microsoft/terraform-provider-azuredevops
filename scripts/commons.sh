#!/usr/bin/env bash

set -euo pipefail

SCRIPTS_DIR="$(cd -P "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BUILD_DIR="${SCRIPTS_DIR}/../bin/"
SOURCE_DIR="${SCRIPTS_DIR}/../"
PROVIDER_NAME_FILE="${SCRIPTS_DIR}/../PROVIDER_NAME.txt"
PROVIDER_VERSION_FILE="${SCRIPTS_DIR}/../PROVIDER_VERSION.txt"

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
if [ "$OS" = "linux" ]; then
  PROC=$(lscpu 2> /dev/null | awk '/Architecture/ {if($2 == "x86_64") {print "amd64"; exit} else if($2 ~ /arm/) {print "arm"; exit} else if($2 ~ /aarch64/) {print "arm"; exit} else {print "386"; exit}}')
  if [ -z $PROC ]; then
    PROC=$(cat /proc/cpuinfo | awk '/model\ name/ {if($0 ~ /ARM/) {print "arm"; exit}}')
  fi
  if [ -z $PROC ]; then
    PROC=$(cat /proc/cpuinfo | awk '/flags/ {if($0 ~ /lm/) {print "amd64"; exit} else {print "386"; exit}}')
  fi
else
  PROC="amd64"
fi
[ "$(echo "$PROC" | grep 'arm')" != '' ] && PROC='arm'  # terraform downloads use "arm" not full arm type

function log() {
    LEVEL="$1"
    shift
    echo "[$LEVEL] $@"
}

function info() {
    log "INFO" $@
}

function fatal() {
    log "FATAL" $@
    exit 1
}
