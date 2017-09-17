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
	csvFormat  = "csv"
	jsonFormat = "json"
	nsFormat   = "ns"
	textFormat = "text"
)

var (
	all      = flag.Bool("a", false, "get all records")
	duration = flag.Duration("d", time.Hour, "get `duration` worth of previous records")
	format   = flag.String("f", textFormat, "format in which to print records (csv, json, ns, or text)")

	egv         = flag.Bool("e", true, "include EGV records")
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
	case csvFormat, jsonFormat, nsFormat, textFormat:
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
	if *format == csvFormat {
		fmt.Printf("Time,Type,Glucose,Sensor,Slope,Intercept,Scale,Decay\n")
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
	switch info := r.Info.(type) {
	case dexcom.SensorInfo:
		printSensor(t, info)
	case dexcom.EGVInfo:
		printEGV(t, info)
	case dexcom.CalibrationInfo:
		printCalibration(t, info)
	case dexcom.MeterInfo:
		printMeter(t, info)
	default:
		panic(fmt.Sprintf("unexpected record %+v", r))
	}
}

func printSensor(t string, s dexcom.SensorInfo) {
	switch *format {
	case csvFormat:
		fmt.Printf("%s,G,,%d\n", t, s.Unfiltered)
	case textFormat:
		fmt.Printf("%s            %6d  %6d  %3d\n", t, s.Unfiltered, s.Filtered, s.RSSI)
	}
}

func printEGV(t string, e dexcom.EGVInfo) {
	switch *format {
	case csvFormat:
		fmt.Printf("%s,G,%d,\n", t, e.Glucose)
	case textFormat:
		fmt.Printf("%s  %3d  %3d\n", t, e.Glucose, e.Noise)
	}
}

func printCalibration(t string, cal dexcom.CalibrationInfo) {
	switch *format {
	case csvFormat:
		fmt.Printf("%s,%s,,,%g,%g,%g,%g\n", t, "C", cal.Slope, cal.Intercept, cal.Scale, cal.Decay)
		for _, d := range cal.Data {
			t = d.TimeEntered.Format(dexcom.UserTimeLayout)
			fmt.Printf("%s,%s,%d,%d\n", t, "D", d.Glucose, d.Raw)
		}
	case textFormat:
		fmt.Printf("%s  %-5s  %g  %g  %g  %g\n", t, "CAL", cal.Slope, cal.Intercept, cal.Scale, cal.Decay)
		for _, d := range cal.Data {
			t = d.TimeEntered.Format(dexcom.UserTimeLayout)
			fmt.Printf("%s  %-5s  %3d  %6d\n", t, "DATA", d.Glucose, d.Raw)
		}
	}
}

func printMeter(t string, m dexcom.MeterInfo) {
	switch *format {
	case csvFormat:
		fmt.Printf("%s,%s,%d\n", t, "M", m.Glucose)
	case textFormat:
		fmt.Printf("%s  %-5s  %3d\n", t, "METER", m.Glucose)
	}
}
