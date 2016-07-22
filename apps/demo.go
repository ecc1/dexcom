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
	cgm := dexcom.Open()
	if cgm.Error() != nil {
		log.Fatal(cgm.Error())
	}
	fmt.Println("display time:", cgm.ReadDisplayTime())
	fmt.Println("transmitter ID:", string(cgm.Cmd(dexcom.READ_TRANSMITTER_ID)))
	printXMLData("firmware header", cgm.ReadFirmwareHeader())

	hw := cgm.ReadXMLRecord(dexcom.MANUFACTURING_DATA)
	printXMLData("manufacturing data", hw.XML)
	fmt.Printf("    %+v\n", hw.Timestamp)

	sw := cgm.ReadXMLRecord(dexcom.PC_SOFTWARE_PARAMETER)
	printXMLData("PC software parameter", sw.XML)
	fmt.Printf("    %+v\n", sw.Timestamp)
}
