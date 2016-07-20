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
	proc := func(v []byte, _ dexcom.RecordContext) error {
		r := dexcom.EGVRecord{}
		err := r.Unmarshal(v)
		if err != nil {
			return err
		}
		_, err = stmt.Exec(glucoseRow(r))
		if err != nil {
			return err
		}
		fmt.Print(".")
		n++
		if n%80 == 0 {
			fmt.Println()
		}
		return nil
	}
	err = dexcom.ReadRecords(dexcom.EGV_DATA, proc)
	if n%80 != 0 {
		fmt.Println()
	}
	if err != nil {
		db.Close()
		log.Fatal(err)
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
func glucoseRow(egv dexcom.EGVRecord) (int64, uint16) {
	return egv.Timestamp.DisplayTime.Round(time.Second).Unix(), egv.Glucose
}
