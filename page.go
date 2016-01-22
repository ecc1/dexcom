package dexcom

import (
	"bytes"
	"fmt"
)

// A RecordType specifies a type of record stored by the Dexcom CGM receiver.
type RecordType byte

//go:generate stringer -type=RecordType

const (
	MANUFACTURING_DATA RecordType = iota
	FIRMWARE_PARAMETER_DATA
	PC_SOFTWARE_PARAMETER
	SENSOR_DATA
	EGV_DATA
	CAL_SET
	DEVIATION
	INSERTION_TIME
	RECEIVER_LOG_DATA
	RECEIVER_ERROR_DATA
	METER_DATA
	USER_EVENT_DATA
	USER_SETTING_DATA
	FirstRecordType = MANUFACTURING_DATA
	LastRecordType  = USER_SETTING_DATA
)

// A RecordContext holds information about a range of pages of
// a given record type during iteration over those records.
type RecordContext struct {
	RecordType RecordType
	StartPage  uint32
	EndPage    uint32
	PageNumber uint32
	Index      uint32
}

// ReadPageRange requests the StartPage and EndPage for a given RecordType
// and returns a RecordContext with those values.
func ReadPageRange(recordType RecordType) (*RecordContext, error) {
	v, err := Cmd(READ_DATABASE_PAGE_RANGE, []byte{byte(recordType)})
	if err != nil {
		return nil, err
	}
	context := RecordContext{
		RecordType: recordType,
		StartPage:  UnmarshalUint32(v[:4]),
		EndPage:    UnmarshalUint32(v[4:]),
	}
	return &context, nil
}

// RecordFunc is the type signature of the procedure that is called
// for each record that is read by ReadPage.
type RecordFunc func(record []byte, context *RecordContext) error

func ReadRecords(recordType RecordType, recordFn RecordFunc) error {
	context, err := ReadPageRange(recordType)
	if err != nil {
		return err
	}
	return IterRecords(context, recordFn)
}

// IterRecords reads the pages of the type and range specified by the
// given RecordContext and applies recordFn to each record in each page.
func IterRecords(context *RecordContext, recordFn RecordFunc) error {
	for n := context.StartPage; n <= context.EndPage; n++ {
		context.PageNumber = n
		err := ReadPage(context, recordFn)
		if err != nil {
			return err
		}
	}
	return nil
}

// ReadPage reads a single page specified by the PageNumber field of the
// given RecordContext and applies recordFn to each record in the page.
func ReadPage(context *RecordContext, recordFn RecordFunc) error {
	v, err := Cmd(READ_DATABASE_PAGES, []byte{byte(context.RecordType)}, MarshalUint32(context.PageNumber), []byte{1})
	if err != nil {
		return err
	}

	const headerSize = 28
	if len(v) < headerSize {
		return fmt.Errorf("invalid page length (%d)", len(v))
	}
	crc := v[26:28]
	calc := crc16(v[:26])
	if !bytes.Equal(crc, calc) {
		return fmt.Errorf("bad page CRC (received %X, computed %X) in context %v", crc, calc, *context)
	}

	firstIndex := UnmarshalUint32(v[0:4])
	numRecords := int(UnmarshalUint32(v[4:8]))

	r := RecordType(v[8])
	if r != context.RecordType {
		return fmt.Errorf("unexpected record type %X in context %v", r, *context)
	}

	// rev := v[9]

	p := UnmarshalUint32(v[10:14])
	if p != context.PageNumber {
		return fmt.Errorf("unexpected page number %X in context %v", p, *context)
	}

	// r1 := UnmarshalUint32(v[14:18])
	// r2 := UnmarshalUint32(v[18:22])
	// r3 := UnmarshalUint32(v[22:26])

	data := v[headerSize:]
	dataLen := len(data)

	// remove padding (trailing 0xFF bytes) and compute record length
	for i := 1; i <= dataLen; i++ {
		if data[dataLen-i] != 0xFF {
			dataLen -= i - 1
			break
		}
	}
	data = data[:dataLen]
	recordLen := dataLen / numRecords

	// slice data into records, validate per-record CRCs, and apply recordFn
	for i := 0; i < numRecords; i++ {
		context.Index = firstIndex + uint32(i)
		rec := data[i*recordLen : (i+1)*recordLen]
		crc := rec[recordLen-2:]
		rec = rec[:recordLen-2]
		calc := crc16(rec)
		if !bytes.Equal(crc, calc) {
			return fmt.Errorf("bad record CRC (received %X, computed %X) in context %v", crc, calc, *context)
		}
		err = recordFn(rec, context)
		if err != nil {
			return err
		}
	}
	return nil
}
