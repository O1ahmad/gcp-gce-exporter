<p><img src="https://cdn.worldvectorlogo.com/logos/prometheus.svg" alt="Prometheus logo" title="prometheus" align="left" height="60" /></p>
<p><img src="https://cloud.google.com/images/social-icon-google-cloud-1200-630.png" alt="gcp logo" title="gcp" align="right" height="100" /></p>

# GCP GCE :cloud: Exporter
A prometheus exporter providing metrics for GCP GCE compute resource specifications and capacity profiling.

![GitHub release (latest by date)](https://img.shields.io/github/v/release/0x0I/gcp-gce-exporter?color=yellow)
[![Build Status](https://travis-ci.org/0x0I/gcp-gce-exporter.svg?branch=master)](https://travis-ci.org/0x0I/gcp-gce-exporter)
[![Docker Pulls](https://img.shields.io/docker/pulls/0labs/0x01.gcp-gce-exporter?style=flat)](https://hub.docker.com/repository/docker/0labs/0x01.gcp-gce-exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/0x0I/gcp-gce-exporter)](https://goreportcard.com/report/github.com/0x0I/gcp-gce-exporter)
[![License: MIT](https://img.shields.io/badge/License-MIT-blueviolet.svg)](https://opensource.org/licenses/MIT)

Exposes compute resource statistics of GCP GCE machine-types, images and regions from the GCP Compute Engine API to a Prometheus compatible endpoint.

## Description

The application can be run in a number of ways though its main consumption is via the Docker hub image `0x0I.gcp-gce-exporter`.

**Required**
* `PROJECT`                                   - Google Cloud Platform API authorized and capable Google Cloud project
* `GOOGLE_APPLICATION_CREDENTIALS`            - path to file containing service account key and authentication credentials (**requires file mount into container**)
* `GOOGLE_APPLICATION_CREDENTIALS_JSON`       - json blob containing service account key and authentication credentials (*to be used in place of* **GOOGLE_APPLICATION_CREDENTIALS**)

**Optional**
* `METRICS_PATH`           - Path under which to expose metrics. Defaults to `/metrics`
* `LISTEN_PORT`            - Port on which to expose metrics. Defaults to `9692`
* `REGION`                 - GCP region to scrape. Defaults to `us-east1`
* `LOG_LEVEL`              - Set the logging level. Defaults to `info`

## Install and deploy

Run manually from Docker Hub (**Note:** credentials file must exist within Container):
```
docker run --detach --env GOOGLE_APPLICATION_CREDENTIALS="/path/to/creds.json" \
           --env PROJECT="XXXXXXX" --publish 9692:9692 \
           --volume /home/user/gcp.creds.json:/path/to/creds.json \
           0Iabs/0x01.gcp-gce-exporter
```

Scrape non-default GCP GCE region and increase logging level:
```
docker run --detach --env GOOGLE_APPLICATION_CREDENTIALS="/path/to/creds.json" \
           --env PROJECT="XXXXXXX" \
           --env REGION=asia-southeast1 \
           --env LOG_LEVEL=debug \
           --publish 9692:9692 \
           0Iabs/0x01.gcp-gce-exporter
```

Build and run a local container image:
```
docker build --file build/Containerfile --tag <image-name> .
docker run -d -e GOOGLE_APPLICATION_CREDENTIALS="/path/to/creds.json" -e PROJECT="XXXXXXX" -p 9692:9692 <image-name>
```

## Docker compose

```
version: "3"

services:
 gcp-gce-exporter:
    image: 0labs/0x01.gcp-gce-exporter:0.1.0
    container_name: gcp-exporter
    environment:
      GOOGLE_APPLICATION_CREDENTIALS_JSON: <creds-json>
      PROJECT: "<insert-project>"
    ports:
      - 9692:9692/tcp
```

## Metrics

Metrics will be made available on port **9692** by default, or you can pass environment variable ```LISTEN_ADDRESS``` to override this. An example printout of the metrics you should expect to see can be found in [METRICS.md](https://github.com/0x0I/gcp-gce-exporter/blob/master/METRICS.md).
