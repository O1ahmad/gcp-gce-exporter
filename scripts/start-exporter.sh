#!/bin/bash

if [[ -n $GOOGLE_APPLICATION_CREDENTIALS_JSON ]]; then
    source $HOME/.bash_profile 
fi

gcp-gce-exporter
