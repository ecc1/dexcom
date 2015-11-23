package dexcom

import (
	"time"
)

var (
	baseTime = time.Date(2009, 1, 1, 0, 0, 0, 0, time.Local)
)

func toTime(t uint32) time.Time {
	return baseTime.Add(time.Duration(int64(t)) * time.Second)
}

func UnmarshalTime(v []byte) time.Time {
	return toTime(UnmarshalUint32(v))
}

type Timestamp struct {
	SystemTime  time.Time
	DisplayTime time.Time
}

func UnmarshalTimestamp(v []byte) Timestamp {
	return Timestamp{
		SystemTime:  UnmarshalTime(v[0:4]),
		DisplayTime: UnmarshalTime(v[4:8]),
	}
}

func displayTime(sys uint32, offset int32) time.Time {
	d := int64(sys) + int64(offset)
	return baseTime.Add(time.Duration(d) * time.Second)
}

// SYSTEM_TIME = RTC + SYSTEM_TIME_OFFSET
// DISPLAY_TIME = SYSTEM_TIME + DISPLAY_TIME_OFFSET

func (dev Device) ReadDisplayTime() (time.Time, error) {
	v, err := dev.Cmd(READ_DISPLAY_TIME_OFFSET)
	if err != nil {
		return time.Time{}, err
	}
	displayOffset := UnmarshalInt32(v)
	v, err = dev.Cmd(READ_SYSTEM_TIME)
	if err != nil {
		return time.Time{}, err
	}
	sysTime := UnmarshalUint32(v)
	return displayTime(sysTime, displayOffset), nil
}
