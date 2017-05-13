package dexcom

import (
	"fmt"
	"time"

	"github.com/ecc1/nightscout"
)

// NightscoutEntry converts a Record to a nightscout.Entry.
func (r Record) NightscoutEntry() nightscout.Entry {
	t := r.Time()
	e := nightscout.Entry{
		Date:       t.UnixNano() / 1000000,
		DateString: t.Format(nightscout.DateStringLayout),
		Device:     nightscout.Hostname(),
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
		e.MBG = int(r.Meter.Glucose)
		return e
	}
	if r.Sensor != nil || r.EGV != nil {
		e.Type = "sgv"
		if r.Sensor != nil {
			e.Unfiltered = int(r.Sensor.Unfiltered)
			e.Filtered = int(r.Sensor.Filtered)
			e.RSSI = int(r.Sensor.RSSI)
		}
		if r.EGV != nil {
			e.SGV = int(r.EGV.Glucose)
			e.Direction = nightscoutTrend(r.EGV.Trend)
			e.Noise = int(r.EGV.Noise)
		}
		return e
	}
	panic(fmt.Sprintf("NightscoutEntry(%+v}", r))
}

func nightscoutTrend(t Trend) string {
	switch t {
	case UpUp:
		return "DoubleUp"
	case Up:
		return "SingeUp"
	case Up45:
		return "FortyFiveUp"
	case Flat:
		return "Flat"
	case Down45:
		return "FortyFiveDown"
	case Down:
		return "SingleDown"
	case DownDown:
		return "DoubleDown"
	default:
		return ""
	}
}

const (
	edgeMargin = 1 * time.Minute
)

// MissingNightscoutEntries returns nightscout.Entry values
// for those records that fall within the given gaps.
func MissingNightscoutEntries(records []Record, gaps []nightscout.Gap) []nightscout.Entry {
	var missing []nightscout.Entry
	i := 0
	for _, g := range gaps {
		// Skip over records that lie outside the gap.
		for i < len(records) {
			t := records[i].Time()
			if t.Before(g.Finish) {
				break
			}
			i++
		}
		// Add records that fall within the gap
		// (by a margin of at least edgeMargin to avoid duplicates).
		for i < len(records) {
			r := records[i]
			t := r.Time()
			if t.Before(g.Start) {
				break
			}
			if t.Sub(g.Start) >= edgeMargin && g.Finish.Sub(t) >= edgeMargin {
				missing = append(missing, r.NightscoutEntry())
			}
			i++
		}
	}
	return missing
}
