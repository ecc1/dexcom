package dexcom

import (
	"bytes"
	"encoding/xml"
)

// XMLInfo maps attribute names to values.
// The Dexcom CGM receiver represents its system data as single XML nodes
// with multiple attributes, so a tree structure is not required.
type XMLInfo map[string]string

func umarshalXMLInfo(r *Record, v []byte) {
	v = v[8:]
	i := bytes.IndexByte(v, 0x00)
	if i != -1 {
		v = v[:i]
	}
	r.Info = umarshalXMLBytes(v)
}

func umarshalXMLBytes(v []byte) XMLInfo {
	m := make(XMLInfo)
	err := xml.Unmarshal(v, &m)
	if err != nil {
		m["InvalidXML"] = string(v)
	}
	return m
}

// UnmarshalXML is called by xml.Unmarshal.
func (ptr *XMLInfo) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	m := *ptr
	for _, attr := range start.Attr {
		m[attr.Name.Local] = attr.Value
	}
	return d.Skip()
}

// ReadFirmwareHeader gets the firmware header from the Dexcom CGM receiver
// and returns it as XMLInfo.
func (cgm *CGM) ReadFirmwareHeader() XMLInfo {
	v := cgm.Cmd(ReadFirmwareHeader)
	if cgm.Error() != nil {
		return nil
	}
	return umarshalXMLBytes(v)
}

// ReadXMLRecord gets the given XML record type from the Dexcom CGM receiver.
func (cgm *CGM) ReadXMLRecord(pageType PageType) Record {
	x := Record{}
	proc := func(r Record) error {
		x = r
		return IterationDone
	}
	cgm.IterRecords(pageType, 0, 0, proc)
	return x
}
