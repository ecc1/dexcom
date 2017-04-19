package main

import (
	"fmt"
	"log"

	"github.com/ecc1/dexcom"
)

func main() {
	cgm := dexcom.Open()
	if cgm.Error() != nil {
		log.Fatal(cgm.Error())
	}
	fmt.Printf("                         start   end\n")
	for t := dexcom.FirstPageType; t <= dexcom.LastPageType; t++ {
		first, last := cgm.ReadPageRange(t)
		if cgm.Error() != nil {
			fmt.Printf("%v: %v\n", t, cgm.Error())
			continue
		}
		fmt.Printf("%-24s  %4d  %4d\n", t.String(), first, last)
	}
}
