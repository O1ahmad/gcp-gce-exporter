#!/bin/bash

set -eo pipefail
set +x

if [[ -n $GOOGLE_APPLICATION_CREDENTIALS_JSON ]]; then
    echo "${GOOGLE_APPLICATION_CREDENTIALS_JSON}" > /tmp/gcp_creds.json

    echo "export GOOGLE_APPLICATION_CREDENTIALS=/tmp/gcp_creds.json" >> $HOME/.bash_profile
else
    echo "GCLOUD credentials JSON not provided. GOOGLE_APPLICATION_CREDENTIALS setting override disabled."
fi
