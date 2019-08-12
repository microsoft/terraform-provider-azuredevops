#!/usr/bin/env bash

set -euo pipefail

docker build -f Dockerfile . -t azdotf
