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
	fmt.Println("display time:", cgm.ReadDisplayTime())
	fmt.Println("transmitter ID:", string(cgm.Cmd(dexcom.ReadTransmitterID)))
	printXMLInfo("firmware header", cgm.ReadFirmwareHeader())
	printXMLRecord(cgm, dexcom.ManufacturingData, "manufacturing data")
	printXMLRecord(cgm, dexcom.SoftwareData, "PC software parameter")
}

func printXMLInfo(name string, xml dexcom.XMLInfo) {
	fmt.Printf("%s:\n", name)
	for k, v := range xml {
		fmt.Printf("    %s: %s\n", k, v)
	}
}

func printXMLRecord(cgm *dexcom.CGM, pageType dexcom.PageType, description string) {
	r := cgm.ReadXMLRecord(pageType)
	xml := r.Info.(dexcom.XMLInfo)
	printXMLInfo(description, xml)
	fmt.Printf("    %+v\n", r.Timestamp)
}
