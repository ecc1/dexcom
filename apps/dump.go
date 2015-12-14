package main

import (
	"fmt"
	"github.com/ecc1/dexcom"
	"log"
)

func main() {
	dev, err := dexcom.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer dev.Close()

	db := dexcom.OpenDB()
	defer db.Close()

	xact, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := xact.Prepare(dexcom.InsertStmt)
	if err != nil {
		log.Fatal(err)
	}
	n := 0
	proc := func(v []byte, context *dexcom.RecordContext) error {
		_, err = stmt.Exec(dexcom.GlucoseRow(dexcom.UnmarshalEGVRecord(v)))
		fmt.Print(".")
		n++
		if n%80 == 0 {
			fmt.Print("\n")
		}
		return err
	}
	dev.ReadRecords(dexcom.EGV_DATA, proc)
	if n%80 != 0 {
		fmt.Print("\n")
	}
	xact.Commit()
}
