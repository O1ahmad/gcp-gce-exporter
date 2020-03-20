package main

import (
    "strings"
    "golang.org/x/net/context"
    "golang.org/x/oauth2/google"
    "google.golang.org/api/compute/v1"
	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
)

func (e *Exporter) gatherInstanceMetrics(ch chan<- prometheus.Metric) (*compute.MachineTypesListCall, error) {

    ctx := context.Background()
    computeClient, err := google.DefaultClient(ctx, compute.CloudPlatformScope)
    if err != nil {
        log.Fatal(err)
    }
    computeService, err := compute.New(computeClient)
    if err != nil {
        log.Fatal(err)
    }

    result := computeService.MachineTypes.List(project, "us-east1-a")
    if err := result.Pages(ctx, func(page *compute.MachineTypeList) error {
        log.Debug("MachineTypes <RESULT>:", page.Items)
        for _, machineType := range page.Items {
            log.Debug("Data <machine>:", machineType)

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

	return result, err

}

func (e *Exporter) gatherImageMetrics(ch chan<- prometheus.Metric) (*compute.ImagesListCall, error) {

    ctx := context.Background()
    computeClient, err := google.DefaultClient(ctx, compute.CloudPlatformScope)
    if err != nil {
        log.Fatal(err)
    }
    computeService, err := compute.New(computeClient)
    if err != nil {
        log.Fatal(err)
    }

    result := computeService.Images.List(project)
    if err := result.Pages(ctx, func(page *compute.ImageList) error {
        log.Debug("ImageList <RESULT>:", page.Items)
        for _, image := range page.Items {
            log.Debug("Data <image>:", image)

            e.gaugeVecs["imgArchiveSize"].With(prometheus.Labels{
                "name":          image.Name,
                "family":        image.Family,
                "src_img":       image.SourceImage,
            }).Set(float64(image.ArchiveSizeBytes))
            e.gaugeVecs["imgDiskSize"].With(prometheus.Labels{
                "name":          image.Name,
                "family":        image.Family,
                "src_img":       image.SourceImage,
            }).Set(float64(image.DiskSizeGb))
            e.counterVecs["totalImages"].With(prometheus.Labels{
                "family":        image.Family,
                "src_img":       image.SourceImage,
			}).Inc()
        }

        return nil
    }); err != nil {
         log.Fatal(err)
    }

	return result, err

}

func (e *Exporter) gatherDiskMetrics(ch chan<- prometheus.Metric) (*compute.DiskTypesListCall, error) {

    ctx := context.Background()
    computeClient, err := google.DefaultClient(ctx, compute.CloudPlatformScope)
    if err != nil {
        log.Fatal(err)
    }
    computeService, err := compute.New(computeClient)
    if err != nil {
        log.Fatal(err)
    }

    result := computeService.DiskTypes.List(project, zone)
    if err := result.Pages(ctx, func(page *compute.DiskTypeList) error {
        log.Debug("DiskTypeList <RESULT>:", page.Items)
        for _, dt := range page.Items {
            log.Debug("Data <image>:", dt)

            s := strings.Split(dt.Zone, "/")
            z := s[len(s)-1]

            e.gaugeVecs["diskTypeSize"].With(prometheus.Labels{
                "name":          dt.Name,
                "valid_disk_size": dt.ValidDiskSize,
                "region":        region,
                "zone":          z,
            }).Set(float64(dt.DefaultDiskSizeGb))
            e.counterVecs["totalDiskTypes"].With(prometheus.Labels{
                "valid_disk_size": dt.ValidDiskSize,
                "region":        region,
                "zone":          z,
			}).Inc()

        }

        return nil
    }); err != nil {
         log.Fatal(err)
    }

	return result, err

}

func (e *Exporter) gatherRegionMetrics(ch chan<- prometheus.Metric) (*compute.RegionsListCall, error) {

    ctx := context.Background()
    computeClient, err := google.DefaultClient(ctx, compute.CloudPlatformScope)
    if err != nil {
        log.Fatal(err)
    }
    computeService, err := compute.New(computeClient)
    if err != nil {
        log.Fatal(err)
    }

    result := computeService.Regions.List(project)
    if err := result.Pages(ctx, func(page *compute.RegionList) error {
        log.Debug("RegionList <RESULT>:", page.Items)
        for _, region := range page.Items {
            log.Debug("Data <region>:", region)

            e.counterVecs["totalRegions"].With(prometheus.Labels{
                "name": region.Name,
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

    ctx := context.Background()
    computeClient, err := google.DefaultClient(ctx, compute.CloudPlatformScope)
    if err != nil {
        log.Fatal(err)
    }
    computeService, err := compute.New(computeClient)
    if err != nil {
        log.Fatal(err)
    }

    result := computeService.Zones.List(project)
    if err := result.Pages(ctx, func(page *compute.ZoneList) error {
        log.Debug("ZoneList <RESULT>:", page.Items)
        for _, zone := range page.Items {
            log.Debug("Data <zone>:", zone)

            s := strings.Split(zone.Region, "/")
            r := s[len(s)-1]

            e.counterVecs["totalZones"].With(prometheus.Labels{
                "name": zone.Name,
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
