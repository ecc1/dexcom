package main

import (
	"database/sql"
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
	db := OpenDB()
	xact, err := db.Begin()
	if err != nil {
		db.Close()
		log.Fatal(err)
	}
	cgm := dexcom.Open()
	egv := cgm.ReadHistory(dexcom.EGV_DATA, time.Time{})
	if cgm.Error() != nil {
		log.Fatal(cgm.Error())
	}
	stmt, err := xact.Prepare("insert into glucose values (?, ?)")
	if err != nil {
		db.Close()
		log.Fatal(err)
	}
	for _, r := range egv {
		_, err = stmt.Exec(glucoseRow(r))
		if err != nil {
			db.Close()
			log.Fatal(err)
		}
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
func glucoseRow(r dexcom.Record) (int64, uint16) {
	t := r.Timestamp.DisplayTime
	return t.Unix(), r.Egv.Glucose
}
