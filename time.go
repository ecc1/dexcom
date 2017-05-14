package dexcom

import (
	"time"
)

const (
	// UserTimeLayout specifies a consistent, human-readable format for local time.
	UserTimeLayout = "2006-01-02 15:04:05"
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

func fromTime(t time.Time) int64 {
	// Construct the corresponding value in UTC.
	year, month, day := t.Date()
	hour, min, sec := t.Clock()
	u := time.Date(year, month, day, hour, min, sec, 0, time.UTC)
	return int64(u.Sub(dexcomEpoch) / time.Second)
}

func unmarshalTime(v []byte) time.Time {
	return toTime(int64(unmarshalUint32(v)))
}

// A Timestamp contains system and display time values.
type Timestamp struct {
	SystemTime  time.Time
	DisplayTime time.Time
}

func (r *Timestamp) unmarshal(v []byte) {
	r.SystemTime = unmarshalTime(v[0:4])
	r.DisplayTime = unmarshalTime(v[4:8])
}

func displayTime(sys uint32, offset int32) time.Time {
	return toTime(int64(sys) + int64(offset))
}

// ReadDisplayTime returns the Dexcom receiver's display time.
//	SystemTime = RTC + SystemTimeOffset
//	DisplayTime = SystemTime + DisplayTimeOffset
func (cgm *CGM) ReadDisplayTime() time.Time {
	v := cgm.Cmd(ReadDisplayTimeOffset)
	if cgm.Error() != nil {
		return time.Time{}
	}
	displayOffset := unmarshalInt32(v)
	v = cgm.Cmd(ReadSystemTime)
	if cgm.Error() != nil {
		return time.Time{}
	}
	sysTime := unmarshalUint32(v)
	return displayTime(sysTime, displayOffset)
}

// SetDisplayTime sets the Dexcom receiver's display time.
func (cgm *CGM) SetDisplayTime(t time.Time) {
	v := cgm.Cmd(ReadSystemTime)
	if cgm.Error() != nil {
		return
	}
	sysTime := unmarshalUint32(v)
	offset := int32(fromTime(t) - int64(sysTime))
	cgm.Cmd(WriteDisplayTimeOffset, marshalInt32(offset)...)
}
