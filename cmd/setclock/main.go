package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ecc1/dexcom"
)

func usage() {
	log.Fatalf("Usage: %s YYYY-MM-DD HH:MM:SS (or \"now\")", os.Args[0])
}

func main() {
	update := true
	t := time.Time{}
	switch len(os.Args) {
	case 1:
		update = false
	case 2:
		if os.Args[1] == "now" {
			t = time.Now()
		} else {
			usage()
		}
	case 3:
		t = parseTime(os.Args[1] + " " + os.Args[2])
	default:
		usage()
	}
	cgm := dexcom.Open()
	if update {
		cgm.SetDisplayTime(t)
	}
	if cgm.Error() != nil {
		log.Fatal(cgm.Error())
	}
	fmt.Println(cgm.ReadDisplayTime().Format(dexcom.UserTimeLayout))
}

func parseTime(date string) time.Time {
	t, err := time.ParseInLocation(dexcom.UserTimeLayout, date, time.Local)
	if err != nil {
		usage()
	}
	return t
}
