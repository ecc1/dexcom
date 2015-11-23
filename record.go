package dexcom

import (
	"fmt"
	"time"
)

// A SensorRecord contains a reading received from a Dexcom CGM sensor.
type SensorRecord struct {
	Timestamp  Timestamp
	Unfiltered uint32
	Filtered   uint32
	RSSI       uint16
}

// UnmarshalSensorRecord unmarshals a byte array into a SensorRecord.
func UnmarshalSensorRecord(v []byte) SensorRecord {
	if len(v) != 18 {
		panic(fmt.Sprintf("SensorRecord: wrong length (%d)", len(v)))
	}
	return SensorRecord{
		Timestamp:  UnmarshalTimestamp(v[0:8]),
		Unfiltered: UnmarshalUint32(v[8:12]),
		Filtered:   UnmarshalUint32(v[12:16]),
		RSSI:       UnmarshalUint16(v[16:18]),
	}
}

// SpecialGlucose values are used to encode various exceptional conditions.
type SpecialGlucose uint16

//go:generate stringer -type=SpecialGlucose

const (
	SENSOR_NOT_ACTIVE SpecialGlucose = 1 + iota
	MINIMAL_DEVIATION
	NO_ANTENNA
	_
	SENSOR_NOT_CALIBRATED
	COUNTS_DEVIATION
	_
	_
	ABSOLUTE_DEVIATION
	POWER_DEVIATION
	_
	BAD_RF
	specialLimit
)

// IsSpecial checks whether a glucose value falls in the SpecialGlucose range.
func IsSpecial(glucose uint16) bool {
	return glucose < uint16(specialLimit)
}

// The Trend type represents the directional arrows
// displayed by the Dexcom CGM receiver.
type Trend byte

//go:generate stringer -type=Trend

const (
	UP_UP Trend = 1 + iota
	UP
	UP_45
	FLAT
	DOWN_45
	DOWN
	DOWN_DOWN
	NOT_COMPUTABLE
	OUT_OF_RANGE
)

var trendSymbol = map[Trend]string{
	UP_UP:          "⇈",
	UP:             "↑",
	UP_45:          "↗",
	FLAT:           "→",
	DOWN_45:        "↘",
	DOWN:           "↓",
	DOWN_DOWN:      "⇊",
	NOT_COMPUTABLE: "⁇",
	OUT_OF_RANGE:   "⋯",
}

// Symbol returns the arrow symbol corresponding to a Trend value.
func (t Trend) Symbol() string {
	return trendSymbol[t]
}

// An EGVRecord contains a glucose reading calculated by the Dexcom CGM receiver.
type EGVRecord struct {
	Timestamp   Timestamp
	Glucose     uint16
	DisplayOnly bool
	Trend       Trend
}

const (
	EGV_DISPLAY_ONLY     = (1 << 15)
	EGV_VALUE_MASK       = (1 << 10) - 1
	EGV_TREND_ARROW_MASK = (1 << 4) - 1
)

// UnmarshalEGVRecord unmarshals a byte array into an EGVRecord.
func UnmarshalEGVRecord(v []byte) EGVRecord {
	if len(v) != 11 {
		panic(fmt.Sprintf("EGVRecord: wrong length (%d)", len(v)))
	}
	g := UnmarshalUint16(v[8:10])
	return EGVRecord{
		Timestamp:   UnmarshalTimestamp(v[0:8]),
		Glucose:     g & EGV_VALUE_MASK,
		DisplayOnly: (g & EGV_DISPLAY_ONLY) != 0,
		Trend:       Trend(v[10] & EGV_TREND_ARROW_MASK),
	}
}

// A MeterRecord contains a glucometer reading.
type MeterRecord struct {
	Timestamp Timestamp
	Glucose   uint16
	MeterTime time.Time
}

// UnmarshalMeterRecord unmarshals a byte array into a MeterRecord.
func UnmarshalMeterRecord(v []byte) MeterRecord {
	if len(v) != 14 {
		panic(fmt.Sprintf("MeterRecord: wrong length (%d)", len(v)))
	}
	return MeterRecord{
		Timestamp: UnmarshalTimestamp(v[0:8]),
		Glucose:   UnmarshalUint16(v[8:10]),
		MeterTime: UnmarshalTime(v[10:14]),
	}
}
