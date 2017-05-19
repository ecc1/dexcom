package dexcom

import (
	"bytes"
	"fmt"
)

// PageType specifies a record page type stored by the Dexcom G4 receiver.
type PageType byte

//go:generate stringer -type PageType

// Types of CGM records stored by the Dexcom G4 receiver.
const (
	ManufacturingData PageType = 0
	FirmwareData      PageType = 1
	SoftwareData      PageType = 2
	SensorData        PageType = 3
	EGVData           PageType = 4
	CalibrationData   PageType = 5
	DeviationData     PageType = 6
	InsertionTimeData PageType = 7
	ReceiverLogData   PageType = 8
	ReceiverErrorData PageType = 9
	MeterData         PageType = 10
	UserEventData     PageType = 11
	UserSettingData   PageType = 12

	FirstPageType = ManufacturingData
	LastPageType  = UserSettingData

	InvalidPage PageType = 0xFF
)

// Lengths of fixed-size records (including 2-byte CRC),
// or 0 if veriable-length (terminated by null bytes).
var recordLength = map[PageType]int{
	ManufacturingData: 0,
	FirmwareData:      0,
	SoftwareData:      0,
	SensorData:        20,
	EGVData:           13,
	CalibrationData:   249,
	InsertionTimeData: 15,
	MeterData:         16,
}

// ReadPageRange returns the starting and ending page for a given PageType.
// The page numbers can be -1 if there are no entries (for example, USER_EVENT_DATA).
func (cgm *CGM) ReadPageRange(pageType PageType) (int, int) {
	v := cgm.Cmd(ReadDatabasePageRange, byte(pageType))
	if cgm.Error() != nil {
		return -1, -1
	}
	return int(unmarshalInt32(v[:4])), int(unmarshalInt32(v[4:]))
}

// RecordFunc represents a function that IterRecords applies
// to each record that it reads, until it returns true or an error.
type RecordFunc func(Record) (bool, error)

// CRCError indicates that a CRC error was detected.
type CRCError struct {
	Kind               string
	Received, Computed uint16
	PageType           PageType
	PageNumber         int
	Data               []byte
}

func (e CRCError) Error() string {
	if e.PageType == InvalidPage {
		return fmt.Sprintf("bad %s CRC (received %02X, computed %02X); data = % X", e.Kind, e.Received, e.Computed, e.Data)
	}
	return fmt.Sprintf("bad %s CRC (received %02X, computed %02X) for %v page %d; data = % X", e.Kind, e.Received, e.Computed, e.PageType, e.PageNumber, e.Data)
}

// ReadPage reads the specified page.
func (cgm *CGM) ReadPage(pageType PageType, pageNumber int) []byte {
	buf := bytes.Buffer{}
	buf.WriteByte(byte(pageType))
	buf.Write(marshalInt32(int32(pageNumber)))
	buf.WriteByte(1)
	return cgm.Cmd(ReadDatabasePages, buf.Bytes()...)
}

// ReadRawRecords reads the specified page and returns its records as raw byte slices.
func (cgm *CGM) ReadRawRecords(pageType PageType, pageNumber int) [][]byte {
	v := cgm.ReadPage(pageType, pageNumber)
	if cgm.Error() != nil {
		return nil
	}
	p, n, data, err := unmarshalPage(v)
	if err != nil {
		cgm.SetError(fmt.Errorf("%v page %d: %v", pageType, pageNumber, err))
		return nil
	}
	if p != pageType {
		cgm.SetError(fmt.Errorf("%v page %d: unexpected page type (%d)", pageType, pageNumber, p))
		return nil
	}
	if n != pageNumber {
		cgm.SetError(fmt.Errorf("%v page %d: unexpected page number (%d)", pageType, pageNumber, n))
		return nil
	}
	return data
}

// ReadRecords reads the specified page and returns its records.
func (cgm *CGM) ReadRecords(pageType PageType, pageNumber int) []Record {
	data := cgm.ReadRawRecords(pageType, pageNumber)
	if cgm.Error() != nil {
		return nil
	}
	records, err := unmarshalRecords(pageType, data)
	if err != nil {
		cgm.SetError(fmt.Errorf("%v page %d: %v", pageType, pageNumber, err))
		return nil
	}
	return records
}

// unmarshalPage returns the page type, page number, and raw records in the given page data.
func unmarshalPage(v []byte) (pageType PageType, pageNumber int, records [][]byte, err error) {
	const headerSize = 28
	if len(v) < headerSize {
		err = fmt.Errorf("invalid page length (%d)", len(v))
		return
	}
	crc := unmarshalUint16(v[headerSize-2 : headerSize])
	calc := crc16(v[:headerSize-2])
	if crc != calc {
		err = CRCError{
			Kind:       "page",
			Received:   crc,
			Computed:   calc,
			PageType:   pageType,
			PageNumber: pageNumber,
			Data:       v,
		}
		return
	}
	// firstIndex := int(unmarshalInt32(v[0:4]))
	numRecords := int(unmarshalInt32(v[4:8]))
	pageType = PageType(v[8])
	// rev := v[9]
	pageNumber = int(unmarshalInt32(v[10:14]))
	// r1 := unmarshalInt32(v[14:18])
	// r2 := unmarshalInt32(v[18:22])
	// r3 := unmarshalInt32(v[22:26])
	v = v[headerSize:]
	recordLen := recordLength[pageType]
	if recordLen == 0 {
		if numRecords != 1 {
			err = fmt.Errorf("unexpected number of records (%d)", numRecords)
			return
		}
		recordLen = len(v)
	}
	records = make([][]byte, 0, numRecords)
	// Collect records in reverse chronological order.
	for i := numRecords - 1; i >= 0; i-- {
		rec := v[i*recordLen : (i+1)*recordLen]
		crc := unmarshalUint16(rec[recordLen-2 : recordLen])
		rec = rec[:recordLen-2]
		calc := crc16(rec)
		if crc != calc {
			err = CRCError{
				Kind:       "record",
				Received:   crc,
				Computed:   calc,
				PageType:   pageType,
				PageNumber: pageNumber,
				Data:       rec,
			}
			return
		}
		records = append(records, rec)
	}
	return
}

func unmarshalRecords(pageType PageType, data [][]byte) ([]Record, error) {
	records := make([]Record, 0, len(data))
	for _, rec := range data {
		r := Record{}
		err := r.unmarshal(pageType, rec)
		if err != nil {
			return records, err
		}
		records = append(records, r)
	}
	return records, nil
}

// IterRecords reads the specified page range and applies recordFn to each
// record in each page.  Pages are visited in reverse order to facilitate
// scanning for recent records.
func (cgm *CGM) IterRecords(pageType PageType, firstPage, lastPage int, recordFn RecordFunc) {
	for n := lastPage; n >= firstPage; n-- {
		records := cgm.ReadRecords(pageType, n)
		if cgm.Error() != nil {
			return
		}
		for _, r := range records {
			done, err := recordFn(r)
			if err != nil {
				cgm.SetError(fmt.Errorf("%v page %d: %v", pageType, n, err))

				return
			}
			if done {
				return
			}
		}
	}
}
