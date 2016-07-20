package main

import (
	"fmt"
	"log"

	"github.com/ecc1/dexcom"
)

func main() {
	err := dexcom.Open()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("                         start   end\n")
	for t := dexcom.FirstRecordType; t <= dexcom.LastRecordType; t++ {
		context, err := dexcom.ReadPageRange(t)
		if err != nil {
			fmt.Printf("%v: %v\n", t, err)
			continue
		}
		fmt.Printf("%-24s  %4d  %4d\n", t.String(), context.StartPage, context.EndPage)
	}
}
