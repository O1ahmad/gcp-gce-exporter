package main

import "github.com/prometheus/client_golang/prometheus"

// AddMetrics creates gauge and counter metric vectors based on resource properties
func AddMetrics() (map[string]*prometheus.GaugeVec, map[string]*prometheus.CounterVec) {

	gaugeVecs := make(map[string]*prometheus.GaugeVec)
	counterVecs := make(map[string]*prometheus.CounterVec)

	// instance metrics
	gaugeVecs["totalvCPUs"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "total_vcpus",
			Help:      "Total virtual CPUs capacity provided by a machine-type",
		}, []string{"instance_type", "region", "zone"})
	gaugeVecs["totalMem"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "total_mem",
			Help:      "Total memory capacity(GiB) provided by a machine-type",
		}, []string{"instance_type", "region", "zone"})
	gaugeVecs["maxStorage"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "max_storage",
			Help:      "Maximum persistent disk storage capacity(GiB) provided by a machine-type",
		}, []string{"instance_type", "region", "zone"})
	gaugeVecs["maxDisks"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "max_disk",
			Help:      "Maximum persistent disks allowed by a machine-type",
		}, []string{"instance_type", "region", "zone"})

	// image metrics
	counterVecs["totalImages"] = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "total_images",
			Help:      "Total count of publically available images",
		}, []string{"project", "family", "src_img"})
	gaugeVecs["imgArchiveSize"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "image_archive_size",
			Help:      "Size of the image tar.gz archive stored in Google Cloud Storage (in bytes).",
		}, []string{"project", "name", "family", "src_img"})
	gaugeVecs["imgDiskSize"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "image_disk_size",
			Help:      "Size of the image when restored onto a persistent disk (in GB).",
		}, []string{"project", "name", "family", "src_img"})

	// disk type metrics
	gaugeVecs["diskTypeSize"] = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "disk_type_size",
			Help:      "Server-defined default disk size in GB",
		}, []string{"name", "valid_disk_size", "region", "zone"})

	counterVecs["totalDiskTypes"] = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "total_disk_types",
			Help:      "Total count of disk types per zone",
		}, []string{"valid_disk_size", "region", "zone"})

	// region metrics
	counterVecs["totalRegions"] = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "total_regions",
			Help:      "Total count of publically accessible GCP regions",
		}, []string{"name", "status"})

	// zone metrics
	counterVecs["totalZones"] = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "total_zones",
			Help:      "Total count of publically accessible GCP zones",
		}, []string{"name", "region", "status"})

	return gaugeVecs, counterVecs

}
