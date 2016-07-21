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
	err := dexcom.Open()
	if err != nil {
		log.Fatal(err)
	}
	if *all {
		log.Printf("retrieving all glucose records")
	} else {
		cutoff = time.Now().Add(-time.Duration(*numMinutes) * time.Minute)
		log.Printf("retrieving records since %s", cutoff.Format(time.RFC3339))
	}
	readings, err := dexcom.GlucoseReadings(cutoff)
	if err != nil {
		log.Fatal(err)
	}
	for _, r := range readings {
		printReading(r)
	}
}

func printReading(r dexcom.GlucoseReading) {
	t := r.Sensor.Timestamp.DisplayTime
	if t.IsZero() {
		t = r.Egv.Timestamp.DisplayTime
	}
	fmt.Printf("%s  %3d  %6d  %6d\n", t.Format(timeLayout), r.Egv.Glucose, r.Sensor.Unfiltered, r.Sensor.Filtered)
}
