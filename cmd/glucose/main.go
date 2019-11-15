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
	jsonFormat = "json"
	nsFormat   = "ns"
	textFormat = "text"
)

var (
	all      = flag.Bool("a", false, "get all records")
	duration = flag.Duration("d", time.Hour, "get `duration` worth of previous records")
	since    = flag.String("t", "", "get records since the specified `time` in RFC3339 format")
	format   = flag.String("f", textFormat, "format in which to print records (json, ns, or text)")

	egv         = flag.Bool("e", true, "include EGV records")
	sensor      = flag.Bool("s", false, "include sensor records")
	calibration = flag.Bool("c", false, "include calibration records")
	insertion   = flag.Bool("i", false, "include insertion records")
	meter       = flag.Bool("m", false, "include meter records")

	recordTypes = []struct {
		flag *bool
		page dexcom.PageType
	}{
		{egv, dexcom.EGVData},
		{sensor, dexcom.SensorData},
		{calibration, dexcom.CalibrationData},
		{insertion, dexcom.InsertionTimeData},
		{meter, dexcom.MeterData},
	}
)

func main() {
	flag.Parse()
	switch *format {
	case jsonFormat, nsFormat, textFormat:
	default:
		flag.Usage()
		return
	}
	var cutoff time.Time
	var err error
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
	} else if *since != "" {
		cutoff, err = time.Parse(dexcom.JSONTimeLayout, *since)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		cutoff = time.Now().Add(-*duration)
	}
	if !*all {
		log.Printf("retrieving records since %s", cutoff.Format(dexcom.UserTimeLayout))
	}
	scans := scanRecords(cgm, cutoff)
	results := dexcom.MergeHistory(scans...)
	if len(results) == 0 {
		return
	}
	if *format == nsFormat {
		printJSON(dexcom.NightscoutEntries(results))
		return
	}
	if *format == jsonFormat {
		printJSON(results)
		return
	}
	for _, r := range results {
		printRecord(r)
	}
}

func scanRecords(cgm *dexcom.CGM, cutoff time.Time) []dexcom.Records {
	var scans []dexcom.Records
	for _, t := range recordTypes {
		if !*t.flag {
			continue
		}
		scans = append(scans, cgm.ReadHistory(t.page, cutoff))
	}
	if cgm.Error() != nil {
		log.Fatal(cgm.Error())
	}
	if len(scans) == 0 {
		log.Fatal("no records found")
	}
	return scans
}

func printJSON(v interface{}) {
	e := json.NewEncoder(os.Stdout)
	e.SetIndent("", "  ")
	err := e.Encode(v)
	if err != nil {
		log.Fatal(err)
	}
}

func printRecord(r dexcom.Record) {
	t := r.Time().Format(dexcom.UserTimeLayout)
	printSensor(t, r.Sensor)
	printEGV(t, r.EGV)
	printCalibration(t, r.Calibration)
	printInsertion(t, r.Insertion)
	printMeter(t, r.Meter)
}

func printSensor(t string, s *dexcom.SensorInfo) {
	if s == nil {
		return
	}
	fmt.Printf("%s            %6d  %6d  %3d\n", t, s.Unfiltered, s.Filtered, s.RSSI)
}

func printEGV(t string, e *dexcom.EGVInfo) {
	if e == nil {
		return
	}
	fmt.Printf("%s  %3d  %3d\n", t, e.Glucose, e.Noise)
}

func printCalibration(t string, cal *dexcom.CalibrationInfo) {
	if cal == nil {
		return
	}
	fmt.Printf("%s  %-5s  %g  %g  %g  %g\n", t, "CAL", cal.Slope, cal.Intercept, cal.Scale, cal.Decay)
	for _, d := range cal.Data {
		t = d.TimeEntered.Format(dexcom.UserTimeLayout)
		fmt.Printf("%s  %-5s  %3d  %6d\n", t, "DATA", d.Glucose, d.Raw)
	}
}

func printInsertion(t string, i *dexcom.InsertionInfo) {
	if i == nil {
		return
	}
	fmt.Printf("%s  %-5s  %s\n", t, "SITE", i.Event)
}

func printMeter(t string, m *dexcom.MeterInfo) {
	if m == nil {
		return
	}
	fmt.Printf("%s  %-5s  %3d\n", t, "METER", m.Glucose)
}
