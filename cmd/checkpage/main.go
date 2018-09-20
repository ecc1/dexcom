package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ecc1/dexcom"
)

var (
	nsFlag = flag.Bool("n", false, "print records in Nightscout format")
)

func main() {
	flag.Parse()
	for _, file := range flag.Args() {
		f, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		data, err := readBytes(f)
		_ = f.Close()
		if err != nil {
			log.Fatal(err)
		}
		readRecords(data)
	}
}

func readRecords(data []byte) {
	page, err := dexcom.UnmarshalPage(data)
	if err != nil {
		log.Fatal(err)
	}
	decoded, err := dexcom.UnmarshalRecords(page.Type, page.Records)
	printResults(decoded)
	if err != nil {
		log.Fatal(err)
	}
}

func printResults(results dexcom.Records) {
	e := json.NewEncoder(os.Stdout)
	e.SetIndent("", "  ")
	var err error
	if *nsFlag {
		err = e.Encode(dexcom.NightscoutEntries(results))
	} else {
		err = e.Encode(results)
	}
	if err != nil {
		log.Fatal(err)
	}
}

func readBytes(r io.Reader) ([]byte, error) {
	var data []byte
	for {
		var b byte
		n, err := fmt.Fscanf(r, "%02x", &b)
		if n == 0 {
			break
		}
		if err != nil {
			return data, err
		}
		data = append(data, b)
	}
	return data, nil
}
