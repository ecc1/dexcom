package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ecc1/dexcom"
)

var (
	pageTypeFlag = flag.Int("p", int(dexcom.EGVData), "`page type` to read")
	numRecords   = flag.Int("n", 10, "number of `records` to get")
)

// nolint: gas
func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "Page Types:\n")
		for p := dexcom.FirstPageType; p <= dexcom.LastPageType; p++ {
			fmt.Fprintf(os.Stderr, "  %2d = %v\n", int(p), p)
		}
	}
	flag.Parse()
	pageType := dexcom.PageType(*pageTypeFlag)
	if pageType < dexcom.FirstPageType || dexcom.LastPageType < pageType {
		fmt.Fprintf(os.Stderr, "invalid page type (%d)\n", *pageTypeFlag)
		flag.Usage()
		os.Exit(1)
	}
	cgm := dexcom.Open()
	results := cgm.ReadCount(pageType, *numRecords)
	if cgm.Error() != nil {
		log.Fatal(cgm.Error())
	}
	for _, r := range results {
		b, err := json.MarshalIndent(r, "", "  ")
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(string(b))
	}
}
