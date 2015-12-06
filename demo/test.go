package main

import (
	"fmt"
	"log"

	"github.com/ecc1/dexcom"
)

func printXMLData(name string, xmlData dexcom.XMLData) {
	fmt.Printf("%s:\n", name)
	for k, v := range xmlData {
		fmt.Printf("    %s: %s\n", k, v)
	}
}

func main() {
	dev, err := dexcom.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer dev.Close()

	displayTime, err := dev.ReadDisplayTime()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("display time:", displayTime)

	id, err := dev.Cmd(dexcom.READ_TRANSMITTER_ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("transmitter ID: %s\n", string(id))

	fw, err := dev.ReadFirmwareHeader()
	if err != nil {
		log.Fatal(err)
	}
	printXMLData("firmware header", fw)

	xr, err := dev.ReadXMLRecord(dexcom.MANUFACTURING_DATA)
	if err != nil {
		log.Fatal(err)
	}
	printXMLData("manufacturing data", xr.XML)
	fmt.Printf("    %v\n", xr.Timestamp)

	xr, err = dev.ReadXMLRecord(dexcom.PC_SOFTWARE_PARAMETER)
	if err != nil {
		log.Fatal(err)
	}
	printXMLData("PC software parameter", xr.XML)
	fmt.Printf("    %v\n", xr.Timestamp)
}
