package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/ecc1/dexcom"
)

const (
	userTimeLayout = "2006-01-02 15:04:05"
)

var (
	all        = flag.Bool("a", false, "get all records")
	numMinutes = flag.Int("n", 30, "number of `minutes` to get")
)

func main() {
	flag.Parse()
	cutoff := time.Time{}
	cgm := dexcom.Open()
	if cgm.Error() != nil {
		log.Fatal(cgm.Error())
	}
	if *all {
		log.Printf("retrieving all glucose records")
	} else {
		cutoff = time.Now().Add(-time.Duration(*numMinutes) * time.Minute)
		log.Printf("retrieving records since %s", cutoff.Format(userTimeLayout))
	}
	readings := cgm.GlucoseReadings(cutoff)
	if cgm.Error() != nil {
		log.Fatal(cgm.Error())
	}
	for _, r := range readings {
		printReading(r)
	}
}

func printReading(r dexcom.Record) {
	t := r.Time().Format(userTimeLayout)
	glucose, noise, unfiltered, filtered, rssi := "", "", "", "", ""
	if r.Egv != nil {
		glucose = fmt.Sprintf("%d", r.Egv.Glucose)
		noise = fmt.Sprintf("%d", r.Egv.Noise)
	}
	if r.Sensor != nil {
		unfiltered = fmt.Sprintf("%6d", r.Sensor.Unfiltered)
		filtered = fmt.Sprintf("%6d", r.Sensor.Filtered)
		rssi = fmt.Sprintf("%3d", r.Sensor.Rssi)
	}
	fmt.Printf("%s  %3s  %3s  %6s  %6s  %3s\n", t, glucose, noise, unfiltered, filtered, rssi)
}
