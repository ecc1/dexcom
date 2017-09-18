package main

// Backfill missing CGM readings to Nightscout.

import (
	"flag"
	"log"
	"time"

	"github.com/ecc1/dexcom"
	"github.com/ecc1/nightscout"
)

var (
	checkDuration = flag.Duration("c", time.Hour, "`duration` to check")
	gapsOnlyFlag  = flag.Bool("g", false, "list Nightscout gaps only")
	noUploadFlag  = flag.Bool("s", false, "simulate Nightscout uploads")
	verboseFlag   = flag.Bool("v", false, "verbose mode")

	gapDuration = 7 * time.Minute
	timeLayout  = "2006-01-02 15:04:05"

	pageTypes = []dexcom.PageType{
		dexcom.SensorData,
		dexcom.EGVData,
		dexcom.MeterData,
		dexcom.CalibrationData,
	}
)

func main() {
	flag.Parse()
	nightscout.SetNoUpload(*noUploadFlag)
	nightscout.SetVerbose(*verboseFlag)
	gaps, cutoff := findGaps()
	if len(gaps) == 0 {
		return
	}
	upload(nightscout.Missing(getRecords(cutoff), gaps))
}

func findGaps() ([]nightscout.Gap, time.Time) {
	now := time.Now()
	cutoff := now.Add(-*checkDuration)
	gaps, err := nightscout.Gaps(cutoff, gapDuration)
	if err != nil {
		log.Fatal(err)
	}
	if len(gaps) == 0 {
		log.Printf("no gaps found")
		return nil, cutoff
	}
	if *gapsOnlyFlag {
		printGaps(gaps)
		return nil, cutoff
	}
	if *verboseFlag {
		printGaps(gaps)
	}
	// No need to retrieve records further than beginning of earliest gap.
	earliest := gaps[len(gaps)-1].Start
	if cutoff.Before(earliest) {
		cutoff = earliest
	}
	return gaps, cutoff
}

func getRecords(cutoff time.Time) nightscout.Entries {
	cgm := dexcom.Open()
	if cgm.Error() != nil {
		log.Fatal(cgm.Error())
	}
	log.Printf("retrieving Dexcom records since %s", cutoff)
	var records []dexcom.Records
	for _, page := range pageTypes {
		records = append(records, cgm.ReadHistory(page, cutoff))
	}
	if cgm.Error() != nil {
		log.Fatal(cgm.Error())
	}
	return dexcom.NightscoutEntries(dexcom.MergeHistory(records...))
}

func upload(entries nightscout.Entries) {
	log.Printf("uploading %d entries to Nightscout", len(entries))
	for _, e := range entries {
		err := nightscout.Upload("POST", "entries", e)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func printGaps(gaps []nightscout.Gap) {
	for _, g := range gaps {
		start := g.Start
		finish := g.Finish
		gap := finish.Sub(start)
		log.Printf("%v gap from %s to %s", gap, start.Format(timeLayout), finish.Format(timeLayout))
	}
}
