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
	switch info := r.Info.(type) {
	case CalibrationInfo:
		e.Type = "cal"
		e.Slope = info.Slope
		e.Intercept = info.Intercept
		e.Scale = info.Scale
		return e
	case MeterInfo:
		e.Type = "mbg"
		e.MBG = int(info.Glucose)
		return e
	case BGInfo:
		e.Type = "sgv"
		e.Unfiltered = int(info.Sensor.Unfiltered)
		e.Filtered = int(info.Sensor.Filtered)
		e.RSSI = int(info.Sensor.RSSI)
		e.SGV = int(info.EGV.Glucose)
		e.Direction = nightscoutTrend(info.EGV.Trend)
		e.Noise = int(info.EGV.Noise)
		return e
	}
	panic(fmt.Sprintf("NightscoutEntry %+v", r))
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
func MissingNightscoutEntries(records Records, gaps []nightscout.Gap) []nightscout.Entry {
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
