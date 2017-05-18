package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ecc1/dexcom"
)

var (
	pageNumber   = flag.Int("n", -1, "`page` number to read; -1 for most recent")
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
	pageNum := *pageNumber
	if pageNum == -1 {
		_, pageNum = cgm.ReadPageRange(pageType)
	}
	if cgm.Error() != nil {
		log.Fatal(cgm.Error())
	}
	log.Printf("reading %v page %d", pageType, pageNum)
	v := cgm.ReadPage(pageType, pageNum)
	if cgm.Error() != nil {
		log.Fatal(cgm.Error())
	}
	fmt.Printf("% X\n", v)
}
