package main

import (
	"database/sql"
	"fmt"
	"log"
	"os/user"
	"path/filepath"
	"time"

	"github.com/ecc1/dexcom"
	_ "github.com/mattn/go-sqlite3"
)

const (
	dbName = "glucose.db"
)

func main() {
	err := dexcom.Open()
	if err != nil {
		log.Fatal(err)
	}
	db := OpenDB()
	xact, err := db.Begin()
	if err != nil {
		db.Close()
		log.Fatal(err)
	}
	stmt, err := xact.Prepare("insert into glucose values (?, ?)")
	if err != nil {
		db.Close()
		log.Fatal(err)
	}
	n := 0
	dexcom.ReadRecords(
		&dexcom.EGVRecord{},
		func(record dexcom.Record, context *dexcom.RecordContext) error {
			_, err = stmt.Exec(glucoseRow(record))
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
	db.Close()
}

func OpenDB() *sql.DB {
	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	home := u.HomeDir
	db, err := sql.Open("sqlite3", filepath.Join(home, dbName))
	if err != nil {
		log.Fatal(err)
	}
	// use integer (Unix time); datetime type stores value as a string
	stmt := `
	  create table if not exists glucose
	    (time integer primary key on conflict ignore,
	     value integer)
	`
	_, err = db.Exec(stmt)
	if err != nil {
		db.Close()
		log.Fatal(err)
	}
	return db
}

// glucoseRow returns the time and glucose value from an EGVRecord
// suitable for inserting into the database.
func glucoseRow(record dexcom.Record) (int64, uint16) {
	egv := record.(*dexcom.EGVRecord)
	return egv.Timestamp.DisplayTime.Round(time.Second).Unix(), egv.Glucose
}
