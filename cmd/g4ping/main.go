package main

import (
	"log"

	"github.com/ecc1/dexcom"
)

func main() {
	cgm := dexcom.Open()
	cgm.Cmd(dexcom.Ping)
	if cgm.Error() != nil {
		log.Fatal(cgm.Error())
	}
}
