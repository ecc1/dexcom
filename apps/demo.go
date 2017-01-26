package main

import (
	"fmt"
	"log"

	"github.com/ecc1/dexcom"
)

func printXMLInfo(name string, xml *dexcom.XMLInfo) {
	fmt.Printf("%s:\n", name)
	for k, v := range *xml {
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
	printXMLInfo("firmware header", cgm.ReadFirmwareHeader())

	hw := cgm.ReadXMLRecord(dexcom.MANUFACTURING_DATA)
	printXMLInfo("manufacturing data", hw.XML)
	fmt.Printf("    %+v\n", hw.Timestamp)

	sw := cgm.ReadXMLRecord(dexcom.PC_SOFTWARE_PARAMETER)
	printXMLInfo("PC software parameter", sw.XML)
	fmt.Printf("    %+v\n", sw.Timestamp)
}
