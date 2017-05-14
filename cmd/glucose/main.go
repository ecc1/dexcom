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
	csvFormat      = "csv"
	textFormat     = "text"
	jsonFormat     = "json"
)

var (
	all         = flag.Bool("a", false, "get all records")
	duration    = flag.Duration("d", time.Hour, "get `duration` worth of previous records")
	format      = flag.String("f", textFormat, "format in which to print records (csv, json, or text)")
	egv         = flag.Bool("g", true, "include glucose records")
	sensor      = flag.Bool("s", false, "include sensor records")
	calibration = flag.Bool("c", false, "include calibration records")
	meter       = flag.Bool("m", false, "include meter records")

	recordTypes = []struct {
		flag *bool
		page dexcom.PageType
	}{
		{egv, dexcom.EGVData},
		{sensor, dexcom.SensorData},
		{calibration, dexcom.CalibrationData},
		{meter, dexcom.MeterData},
	}
)

func main() {
	flag.Parse()
	switch *format {
	case csvFormat, jsonFormat, textFormat:
	default:
		flag.Usage()
		return
	}
	var cutoff time.Time
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
	scans := scanRecords(cgm, cutoff)
	results := dexcom.MergeHistory(scans...)
	if *format == jsonFormat {
		e := json.NewEncoder(os.Stdout)
		e.SetIndent("", "  ")
		err := e.Encode(results)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	if *format == csvFormat {
		fmt.Printf("Time,Type,Glucose,Sensor,Slope,Intercept,Scale,Decay\n")
	}
	for _, r := range results {
		printRecord(r)
	}
}

func scanRecords(cgm *dexcom.CGM, cutoff time.Time) [][]dexcom.Record {
	var scans [][]dexcom.Record
	for _, t := range recordTypes {
		if !*t.flag {
			continue
		}
		var v []dexcom.Record
		// Special case when both EGV and sensor records are requested.
		if t.page == dexcom.EGVData && *sensor {
			v = cgm.GlucoseReadings(cutoff)
		} else if t.page == dexcom.SensorData && *egv {
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
	return scans
}

func printRecord(r dexcom.Record) {
	t := r.Time().Format(userTimeLayout)
	printGlucose(t, r.EGV, r.Sensor)
	printCalibration(t, r.Calibration)
	printMeter(t, r.Meter)
}

func printGlucose(t string, e *dexcom.EGVInfo, s *dexcom.SensorInfo) {
	if e == nil && s == nil {
		return
	}
	glucose, noise, unfiltered, filtered, rssi := "", "", "", "", ""
	if e != nil {
		glucose = fmt.Sprintf("%d", e.Glucose)
		noise = fmt.Sprintf("%d", e.Noise)
	}
	if s != nil {
		unfiltered = fmt.Sprintf("%d", s.Unfiltered)
		filtered = fmt.Sprintf("%d", s.Filtered)
		rssi = fmt.Sprintf("%d", s.RSSI)
	}
	switch *format {
	case csvFormat:
		fmt.Printf("%s,%s,%s,%s\n", t, "G", glucose, unfiltered)
	case textFormat:
		fmt.Printf("%s  %3s  %3s  %6s  %6s  %3s\n", t, glucose, noise, unfiltered, filtered, rssi)
	}
}

func printCalibration(t string, cal *dexcom.CalibrationInfo) {
	if cal == nil {
		return
	}
	switch *format {
	case csvFormat:
		fmt.Printf("%s,%s,,,%g,%g,%g,%g\n", t, "C", cal.Slope, cal.Intercept, cal.Scale, cal.Decay)
		for _, d := range cal.Data {
			t = d.TimeEntered.Format(userTimeLayout)
			fmt.Printf("%s,%s,%d,%d\n", t, "D", d.Glucose, d.Raw)
		}
	case textFormat:
		fmt.Printf("%s  %-5s  %g  %g  %g  %g\n", t, "CAL", cal.Slope, cal.Intercept, cal.Scale, cal.Decay)
		for _, d := range cal.Data {
			t = d.TimeEntered.Format(userTimeLayout)
			fmt.Printf("%s  %-5s  %3d  %6d\n", t, "DATA", d.Glucose, d.Raw)
		}
	}
}

func printMeter(t string, m *dexcom.MeterInfo) {
	if m == nil {
		return
	}
	switch *format {
	case csvFormat:
		fmt.Printf("%s,%s,%d\n", t, "M", m.Glucose)
	case textFormat:
		fmt.Printf("%s  %-5s  %3d\n", t, "METER", m.Glucose)
	}
}
