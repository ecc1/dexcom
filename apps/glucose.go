package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/ecc1/dexcom"
)

const (
	timeLayout = "2006/01/02 15:04:05"
)

var (
	all        = flag.Bool("a", false, "get all records")
	numMinutes = flag.Int("n", 30, "number of `minutes` to get")
)

func main() {
	flag.Parse()
	cutoff := time.Time{}
	cgm := dexcom.Open()
	if *all {
		log.Printf("retrieving all glucose records")
	} else {
		cutoff = time.Now().Add(-time.Duration(*numMinutes) * time.Minute)
		log.Printf("retrieving records since %s", cutoff.Format(time.RFC3339))
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
	t := r.Timestamp.DisplayTime
	if r.Egv != nil && r.Sensor != nil {
		fmt.Printf("%s  %3d  %6d  %6d\n", t.Format(timeLayout), r.Egv.Glucose, r.Sensor.Unfiltered, r.Sensor.Filtered)
	} else if r.Egv != nil {
		fmt.Printf("%s  %3d\n", t.Format(timeLayout), r.Egv.Glucose)
	} else if r.Sensor != nil {
		fmt.Printf("%s       %6d  %6d\n", t.Format(timeLayout), r.Sensor.Unfiltered, r.Sensor.Filtered)
	}
}
