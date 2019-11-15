package dexcom

import (
	"fmt"
	"log"
	"time"

	"github.com/ecc1/nightscout"
)

// NightscoutEntries converts records (in reverse-chronological order)
// into a Nightscout entries.  Neighboring Sensor and EGV records are merged.
func NightscoutEntries(records Records) nightscout.Entries {
	entries := make(nightscout.Entries, len(records))
	for i, r := range records {
		entries[i] = r.nightscoutEntry()
	}
	return mergeGlucoseEntries(entries)
}

func (r Record) nightscoutEntry() nightscout.Entry {
	t := r.Time()
	e := nightscout.Entry{
		Date:       nightscout.Date(t),
		DateString: t.Format(nightscout.DateStringLayout),
		Device:     nightscout.Device(),
	}
	if r.Sensor != nil {
		info := r.Sensor
		e.Type = nightscout.SGVType
		e.Unfiltered = int(info.Unfiltered)
		e.Filtered = int(info.Filtered)
		e.RSSI = int(info.RSSI)
		return e
	}
	if r.EGV != nil {
		info := r.EGV
		e.Type = nightscout.SGVType
		e.SGV = int(info.Glucose)
		e.Direction = nightscoutTrend(info.Trend)
		e.Noise = int(info.Noise)
		return e
	}
	if r.Calibration != nil {
		info := r.Calibration
		e.Type = nightscout.CalType
		e.Slope = info.Slope
		e.Intercept = info.Intercept
		e.Scale = info.Scale
		return e
	}
	if r.Insertion != nil {
		info := r.Insertion
		switch info.Event {
		case Stopped:
			e.Type = "Sensor Change"
		case Started:
			e.Type = "Sensor Start"
		default:
			e.Type = fmt.Sprintf("%d", info.Event)
		}
		return e
	}
	if r.Meter != nil {
		info := r.Meter
		e.Type = nightscout.MBGType
		e.MBG = int(info.Glucose)
		return e
	}
	panic(fmt.Sprintf("nightscoutEntry %+v", r))
}

func nightscoutTrend(t Trend) string {
	switch t {
	case UpUp:
		return "DoubleUp"
	case Up:
		return "SingleUp"
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
	// Time window within which sensor and EGV readings will be merged.
	glucoseReadingWindow = 10 * time.Second
)

func mergeGlucoseEntries(entries nightscout.Entries) nightscout.Entries {
	merged := make(nightscout.Entries, 0, len(entries))
	i := 0
	for i < len(entries) {
		e := entries[i]
		if e.Type == nightscout.SGVType && i+1 < len(entries) {
			f := entries[i+1]
			if f.Type == nightscout.SGVType {
				delta := e.Time().Sub(f.Time())
				if 0 <= delta && delta < glucoseReadingWindow {
					e = combineEntries(e, f)
					i++
				}
			}
		}
		merged = append(merged, e)
		i++
	}
	return merged
}

func combineEntries(a, b nightscout.Entry) nightscout.Entry {
	if a.Type != nightscout.SGVType || b.Type != nightscout.SGVType {
		log.Panicf("combining %s and %s", a.Type, b.Type)
	}
	if b.Time().Before(a.Time()) {
		// Use b's earlier time.
		a.Date = b.Date
		a.DateString = b.DateString
	}
	// Update a with non-zero sgv values from b.
	if b.SGV != 0 {
		a.SGV = b.SGV
	}
	if b.Direction != "" {
		a.Direction = b.Direction
	}
	if b.Filtered != 0 {
		a.Filtered = b.Filtered
	}
	if b.Unfiltered != 0 {
		a.Unfiltered = b.Unfiltered
	}
	if b.RSSI != 0 {
		a.RSSI = b.RSSI
	}
	if b.Noise != 0 {
		a.Noise = b.Noise
	}
	return a
}
