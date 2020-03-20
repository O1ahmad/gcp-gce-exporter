package main

import (
	"net/http"

	c "github.com/0x0I/gcp-gce-exporter/src/config"
	log "github.com/Sirupsen/logrus"
    "golang.org/x/net/context"
    "golang.org/x/oauth2/google"
    "google.golang.org/api/compute/v1"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	namespace = "gcp_gce" // Used to prepand Prometheus metrics
)

// Runtime variables
var (
	metricsPath = c.GetEnv("METRICS_PATH", "/metrics") // Path under which to expose metrics
	listenPort  = c.GetEnv("LISTEN_PORT", ":9687")     // Port on which to expose metrics
	logLevel    = c.GetEnv("LOG_LEVEL", "info")

	region      = c.GetEnv("REGION", "us-east1")
    project     = c.GetEnv("PROJECT", "")

    // GCP public image projects with and without Shielded VM support
    imageProjects []string = []string{
        "gce-uefi-images",
        "centos-cloud",
        "cos-cloud",
        "coreos-cloud",
        "debian-cloud",
        "rhel-cloud",
        "rhel-sap-cloud",
        "suse-cloud",
        "suse-sap-cloud",
        "ubuntu-os-cloud",
        "windows-cloud",
        "windows-sql-cloud",
    }

    // Google API access global variables for single context
    ctx = context.Background()
    computeClient *http.Client
    computeService *compute.Service
)

func main() {
	c.CheckConfig()

    var err error
    computeClient, err = google.DefaultClient(ctx, compute.CloudPlatformScope)
    if err != nil {
        log.Fatal(err)
    }
    computeService, err = compute.New(computeClient)
    if err != nil {
        log.Fatal(err)
    }

	setLogLevel(logLevel)
	log.Info("Starting Prometheus GCP GCE Exporter")

	exporter := newExporter()
	prometheus.MustRegister(exporter)

	http.Handle(metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
      <head><title>GCP GCE Exporter</title></head>
      <body>
      <h1>GCP GCE Exporter</h1>
      <p><a href=` + metricsPath + `>Metrics</a></p>
      </body>
      </html>`))
	})

	log.Printf("Starting Server on port %s and path %s", listenPort, metricsPath)
	log.Fatal(http.ListenAndServe(listenPort, nil))
}
