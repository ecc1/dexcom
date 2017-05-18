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

var (
	all          = flag.Bool("a", false, "get all records")
	duration     = flag.Duration("d", time.Hour, "get `duration` worth of previous records")
	pageNumber   = flag.Int("n", -1, "`page` number to read")
	pageTypeFlag = flag.Int("p", int(dexcom.EGVData), "page `type` to read")
)

// nolint
func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "Page Types:\n")
	for p := dexcom.FirstPageType; p <= dexcom.LastPageType; p++ {
		fmt.Fprintf(os.Stderr, "  %2d = %v\n", int(p), p)
	}
}

func main() {
	flag.Usage = usage
	flag.Parse()
	pageType := dexcom.PageType(*pageTypeFlag)
	if pageType < dexcom.FirstPageType || dexcom.LastPageType < pageType {
		// nolint
		fmt.Fprintf(os.Stderr, "invalid page type (%d)\n", *pageTypeFlag)
		flag.Usage()
		os.Exit(1)
	}
	cgm := dexcom.Open()
	if cgm.Error() != nil {
		log.Fatal(cgm.Error())
	}
	var results []dexcom.Record
	if *pageNumber != -1 {
		results = cgm.ReadRecords(pageType, *pageNumber)
	} else {
		var cutoff time.Time
		if *all {
			log.Printf("retrieving entire record history")
		} else {
			cutoff = time.Now().Add(-*duration)
			log.Printf("retrieving records since %s", cutoff.Format(dexcom.UserTimeLayout))
		}
		results = cgm.ReadHistory(pageType, cutoff)
	}
	if cgm.Error() != nil {
		log.Fatal(cgm.Error())
	}
	e := json.NewEncoder(os.Stdout)
	e.SetIndent("", "  ")
	err := e.Encode(results)
	if err != nil {
		log.Fatal(err)
	}
}
