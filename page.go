package dexcom

import (
	"bytes"
	"fmt"
)

// A PageType specifies a type of record stored by the Dexcom G4 receiver.
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

// ReadPageRange returns the starting and ending page for a given PageType.
// The page numbers can be -1 if there are no entries (for example, USER_EVENT_DATA).
func (cgm *CGM) ReadPageRange(pageType PageType) (int, int) {
	v := cgm.Cmd(ReadDatabasePageRange, byte(pageType))
	if cgm.Error() != nil {
		return -1, -1
	}
	return int(unmarshalInt32(v[:4])), int(unmarshalInt32(v[4:]))
}

// RecordFunc represents a function that ReadRecords applies
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

// ReadRecords reads the specified page and applies recordFn to each record.
// ReadRecords returns true when an error is encountered or an invocation of
// recordFn returns true, otherwise it returns false.
func (cgm *CGM) ReadRecords(pageType PageType, pageNumber int, recordFn RecordFunc) bool {
	v := cgm.ReadPage(pageType, pageNumber)
	if cgm.Error() != nil {
		return true
	}
	const headerSize = 28
	if len(v) < headerSize {
		cgm.SetError(fmt.Errorf("invalid page length (%d) for %v page %d", len(v), pageType, pageNumber))
		return true
	}
	crc := unmarshalUint16(v[headerSize-2 : headerSize])
	calc := crc16(v[:headerSize-2])
	if crc != calc {
		cgm.SetError(CRCError{
			Kind:       "page",
			Received:   crc,
			Computed:   calc,
			PageType:   pageType,
			PageNumber: pageNumber,
			Data:       v,
		})
		return true
	}
	// firstIndex := int(unmarshalInt32(v[0:4]))
	numRecords := int(unmarshalInt32(v[4:8]))

	r := PageType(v[8])
	if r != pageType {
		cgm.SetError(fmt.Errorf("unexpected page type %d for %v page %d", r, pageType, pageNumber))
		return true
	}

	// rev := v[9]

	p := int(unmarshalInt32(v[10:14]))
	if p != pageNumber {
		cgm.SetError(fmt.Errorf("unexpected page number %d for %v page %d", p, pageType, pageNumber))
		return true
	}

	// r1 := unmarshalInt32(v[14:18])
	// r2 := unmarshalInt32(v[18:22])
	// r3 := unmarshalInt32(v[22:26])

	data := v[headerSize:]
	dataLen := len(data)

	// Remove padding (trailing 0xFF bytes) and compute record length.
	for dataLen > 0 {
		if data[dataLen-1] != 0xFF {
			break
		}
		dataLen--
	}
	// Round dataLen up to a multiple of numRecords so we keep
	// any 0xFF bytes that are part of the last record.
	dataLen = (dataLen + numRecords - 1) / numRecords * numRecords
	recordLen := dataLen / numRecords
	data = data[:dataLen]

	return cgm.processRecords(pageType, pageNumber, recordFn, data, recordLen, numRecords)
}

// Slice data into records, validate per-record CRCs, unmarshal record,
// and apply recordFn.
// Iterate in reverse order to facilitate scanning for recent records.
func (cgm *CGM) processRecords(pageType PageType, pageNumber int, recordFn RecordFunc, data []byte, recordLen int, numRecords int) bool {
	for i := numRecords - 1; i >= 0; i-- {
		rec := data[i*recordLen : (i+1)*recordLen]
		crc := unmarshalUint16(rec[recordLen-2 : recordLen])
		rec = rec[:recordLen-2]
		calc := crc16(rec)
		if crc != calc {
			cgm.SetError(CRCError{
				Kind:       "record",
				Received:   crc,
				Computed:   calc,
				PageType:   pageType,
				PageNumber: pageNumber,
				Data:       rec,
			})
			return true
		}
		r := Record{}
		done := false
		err := r.unmarshal(pageType, rec)
		if err == nil {
			done, err = recordFn(r)
		}
		if err != nil || done {
			cgm.SetError(err)
			return true
		}
	}
	return false
}

// IterRecords reads the specified page range and applies recordFn to each
// record in each page.  Pages are visited in reverse order to facilitate
// scanning for recent records.
func (cgm *CGM) IterRecords(pageType PageType, firstPage, lastPage int, recordFn RecordFunc) {
	for n := lastPage; n >= firstPage; n-- {
		done := cgm.ReadRecords(pageType, n, recordFn)
		if cgm.Error() != nil || done {
			return
		}
	}
}
