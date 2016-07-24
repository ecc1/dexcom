package dexcom

import (
	"fmt"
	"time"
)

type (
	Record struct {
		Timestamp   Timestamp
		Sensor      *SensorInfo      `json:",omitempty"`
		Egv         *EgvInfo         `json:",omitempty"`
		Calibration *CalibrationInfo `json:",omitempty"`
		Meter       *MeterInfo       `json:",omitempty"`
		Xml         *XmlInfo         `json:",omitempty"`
	}

	SensorInfo struct {
		Unfiltered uint32
		Filtered   uint32
		Rssi       uint16
	}

	EgvInfo struct {
		Glucose     uint16
		DisplayOnly bool
		Trend       Trend
	}

	MeterInfo struct {
		Glucose   uint16
		MeterTime time.Time
	}

	CalibrationInfo struct {
		Slope     float64
		Intercept float64
		Scale     float64
		Decay     float64
		Data      []CalibrationData
	}

	CalibrationData struct {
		TimeEntered time.Time
		Glucose     int32
		Raw         int32
		TimeApplied time.Time
	}
)

var recordUnmarshal = map[PageType]struct {
	length    int
	unmarshal func(*Record, []byte)
}{
	MANUFACTURING_DATA:      {-1, unmarshalXmlInfo},
	FIRMWARE_PARAMETER_DATA: {-1, unmarshalXmlInfo},
	PC_SOFTWARE_PARAMETER:   {-1, unmarshalXmlInfo},
	SENSOR_DATA:             {18, unmarshalSensorInfo},
	EGV_DATA:                {11, unmarshalEgvInfo},
	CAL_SET:                 {-1, unmarshalCalibrationInfo},
	METER_DATA:              {14, unmarshalMeterInfo},
}

func (r *Record) Unmarshal(pageType PageType, v []byte) error {
	u, found := recordUnmarshal[pageType]
	if !found {
		return fmt.Errorf("unmarshaling of %v records is unimplemented", pageType)
	}
	if u.length > 0 && len(v) != u.length {
		return fmt.Errorf("wrong length (%d) for %d-byte %v record", len(v), u.length, pageType)
	}
	r.Timestamp.Unmarshal(v[0:8])
	u.unmarshal(r, v)
	return nil
}

func unmarshalSensorInfo(r *Record, v []byte) {
	r.Sensor = &SensorInfo{
		Unfiltered: UnmarshalUint32(v[8:12]),
		Filtered:   UnmarshalUint32(v[12:16]),
		Rssi:       UnmarshalUint16(v[16:18]),
	}
}

const (
	SENSOR_NOT_ACTIVE     SpecialGlucose = 1
	MINIMAL_DEVIATION     SpecialGlucose = 2
	NO_ANTENNA            SpecialGlucose = 3
	SENSOR_NOT_CALIBRATED SpecialGlucose = 5
	COUNTS_DEVIATION      SpecialGlucose = 6
	ABSOLUTE_DEVIATION    SpecialGlucose = 9
	POWER_DEVIATION       SpecialGlucose = 10
	BAD_RF                SpecialGlucose = 12

	specialLimit = BAD_RF
)

// SpecialGlucose represents a gucose value used to encode various exceptional conditions.
type SpecialGlucose uint16

//go:generate stringer -type SpecialGlucose

// IsSpecial checks whether a glucose value falls in the SpecialGlucose range.
func IsSpecial(glucose uint16) bool {
	return glucose <= uint16(specialLimit)
}

// Trend represents a directional arrow displayed by the Dexcom CGM receiver.
type Trend byte

//go:generate stringer -type Trend

const (
	UP_UP          Trend = 1
	UP             Trend = 2
	UP_45          Trend = 3
	FLAT           Trend = 4
	DOWN_45        Trend = 5
	DOWN           Trend = 6
	DOWN_DOWN      Trend = 7
	NOT_COMPUTABLE Trend = 8
	OUT_OF_RANGE   Trend = 9
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

func (t Trend) Symbol() string {
	return trendSymbol[t]
}

const (
	EGV_DISPLAY_ONLY     = 1 << 15
	EGV_VALUE_MASK       = 1<<10 - 1
	EGV_TREND_ARROW_MASK = 1<<4 - 1
)

func unmarshalEgvInfo(r *Record, v []byte) {
	g := UnmarshalUint16(v[8:10])
	r.Egv = &EgvInfo{
		Glucose:     g & EGV_VALUE_MASK,
		DisplayOnly: g&EGV_DISPLAY_ONLY != 0,
		Trend:       Trend(v[10] & EGV_TREND_ARROW_MASK),
	}
}

func unmarshalCalibrationInfo(r *Record, v []byte) {
	cal := &CalibrationInfo{
		Slope:     UnmarshalFloat64(v[8:16]),
		Intercept: UnmarshalFloat64(v[16:24]),
		Scale:     UnmarshalFloat64(v[24:32]),
		Decay:     UnmarshalFloat64(v[35:43]),
	}
	n := int(v[43])
	cal.Data = make([]CalibrationData, n)
	v = v[44:]
	offset := r.Timestamp.DisplayTime.Sub(r.Timestamp.SystemTime)
	for i := 0; i < n; i++ {
		cal.Data[i].Unmarshal(v)
		cal.Data[i].TimeEntered = cal.Data[i].TimeEntered.Add(offset)
		cal.Data[i].TimeApplied = cal.Data[i].TimeApplied.Add(offset)
		v = v[17:]
	}
	r.Calibration = cal
}

func (r *CalibrationData) Unmarshal(v []byte) {
	r.TimeEntered = UnmarshalTime(v[0:4])
	r.Glucose = UnmarshalInt32(v[4:8])
	r.Raw = UnmarshalInt32(v[8:12])
	r.TimeApplied = UnmarshalTime(v[12:16])
}

func unmarshalMeterInfo(r *Record, v []byte) {
	r.Meter = &MeterInfo{
		Glucose:   UnmarshalUint16(v[8:10]),
		MeterTime: UnmarshalTime(v[10:14]),
	}
}
