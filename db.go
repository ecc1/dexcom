package dexcom

import (
	"database/sql"
	"log"
	"os/user"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const dbName = "glucose.db"

// InsertStmt can be used in a Query or Prepare call to insert
// a glucose row into the database.
const InsertStmt = "insert into glucose values (?, ?)"

// OpenDB opens the glucose database.
func OpenDB() (*sql.DB, error) {
	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	home := u.HomeDir
	db, err := sql.Open("sqlite3", filepath.Join(home, dbName))
	if err != nil {
		return nil, err
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
		return nil, err
	}
	return db, nil
}

// GlucoseRow returns the time and glucose value from an EGVRecord
// suitable for inserting into the database.
func GlucoseRow(record Record) (int64, uint16) {
	egv := record.(*EGVRecord)
	return egv.Timestamp.DisplayTime.Round(time.Second).Unix(), egv.Glucose
}
