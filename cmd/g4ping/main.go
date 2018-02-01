package main

import (
	"log"

	"github.com/ecc1/dexcom"
	"github.com/ecc1/papertrail"
)

func main() {
	papertrail.StartLogging()
	cgm := dexcom.Open()
	cgm.Cmd(dexcom.Ping)
	if cgm.Error() != nil {
		log.Fatal(cgm.Error())
	}
}
