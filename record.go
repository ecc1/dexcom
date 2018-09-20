package dexcom

import (
	"bytes"
	"fmt"
	"time"
)

type (
	// Record represents a time-stamped Dexcom receiver record.
	Record struct {
		Timestamp   Timestamp
		XML         XMLInfo          `json:",omitempty"`
		Sensor      *SensorInfo      `json:",omitempty"`
		EGV         *EGVInfo         `json:",omitempty"`
		Calibration *CalibrationInfo `json:",omitempty"`
		Insertion   *InsertionInfo   `json:",omitempty"`
		Meter       *MeterInfo       `json:",omitempty"`
	}

	// Records represents a sequence of records.
	Records []Record

	// SensorInfo represents a sensor reading.
	SensorInfo struct {
		Unfiltered uint32
		Filtered   uint32
		RSSI       int8
		Unknown    byte
	}

	// EGVInfo represents an estimated glucose value.
	EGVInfo struct {
		Glucose     uint16
		DisplayOnly bool
		Noise       uint8
		Trend       Trend
	}

	// InsertionInfo represents a sensor change event.
	InsertionInfo struct {
		SystemTime time.Time
		Event      SensorChange
	}

	// MeterInfo represents a meter reading.
	MeterInfo struct {
		Glucose   uint16
		MeterTime time.Time
	}

	// CalibrationInfo represents a calibration event.
	CalibrationInfo struct {
		Slope     float64
		Intercept float64
		Scale     float64
		Decay     float64
		Data      []CalibrationRecord
	}

	// CalibrationRecord represents a calibration data point.
	CalibrationRecord struct {
		TimeEntered time.Time
		Glucose     int32
		Raw         int32
		TimeApplied time.Time
	}
)

// Time returns the record's display time.
func (r Record) Time() time.Time {
	return r.Timestamp.DisplayTime
}

// Glucose returns the glucose field from an EGV record.
func (r Record) Glucose() uint16 {
	return r.EGV.Glucose
}

// Len returns the number of records.
func (v Records) Len() int {
	return len(v)
}

// Time returns the time of the record at index i.
func (v Records) Time(i int) time.Time {
	return v[i].Time()
}

var recordUnmarshal = map[PageType]func(*Record, []byte){
	ManufacturingData: umarshalXMLInfo,
	FirmwareData:      umarshalXMLInfo,
	SoftwareData:      umarshalXMLInfo,
	SensorData:        unmarshalSensorInfo,
	EGVData:           umarshalEGVInfo,
	CalibrationData:   unmarshalCalibrationInfo,
	InsertionTimeData: unmarshalInsertionInfo,
	MeterData:         unmarshalMeterInfo,
}

func (r *Record) unmarshal(pageType PageType, v []byte) error {
	f, found := recordUnmarshal[pageType]
	if !found {
		return fmt.Errorf("unmarshaling of %v records is unimplemented: % X", pageType, v)
	}
	r.Timestamp.unmarshal(v[0:8])
	f(r, v)
	return nil
}

func unmarshalSensorInfo(r *Record, v []byte) {
	r.Sensor = &SensorInfo{
		Unfiltered: unmarshalUint32(v[8:12]),
		Filtered:   unmarshalUint32(v[12:16]),
		RSSI:       int8(v[16]),
		Unknown:    v[17],
	}
}

// SpecialGlucose represents a glucose value that indicates an exceptional condition.
type SpecialGlucose uint16

//go:generate stringer -type SpecialGlucose

// Exceptional conditions.
const (
	SensorNotActive     SpecialGlucose = 1
	MinimalDeviation    SpecialGlucose = 2
	NoAntenna           SpecialGlucose = 3
	SensorNotCalibrated SpecialGlucose = 5
	CountDeviation      SpecialGlucose = 6
	AbsoluteDeviation   SpecialGlucose = 9
	PowerDeviation      SpecialGlucose = 10
	BadRF               SpecialGlucose = 12

	specialLimit = BadRF
)

// IsSpecial checks whether a glucose value falls in the SpecialGlucose range.
func IsSpecial(glucose uint16) bool {
	return glucose <= uint16(specialLimit)
}

// Trend represents a directional arrow displayed by the Dexcom CGM receiver.
type Trend byte

//go:generate stringer -type Trend

// Trend arrows.
const (
	UpUp          Trend = 1
	Up            Trend = 2
	Up45          Trend = 3
	Flat          Trend = 4
	Down45        Trend = 5
	Down          Trend = 6
	DownDown      Trend = 7
	NotComputable Trend = 8
	OutOfRange    Trend = 9
)

var trendSymbol = map[Trend]string{
	UpUp:          "⇈",
	Up:            "↑",
	Up45:          "↗",
	Flat:          "→",
	Down45:        "↘",
	Down:          "↓",
	DownDown:      "⇊",
	NotComputable: "⁇",
	OutOfRange:    "⋯",
}

// Symbol converts a Trend to a graphical representation.
func (t Trend) Symbol() string {
	return trendSymbol[t]
}

// Constants used to extract EGV, noise, and trend.
const (
	EGVDisplayOnly = 1 << 15
	EGVValueMask   = 0x3FF
	EGVNoiseMask   = 0x70
	EGVTrendMask   = 0xF
)

func umarshalEGVInfo(r *Record, v []byte) {
	g := unmarshalUint16(v[8:10])
	r.EGV = &EGVInfo{
		Glucose:     g & EGVValueMask,
		DisplayOnly: g&EGVDisplayOnly != 0,
		Noise:       v[10] & EGVNoiseMask >> 4,
		Trend:       Trend(v[10] & EGVTrendMask),
	}
}

func unmarshalCalibrationInfo(r *Record, v []byte) {
	cal := &CalibrationInfo{
		Slope:     unmarshalFloat64(v[8:16]),
		Intercept: unmarshalFloat64(v[16:24]),
		Scale:     unmarshalFloat64(v[24:32]),
		Decay:     unmarshalFloat64(v[35:43]),
	}
	n := int(v[43])
	cal.Data = make([]CalibrationRecord, n)
	v = v[44:]
	offset := r.Timestamp.DisplayTime.Sub(r.Timestamp.SystemTime)
	for i := 0; i < n; i++ {
		cal.Data[i].unmarshal(v)
		cal.Data[i].TimeEntered = cal.Data[i].TimeEntered.Add(offset)
		cal.Data[i].TimeApplied = cal.Data[i].TimeApplied.Add(offset)
		v = v[17:]
	}
	r.Calibration = cal
}

func (r *CalibrationRecord) unmarshal(v []byte) {
	r.TimeEntered = unmarshalTime(v[0:4])
	r.Glucose = unmarshalInt32(v[4:8])
	r.Raw = unmarshalInt32(v[8:12])
	r.TimeApplied = unmarshalTime(v[12:16])
}

// SensorChange represents a sensor change.
type SensorChange byte

//go:generate stringer -type SensorChange

// Sensor change values.
const (
	Stopped SensorChange = 1
	Started SensorChange = 7
)

var (
	invalidTime = []byte{0xFF, 0xFF, 0xFF, 0xFF}
)

func unmarshalInsertionInfo(r *Record, v []byte) {
	t := time.Time{}
	u := v[8:12]
	if !bytes.Equal(u, invalidTime) {
		t = unmarshalTime(u)
	}
	r.Insertion = &InsertionInfo{
		SystemTime: t,
		Event:      SensorChange(v[12]),
	}
}

func unmarshalMeterInfo(r *Record, v []byte) {
	r.Meter = &MeterInfo{
		Glucose:   unmarshalUint16(v[8:10]),
		MeterTime: unmarshalTime(v[10:14]),
	}
}
