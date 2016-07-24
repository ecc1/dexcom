package dexcom

import (
	"bytes"
	"fmt"
)

// A PageType specifies a type of record stored by the Dexcom CGM receiver.
type PageType byte

//go:generate stringer -type PageType

const (
	MANUFACTURING_DATA      PageType = 0
	FIRMWARE_PARAMETER_DATA PageType = 1
	PC_SOFTWARE_PARAMETER   PageType = 2
	SENSOR_DATA             PageType = 3
	EGV_DATA                PageType = 4
	CAL_SET                 PageType = 5
	DEVIATION               PageType = 6
	INSERTION_TIME          PageType = 7
	RECEIVER_LOG_DATA       PageType = 8
	RECEIVER_ERROR_DATA     PageType = 9
	METER_DATA              PageType = 10
	USER_EVENT_DATA         PageType = 11
	USER_SETTING_DATA       PageType = 12

	FirstPageType = MANUFACTURING_DATA
	LastPageType  = USER_SETTING_DATA

	INVALID_PAGE PageType = 0xFF
)

// ReadPageRange returns the starting and ending page for a given PageType.
// The page numbers can be -1 if there are no entries (for example, USER_EVENT_DATA).
func (cgm *Cgm) ReadPageRange(pageType PageType) (int, int) {
	v := cgm.Cmd(READ_DATABASE_PAGE_RANGE, byte(pageType))
	if cgm.Error() != nil {
		return -1, -1
	}
	return int(UnmarshalInt32(v[:4])), int(UnmarshalInt32(v[4:]))
}

// The ReadPage function applies a function of type RecordFunc
// to each record that it reads, until it returns true or an error.
type RecordFunc func(Record) (bool, error)

type CrcError struct {
	Kind               string
	Received, Computed uint16
	PageType           PageType
	PageNumber         int
	Data               []byte
}

func (e CrcError) Error() string {
	if e.PageType == INVALID_PAGE {
		return fmt.Sprintf("bad %s CRC (received %02X, computed %02X); data = % X", e.Kind, e.Received, e.Computed, e.Data)
	}
	return fmt.Sprintf("bad %s CRC (received %02X, computed %02X) for %v page %d; data = % X", e.Kind, e.Received, e.Computed, e.PageType, e.PageNumber, e.Data)
}

// ReadPage reads the specified page and applies recordFn to each record.
// ReadPage returns true when an error is encountered or an invocation of
// recordFn returns true, otherwise it returns false.
func (cgm *Cgm) ReadPage(pageType PageType, pageNumber int, recordFn RecordFunc) bool {
	buf := bytes.Buffer{}
	buf.WriteByte(byte(pageType))
	buf.Write(MarshalInt32(int32(pageNumber)))
	buf.WriteByte(1)
	v := cgm.Cmd(READ_DATABASE_PAGES, buf.Bytes()...)
	if cgm.Error() != nil {
		return true
	}
	const headerSize = 28
	if len(v) < headerSize {
		cgm.SetError(fmt.Errorf("invalid page length (%d) for %v page %d", len(v), pageType, pageNumber))
		return true
	}
	crc := UnmarshalUint16(v[headerSize-2 : headerSize])
	calc := crc16(v[:headerSize-2])
	if crc != calc {
		cgm.SetError(CrcError{
			Kind:       "page",
			Received:   crc,
			Computed:   calc,
			PageType:   pageType,
			PageNumber: pageNumber,
			Data:       v,
		})
		return true
	}
	// firstIndex := int(UnmarshalInt32(v[0:4]))
	numRecords := int(UnmarshalInt32(v[4:8]))

	r := PageType(v[8])
	if r != pageType {
		cgm.SetError(fmt.Errorf("unexpected page type %d for %v page %d", r, pageType, pageNumber))
		return true
	}

	// rev := v[9]

	p := int(UnmarshalInt32(v[10:14]))
	if p != pageNumber {
		cgm.SetError(fmt.Errorf("unexpected page number %d for %v page %d", p, pageType, pageNumber))
		return true
	}

	// r1 := UnmarshalInt32(v[14:18])
	// r2 := UnmarshalInt32(v[18:22])
	// r3 := UnmarshalInt32(v[22:26])

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
	dataLen = ((dataLen + numRecords - 1) / numRecords) * numRecords
	recordLen := dataLen / numRecords
	data = data[:dataLen]

	// Slice data into records, validate per-record CRCs,
	// unmarshal record, and apply recordFn.
	// Iterate in reverse order to facilitate scanning for recent records.
	for i := numRecords - 1; i >= 0; i-- {
		rec := data[i*recordLen : (i+1)*recordLen]
		crc := UnmarshalUint16(rec[recordLen-2 : recordLen])
		rec = rec[:recordLen-2]
		calc := crc16(rec)
		if crc != calc {
			cgm.SetError(CrcError{
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
		r.Unmarshal(pageType, rec)
		done, err := recordFn(r)
		if err != nil || done {
			cgm.SetError(err)
			return true
		}
	}
	return false
}

func (cgm *Cgm) ReadRecords(pageType PageType, recordFn RecordFunc) {
	first, last := cgm.ReadPageRange(pageType)
	if cgm.Error() != nil {
		return
	}
	cgm.IterRecords(pageType, first, last, recordFn)
}

// IterRecords reads the specified page range and applies recordFn to each
// record in each page.  Pages are visited in reverse order to facilitate
// scanning for recent records.
func (cgm *Cgm) IterRecords(pageType PageType, firstPage, lastPage int, recordFn RecordFunc) {
	for n := lastPage; n >= firstPage; n-- {
		done := cgm.ReadPage(pageType, n, recordFn)
		if cgm.Error() != nil || done {
			return
		}
	}
}
