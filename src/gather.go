package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/api/compute/v1"
	"strings"
)

func getZonesFromRegion(region string) []string {

	result, err := computeService.Regions.Get(project, region).Context(ctx).Do()
	if err != nil {
		log.Fatal(err)
	}

	var zones []string
	for _, zone := range result.Zones {
		s := strings.Split(zone, "/")
		z := s[len(s)-1]
		zones = append(zones, z)
	}

	return zones
}

func (e *Exporter) gatherInstanceMetrics(ch chan<- prometheus.Metric) (*compute.MachineTypesListCall, error) {

	// Process all zones contained within the specified REGION
	zones := getZonesFromRegion(region)

	var err error
	var result *compute.MachineTypesListCall
	for _, zone := range zones {
		result = computeService.MachineTypes.List(project, zone)
		if err = result.Pages(ctx, func(page *compute.MachineTypeList) error {
			log.Debug("MachineTypes <RESULT>:", page.Items)
			for _, machineType := range page.Items {
				log.Debug("Data <machine>:", machineType)
				log.Debug("Zone <zone>:", zone)

				s := strings.Split(machineType.Zone, "/")
				z := s[len(s)-1]

				e.gaugeVecs["totalvCPUs"].With(prometheus.Labels{
					"instance_type": machineType.Name,
					"region":        region,
					"zone":          z,
				}).Set(float64(machineType.GuestCpus))
				e.gaugeVecs["totalMem"].With(prometheus.Labels{
					"instance_type": machineType.Name,
					"region":        region,
					"zone":          z,
				}).Set(float64(machineType.MemoryMb))
				e.gaugeVecs["maxStorage"].With(prometheus.Labels{
					"instance_type": machineType.Name,
					"region":        region,
					"zone":          z,
				}).Set(float64(machineType.MaximumPersistentDisksSizeGb))
				e.gaugeVecs["maxDisks"].With(prometheus.Labels{
					"instance_type": machineType.Name,
					"region":        region,
					"zone":          z,
				}).Set(float64(machineType.MaximumPersistentDisks))
			}

			return nil

		}); err != nil {
			log.Fatal(err)
		}
	}

	return result, err
}

func (e *Exporter) gatherImageMetrics(ch chan<- prometheus.Metric) (*compute.ImagesListCall, error) {

	var err error
	var result *compute.ImagesListCall
	for _, proj := range imageProjects {
		result := computeService.Images.List(proj)
		if err := result.Pages(ctx, func(page *compute.ImageList) error {
			log.Debug("ImageList <RESULT>:", page.Items)
			for _, image := range page.Items {
				log.Debug("Data <image>:", image)

				e.gaugeVecs["imgArchiveSize"].With(prometheus.Labels{
					"name":    image.Name,
					"family":  image.Family,
					"src_img": image.SourceImage,
				}).Set(float64(image.ArchiveSizeBytes))
				e.gaugeVecs["imgDiskSize"].With(prometheus.Labels{
					"name":    image.Name,
					"family":  image.Family,
					"src_img": image.SourceImage,
				}).Set(float64(image.DiskSizeGb))
				e.counterVecs["totalImages"].With(prometheus.Labels{
					"family":  image.Family,
					"src_img": image.SourceImage,
				}).Inc()
			}

			return nil
		}); err != nil {
			log.Fatal(err)
		}
	}

	return result, err

}

func (e *Exporter) gatherDiskMetrics(ch chan<- prometheus.Metric) (*compute.DiskTypesListCall, error) {

	// Process all zones contained within the specified REGION
	zones := getZonesFromRegion(region)

	var err error
	var result *compute.DiskTypesListCall
	for _, zone := range zones {
		result := computeService.DiskTypes.List(project, zone)
		if err := result.Pages(ctx, func(page *compute.DiskTypeList) error {
			log.Debug("DiskTypeList <RESULT>:", page.Items)
			for _, dt := range page.Items {
				log.Debug("Data <image>:", dt)

				e.gaugeVecs["diskTypeSize"].With(prometheus.Labels{
					"name":            dt.Name,
					"valid_disk_size": dt.ValidDiskSize,
					"region":          region,
					"zone":            zone,
				}).Set(float64(dt.DefaultDiskSizeGb))
				e.counterVecs["totalDiskTypes"].With(prometheus.Labels{
					"valid_disk_size": dt.ValidDiskSize,
					"region":          region,
					"zone":            zone,
				}).Inc()

			}

			return nil
		}); err != nil {
			log.Fatal(err)
		}
	}

	return result, err

}

func (e *Exporter) gatherRegionMetrics(ch chan<- prometheus.Metric) (*compute.RegionsListCall, error) {

	var err error
	result := computeService.Regions.List(project)
	if err := result.Pages(ctx, func(page *compute.RegionList) error {
		log.Debug("RegionList <RESULT>:", page.Items)
		for _, region := range page.Items {
			log.Debug("Data <region>:", region)

			e.counterVecs["totalRegions"].With(prometheus.Labels{
				"name":   region.Name,
				"status": region.Status,
			}).Inc()
		}

		return nil
	}); err != nil {
		log.Fatal(err)
	}

	return result, err

}

func (e *Exporter) gatherZoneMetrics(ch chan<- prometheus.Metric) (*compute.ZonesListCall, error) {

	var err error
	result := computeService.Zones.List(project)
	if err := result.Pages(ctx, func(page *compute.ZoneList) error {
		log.Debug("ZoneList <RESULT>:", page.Items)
		for _, zone := range page.Items {
			log.Debug("Data <zone>:", zone)

			s := strings.Split(zone.Region, "/")
			r := s[len(s)-1]

			e.counterVecs["totalZones"].With(prometheus.Labels{
				"name":   zone.Name,
				"status": zone.Status,
				"region": r,
			}).Inc()
		}

		return nil
	}); err != nil {
		log.Fatal(err)
	}

	return result, err

}
