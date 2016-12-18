package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ecc1/dexcom"
)

const (
	userTimeLayout = "2006-01-02 15:04:05"
)

var (
	all         = flag.Bool("a", false, "get all records")
	duration    = flag.Duration("d", time.Hour, "get `duration` worth of previous records")
	format      = flag.String("f", "text", "format in which to print records (csv, json, or text)")
	egv         = flag.Bool("g", true, "include glucose records")
	sensor      = flag.Bool("s", false, "include sensor records")
	calibration = flag.Bool("c", false, "include calibration records")
	meter       = flag.Bool("m", false, "include meter records")

	recordTypes = []struct {
		flag *bool
		page dexcom.PageType
	}{
		{egv, dexcom.EGV_DATA},
		{sensor, dexcom.SENSOR_DATA},
		{calibration, dexcom.CAL_SET},
		{meter, dexcom.METER_DATA},
	}
)

func main() {
	flag.Parse()
	switch *format {
	case "csv", "json", "text":
	default:
		flag.Usage()
		return
	}
	cutoff := time.Time{}
	cgm := dexcom.Open()
	if cgm.Error() != nil {
		log.Fatal(cgm.Error())
	}
	if *all {
		log.Printf("retrieving entire record history")
		*egv = true
		*sensor = true
		*calibration = true
		*meter = true
	} else {
		cutoff = time.Now().Add(-*duration)
		log.Printf("retrieving records since %s", cutoff.Format(userTimeLayout))
	}
	scans := [][]dexcom.Record{}
	for _, t := range recordTypes {
		if !*t.flag {
			continue
		}
		v := []dexcom.Record{}
		// Special case when both EGV and sensor records are requested.
		if t.page == dexcom.EGV_DATA && *sensor {
			v = cgm.GlucoseReadings(cutoff)
		} else if t.page == dexcom.SENSOR_DATA && *egv {
			continue
		} else {
			v = cgm.ReadHistory(t.page, cutoff)
		}
		if len(v) != 0 {
			scans = append(scans, v)
		}
	}
	if cgm.Error() != nil {
		log.Fatal(cgm.Error())
	}
	if len(scans) == 0 {
		log.Fatal("no records found")
	}
	results := dexcom.MergeHistory(scans...)
	if *format == "json" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		err := enc.Encode(results)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	if *format == "csv" {
		fmt.Printf("Time,Type,Glucose,Sensor,Slope,Intercept,Scale,Decay\n")
	}
	for _, r := range results {
		printRecord(r)
	}
}

func printRecord(r dexcom.Record) {
	t := r.Time().Format(userTimeLayout)
	if r.Egv != nil || r.Sensor != nil {
		glucose, noise, unfiltered, filtered, rssi := "", "", "", "", ""
		if r.Egv != nil {
			glucose = fmt.Sprintf("%d", r.Egv.Glucose)
			noise = fmt.Sprintf("%d", r.Egv.Noise)
		}
		if r.Sensor != nil {
			unfiltered = fmt.Sprintf("%d", r.Sensor.Unfiltered)
			filtered = fmt.Sprintf("%d", r.Sensor.Filtered)
			rssi = fmt.Sprintf("%d", r.Sensor.Rssi)
		}
		switch *format {
		case "csv":
			fmt.Printf("%s,%s,%s,%s\n", t, "G", glucose, unfiltered)
		case "text":
			fmt.Printf("%s  %3s  %3s  %6s  %6s  %3s\n", t, glucose, noise, unfiltered, filtered, rssi)
		}
		return
	}
	if r.Calibration != nil {
		cal := r.Calibration
		switch *format {
		case "csv":
			fmt.Printf("%s,%s,,,%f,%f,%f,%f\n", t, "C", cal.Slope, cal.Intercept, cal.Scale, cal.Decay)
			for _, d := range cal.Data {
				t := d.TimeEntered.Format(userTimeLayout)
				fmt.Printf("%s,%s,%d,%d\n", t, "D", d.Glucose, d.Raw)
			}
		case "text":
			fmt.Printf("%s  %-5s  %.2f  %.2f  %.2f  %.2f\n", t, "CAL", cal.Slope, cal.Intercept, cal.Scale, cal.Decay)
			for _, d := range cal.Data {
				t := d.TimeEntered.Format(userTimeLayout)
				fmt.Printf("%s  %-5s  %3d  %6d\n", t, "DATA", d.Glucose, d.Raw)
			}
		}
		return
	}
	if r.Meter != nil {
		m := r.Meter
		switch *format {
		case "csv":
			fmt.Printf("%s,%s,%d\n", t, "M", m.Glucose)
		case "text":
			fmt.Printf("%s  %-5s  %3d\n", t, "METER", m.Glucose)
		}
		return
	}
}
