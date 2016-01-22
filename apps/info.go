package main

import (
	"fmt"
	"log"

	"github.com/ecc1/dexcom"
)

func main() {
	err := dexcom.Open()
	if err != nil {
		log.Fatal(err)
	}
	var numRecords int
	var startPage, endPage uint32
	proc := func(_ []byte, context *dexcom.RecordContext) error {
		numRecords++
		startPage, endPage = context.StartPage, context.EndPage
		return nil
	}
	for t := dexcom.FirstRecordType; t <= dexcom.LastRecordType; t++ {
		numRecords = 0
		startPage = 0
		endPage = 0
		dexcom.ReadRecords(t, proc)
		fmt.Printf("%24s: %4d records, pp %d–%d\n", t.String(), numRecords, startPage, endPage)
	}
}
