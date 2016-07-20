package dexcom

import (
	"time"
)

var (
	dexcomEpoch = time.Date(2009, 1, 1, 0, 0, 0, 0, time.UTC)
)

func toTime(t int64) time.Time {
	u := dexcomEpoch.Add(time.Duration(t) * time.Second)
	// Construct the corresponding value in the local timezone.
	year, month, day := u.Date()
	hour, min, sec := u.Clock()
	return time.Date(year, month, day, hour, min, sec, 0, time.Local)
}

// UnmarshalTime unmarshals a 4-byte array into a time value.
func UnmarshalTime(v []byte) time.Time {
	return toTime(int64(UnmarshalUint32(v)))
}

// A Timestamp contains system and display time values.
type Timestamp struct {
	SystemTime  time.Time
	DisplayTime time.Time
}

func (r *Timestamp) Type() RecordType {
	return timestampType
}

func (r *Timestamp) Unmarshal(v []byte) error {
	r.SystemTime = UnmarshalTime(v[0:4])
	r.DisplayTime = UnmarshalTime(v[4:8])
	return nil
}

func displayTime(sys uint32, offset int32) time.Time {
	return toTime(int64(sys) + int64(offset))
}

// SYSTEM_TIME = RTC + SYSTEM_TIME_OFFSET
// DISPLAY_TIME = SYSTEM_TIME + DISPLAY_TIME_OFFSET

// ReadDisplayTime gets the current display time value from the Dexcom CGM receiver.
func ReadDisplayTime() (time.Time, error) {
	v, err := Cmd(READ_DISPLAY_TIME_OFFSET)
	if err != nil {
		return time.Time{}, err
	}
	displayOffset := UnmarshalInt32(v)
	v, err = Cmd(READ_SYSTEM_TIME)
	if err != nil {
		return time.Time{}, err
	}
	sysTime := UnmarshalUint32(v)
	return displayTime(sysTime, displayOffset), nil
}
