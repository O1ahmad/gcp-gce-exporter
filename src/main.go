package main

import (
	"net/http"

	c "github.com/0x0I/gcp-gce-exporter/src/config"
	log "github.com/Sirupsen/logrus"

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
	zone        = c.GetEnv("ZONE", "us-east1-a")
    project     = c.GetEnv("PROJECT", "")
)

func main() {
	c.CheckConfig()

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
