package dexcom

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/ecc1/nightscout"
)

func jsonString(v interface{}) string {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(data)
}

func (r Record) String() string { return jsonString(r) }

type (
	Entry   nightscout.Entry
	Entries []Entry
)

func (e Entry) String() string { return jsonString(e) }

func jsonTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

func ts(s string) Timestamp {
	return Timestamp{DisplayTime: jsonTime(s)}
}

func nsDate(s string) int64 {
	return nightscout.Date(jsonTime(s))
}

var (
	r1 = Record{
		Timestamp: ts("2017-09-17T01:13:51-04:00"),
		Info: CalibrationInfo{
			Slope:     939.6817717490421,
			Intercept: 35926.604186515906,
			Scale:     1,
		},
	}
	r2 = Record{
		Timestamp: ts("2017-09-17T01:13:49-04:00"),
		Info: MeterInfo{
			Glucose: 128,
		},
	}
	r3 = Record{
		Timestamp: ts("2017-09-17T11:13:17-04:00"),
		Info: EGVInfo{
			Glucose: 84,
			Trend:   Flat,
			Noise:   1,
		},
	}
	r4 = Record{
		Timestamp: ts("2017-09-17T11:13:16-04:00"),
		Info: SensorInfo{
			Unfiltered: 119088,
			Filtered:   110288,
			RSSI:       -62,
		},
	}

	dev = nightscout.Device()

	e1 = Entry{
		Type:       "cal",
		Date:       nsDate("2017-09-17T01:13:51-04:00"),
		DateString: "2017-09-17T01:13:51-04:00",
		Device:     dev,
		Slope:      939.6817717490421,
		Intercept:  35926.604186515906,
		Scale:      1,
	}
	e2 = Entry{
		Type:       "mbg",
		Date:       nsDate("2017-09-17T01:13:49-04:00"),
		DateString: "2017-09-17T01:13:49-04:00",
		Device:     dev,
		MBG:        128,
	}
	e3 = Entry{
		Type:       "sgv",
		Date:       nsDate("2017-09-17T11:13:17-04:00"),
		DateString: "2017-09-17T11:13:17-04:00",
		Device:     dev,
		SGV:        84,
		Direction:  "Flat",
		Noise:      1,
	}
	e4 = Entry{
		Type:       "sgv",
		Date:       nsDate("2017-09-17T11:13:16-04:00"),
		DateString: "2017-09-17T11:13:16-04:00",
		Device:     dev,
		Unfiltered: 119088,
		Filtered:   110288,
		RSSI:       -62,
	}
	e5 = Entry{
		Type:       "sgv",
		Date:       nsDate("2017-09-17T11:13:16-04:00"),
		DateString: "2017-09-17T11:13:16-04:00",
		Device:     dev,
		Unfiltered: 119088,
		Filtered:   110288,
		RSSI:       -62,
		SGV:        84,
		Direction:  "Flat",
		Noise:      1,
	}
)

func TestNightscoutEntry(t *testing.T) {
	cases := []struct {
		r Record
		e Entry
	}{
		{r1, e1},
		{r2, e2},
		{r3, e3},
		{r4, e4},
	}
	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			e := Entry(c.r.nightscoutEntry())
			if e != c.e {
				t.Errorf("nightscoutEntry(%v) == %v, want %v", c.r, e, c.e)
			}
		})
	}
}

func TestNightscoutEntries(t *testing.T) {
	cases := []struct {
		r Records
		e Entries
	}{
		{
			Records{r3, r4},
			Entries{e5},
		},
	}
	for _, c := range cases {
		t.Run("", func(t *testing.T) {
			e := convertEntries(NightscoutEntries(c.r))
			if !equalEntries(e, c.e) {
				t.Errorf("NightscoutEntries(%v) == %v, want %v", c.r, e, c.e)
			}
		})
	}
}

func convertEntries(v nightscout.Entries) Entries {
	entries := make(Entries, len(v))
	for i := range v {
		entries[i] = Entry(v[i])
	}
	return entries
}

func equalEntries(x, y Entries) bool {
	if len(x) != len(y) {
		return false
	}
	for i := range x {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}
