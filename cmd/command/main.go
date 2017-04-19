package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/ecc1/dexcom"
)

func usage() {
	log.Fatalf("Usage: %s code [param ...]", os.Args[0])
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	code, err := strconv.ParseUint(os.Args[1], 0, 8)
	if err != nil {
		log.Fatal(err)
	}
	params := make([]byte, len(os.Args)-2)
	for i, arg := range os.Args[2:] {
		p, err := strconv.ParseUint(arg, 0, 8)
		if err != nil {
			log.Fatal(err)
		}
		params[i] = byte(p)
	}
	cgm := dexcom.Open()
	if cgm.Error() != nil {
		log.Fatal(cgm.Error())
	}
	fmt.Printf("command: %02X\n", code)
	if len(params) != 0 {
		fmt.Printf(" params: % X\n", params)
	}
	result := cgm.Cmd(dexcom.Command(code), params...)
	if cgm.Error() != nil {
		log.Fatal(cgm.Error())
	}
	fmt.Printf(" result: % X\n", result)
}
