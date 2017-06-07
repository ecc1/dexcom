package dexcom

import (
	"bytes"
	"errors"
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

type pageInfo struct {
	Type    PageType
	Number  int
	Records [][]byte
}

// ReadRawRecords reads the specified page and returns its records as raw byte slices.
func (cgm *CGM) ReadRawRecords(pageType PageType, pageNumber int) [][]byte {
	v := cgm.ReadPage(pageType, pageNumber)
	if cgm.Error() != nil {
		return nil
	}
	page, err := unmarshalPage(v)
	if err != nil {
		if page == nil {
			cgm.SetError(err)
			return nil
		}
		err = fmt.Errorf("%v page %d: %v", page.Type, page.Number, err)
	} else if page.Type != pageType {
		err = fmt.Errorf("%v page %d: unexpected page type (%d)", pageType, pageNumber, page.Type)
	} else if page.Number != pageNumber {
		err = fmt.Errorf("%v page %d: unexpected page number (%d)", pageType, pageNumber, page.Number)
	}
	cgm.SetError(err)
	return page.Records
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
	}
	return records
}

// unmarshalPage validates the CRC of the given page data and
// uses the page type to slice the data into raw records.
func unmarshalPage(v []byte) (*pageInfo, error) {
	const headerSize = 28
	if len(v) < headerSize {
		return nil, fmt.Errorf("invalid page length (%d)", len(v))
	}
	crc := unmarshalUint16(v[headerSize-2 : headerSize])
	calc := crc16(v[:headerSize-2])
	if crc != calc {
		return nil, CRCError{
			Kind:     "page",
			Received: crc,
			Computed: calc,
			Data:     v,
		}
	}
	// firstIndex := int(unmarshalInt32(v[0:4]))
	numRecords := int(unmarshalInt32(v[4:8]))
	pageType := PageType(v[8])
	// rev := v[9]
	pageNumber := int(unmarshalInt32(v[10:14]))
	// r1 := unmarshalInt32(v[14:18])
	// r2 := unmarshalInt32(v[18:22])
	// r3 := unmarshalInt32(v[22:26])
	page := pageInfo{Type: pageType, Number: pageNumber}
	v = v[headerSize:]
	recordLen := recordLength[pageType]
	if recordLen == 0 {
		if numRecords != 1 {
			return &page, fmt.Errorf("unexpected number of records (%d)", numRecords)
		}
		recordLen = len(v)
	}
	page.Records = make([][]byte, 0, numRecords)
	// Collect records in reverse chronological order.
	for i := numRecords - 1; i >= 0; i-- {
		rec := v[i*recordLen : (i+1)*recordLen]
		crc := unmarshalUint16(rec[recordLen-2 : recordLen])
		rec = rec[:recordLen-2]
		calc := crc16(rec)
		if crc != calc {
			return &page, CRCError{
				Kind:       "record",
				Received:   crc,
				Computed:   calc,
				PageType:   pageType,
				PageNumber: pageNumber,
				Data:       rec,
			}
		}
		page.Records = append(page.Records, rec)
	}
	return &page, nil
}

func unmarshalRecords(pageType PageType, data [][]byte) ([]Record, error) {
	records := make([]Record, 0, len(data))
	var err error
	for _, rec := range data {
		r := Record{}
		err = r.unmarshal(pageType, rec)
		if err != nil {
			break
		}
		records = append(records, r)
	}
	return records, err
}

// RecordFunc represents a function that IterRecords applies to each record.
type RecordFunc func(Record) error

// IterationDone can be returned by a RecordFunc to indicate that
// iteration is complete and no further records need to be processed.
var IterationDone = errors.New("iteration done") // nolint

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
			err := recordFn(r)
			if err != nil {
				if err != IterationDone {
					cgm.SetError(fmt.Errorf("%v page %d: %v", pageType, n, err))
				}
				return
			}
		}
	}
}
