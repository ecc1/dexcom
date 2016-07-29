package dexcom

import (
	"fmt"

	"github.com/ecc1/nightscout"
)

func (r Record) NightscoutEntry() nightscout.Entry {
	t := r.Time()
	e := nightscout.Entry{
		Date:       t.UnixNano() / 1000000,
		DateString: t.Format(nightscout.DateStringLayout),
		Device:     "agape://ecc",
		Noise:      1,
	}
	if r.Calibration != nil {
		e.Type = "cal"
		e.Slope = r.Calibration.Slope
		e.Intercept = r.Calibration.Intercept
		e.Scale = r.Calibration.Scale
		return e
	}
	if r.Meter != nil {
		e.Type = "mbg"
		e.Mbg = r.Meter.Glucose
		return e
	}
	if r.Sensor != nil || r.Egv != nil {
		e.Type = "sgv"
		if r.Sensor != nil {
			e.Unfiltered = r.Sensor.Unfiltered
			e.Filtered = r.Sensor.Filtered
			e.Rssi = r.Sensor.Rssi
		}
		if r.Egv != nil {
			e.Sgv = r.Egv.Glucose
			e.Direction = nightscoutTrend(r.Egv.Trend)
		}
		return e
	}
	panic(fmt.Sprintf("NightscoutEntry(%+v}", r))
}

func nightscoutTrend(t Trend) string {
	switch t {
	case UP_UP:
		return "DoubleUp"
	case UP:
		return "SingeUp"
	case UP_45:
		return "FortyFiveUp"
	case FLAT:
		return "Flat"
	case DOWN_45:
		return "FortyFiveDown"
	case DOWN:
		return "SingleDown"
	case DOWN_DOWN:
		return "DoubleDown"
	default:
		return ""
	}
}
