#!/bin/bash

set -euo pipefail

# Print all commands executed if DEBUG mode enabled
[ -n "${DEBUG:-""}" ] && set -x

# [Test-Setup]
docker build --file build/Containerfile --tag gcp-gce-exporter:testing .

# [Test-Run+Validate]
GOSS_FILES_PATH=test dgoss run gcp-gce-exporter:testing
