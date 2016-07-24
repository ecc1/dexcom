package main

import (
	"fmt"
	"log"

	"github.com/ecc1/dexcom"
)

func printXmlInfo(name string, xml *dexcom.XmlInfo) {
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
	printXmlInfo("firmware header", cgm.ReadFirmwareHeader())

	hw := cgm.ReadXMLRecord(dexcom.MANUFACTURING_DATA)
	printXmlInfo("manufacturing data", hw.Xml)
	fmt.Printf("    %+v\n", hw.Timestamp)

	sw := cgm.ReadXMLRecord(dexcom.PC_SOFTWARE_PARAMETER)
	printXmlInfo("PC software parameter", sw.Xml)
	fmt.Printf("    %+v\n", sw.Timestamp)
}
