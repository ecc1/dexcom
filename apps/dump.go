package main

import (
	"fmt"
	"github.com/ecc1/dexcom"
	"log"
)

func main() {
	err := dexcom.Open()
	if err != nil {
		log.Fatal(err)
	}

	db, err := dexcom.OpenDB()
	if err != nil {
		log.Fatal(err)
	}
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
	dexcom.ReadRecords(
		&dexcom.EGVRecord{},
		func(record dexcom.Record, context *dexcom.RecordContext) error {
			_, err = stmt.Exec(dexcom.GlucoseRow(record))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Print(".")
			n++
			if n%80 == 0 {
				fmt.Print("\n")
			}
			return err
		})
	if n%80 != 0 {
		fmt.Print("\n")
	}

	xact.Commit()
}
