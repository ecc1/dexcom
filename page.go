package dexcom

import (
	"bytes"
	"fmt"
)

// A RecordType specifies a type of record stored by the Dexcom CGM receiver.
type RecordType byte

//go:generate stringer -type RecordType

const (
	MANUFACTURING_DATA      RecordType = 0
	FIRMWARE_PARAMETER_DATA RecordType = 1
	PC_SOFTWARE_PARAMETER   RecordType = 2
	SENSOR_DATA             RecordType = 3
	EGV_DATA                RecordType = 4
	CAL_SET                 RecordType = 5
	DEVIATION               RecordType = 6
	INSERTION_TIME          RecordType = 7
	RECEIVER_LOG_DATA       RecordType = 8
	RECEIVER_ERROR_DATA     RecordType = 9
	METER_DATA              RecordType = 10
	USER_EVENT_DATA         RecordType = 11
	USER_SETTING_DATA       RecordType = 12

	FirstRecordType = MANUFACTURING_DATA
	LastRecordType  = USER_SETTING_DATA

	// internal record types for unmarshalling
	timestampType
	xmlDataType
)

// A RecordContext holds information about a range of pages of
// a given record type during iteration over those records.
type RecordContext struct {
	RecordType RecordType
	StartPage  int
	EndPage    int
	PageNumber int
	Index      int
}

// ReadPageRange requests the StartPage and EndPage for a given RecordType
// and returns a RecordContext with those values.  The page numbers
// can be -1 if there are no entries (for example, USER_EVENT_DATA).
func (cgm *Cgm) ReadPageRange(recordType RecordType) RecordContext {
	v := cgm.Cmd(READ_DATABASE_PAGE_RANGE, byte(recordType))
	if cgm.Error() != nil {
		return RecordContext{}
	}
	return RecordContext{
		RecordType: recordType,
		StartPage:  int(UnmarshalInt32(v[:4])),
		EndPage:    int(UnmarshalInt32(v[4:])),
	}
}

// The ReadPage function applies a function of type RecordFunc
// to each record that it reads, until it returns false or an error.
type RecordFunc func([]byte, RecordContext) (bool, error)

type CrcError struct {
	Kind               string
	Received, Computed uint16
	Context            *RecordContext
	Data               []byte
}

func (e CrcError) Error() string {
	if e.Context != nil {
		return fmt.Sprintf("bad %s CRC (received %02X, computed %02X) in context %+v; data = % X", e.Kind, e.Received, e.Computed, *e.Context, e.Data)
	} else {
		return fmt.Sprintf("bad %s CRC (received %02X, computed %02X); data = % X", e.Kind, e.Received, e.Computed, e.Data)
	}
}

// ReadPage reads a single page specified by the PageNumber field of the
// given RecordContext and applies recordFn to each record in the page.
func (cgm *Cgm) ReadPage(context RecordContext, recordFn RecordFunc) bool {
	buf := bytes.Buffer{}
	buf.WriteByte(byte(context.RecordType))
	buf.Write(MarshalInt32(int32(context.PageNumber)))
	buf.WriteByte(1)
	v := cgm.Cmd(READ_DATABASE_PAGES, buf.Bytes()...)
	if cgm.Error() != nil {
		return false
	}
	const headerSize = 28
	if len(v) < headerSize {
		cgm.SetError(fmt.Errorf("invalid page length (%d)", len(v)))
		return false
	}
	crc := UnmarshalUint16(v[headerSize-2 : headerSize])
	calc := crc16(v[:headerSize-2])
	if crc != calc {
		cgm.SetError(CrcError{
			Kind:     "page",
			Received: crc,
			Computed: calc,
			Context:  &context,
			Data:     v,
		})
		return false
	}
	firstIndex := int(UnmarshalInt32(v[0:4]))
	numRecords := int(UnmarshalInt32(v[4:8]))

	r := RecordType(v[8])
	if r != context.RecordType {
		cgm.SetError(fmt.Errorf("unexpected record type %d in context %+v", r, context))
		return false
	}

	// rev := v[9]

	p := int(UnmarshalInt32(v[10:14]))
	if p != context.PageNumber {
		cgm.SetError(fmt.Errorf("unexpected page number %d in context %+v", p, context))
		return false
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

	// Slice data into records, validate per-record CRCs, and apply recordFn.
	// Iterate in reverse order to facilitate scanning for recent records.
	for i := numRecords - 1; i >= 0; i-- {
		context.Index = firstIndex + i
		rec := data[i*recordLen : (i+1)*recordLen]
		crc := UnmarshalUint16(rec[recordLen-2 : recordLen])
		rec = rec[:recordLen-2]
		calc := crc16(rec)
		if crc != calc {
			cgm.SetError(CrcError{
				Kind:     "record",
				Received: crc,
				Computed: calc,
				Context:  &context,
				Data:     rec,
			})
			return false
		}
		keepGoing, err := recordFn(rec, context)
		if err != nil || !keepGoing {
			cgm.SetError(err)
			return false
		}
	}
	return true
}

func (cgm *Cgm) ReadRecords(recordType RecordType, recordFn RecordFunc) {
	context := cgm.ReadPageRange(recordType)
	if cgm.Error() != nil {
		return
	}
	cgm.IterRecords(context, recordFn)
}

// IterRecords reads the pages of the type and range specified by the
// given RecordContext and applies recordFn to each record in each page.
// Pages are visited in reverse order to facilitate scanning for recent records.
func (cgm *Cgm) IterRecords(context RecordContext, recordFn RecordFunc) {
	for n := context.EndPage; n >= context.StartPage; n-- {
		context.PageNumber = n
		keepGoing := cgm.ReadPage(context, recordFn)
		if cgm.Error() != nil || !keepGoing {
			return
		}
	}
}
