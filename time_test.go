package dexcom

import (
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	cases := []struct {
		b []byte
		t time.Time
	}{
		{[]byte{0x00, 0x00, 0x00, 0x00}, parseTime("2009-01-01T00:00:00")},
		{[]byte{0x01, 0x02, 0x03, 0x04}, parseTime("2011-02-19T00:06:25")},
		{[]byte{0x74, 0x90, 0x33, 0x0E}, parseTime("2016-07-20T15:25:40")},
	}
	for _, c := range cases {
		tv := UnmarshalTime(c.b)
		if !tv.Equal(c.t) {
			t.Errorf("UnmarshalTime(% X) == %v, want %v", c.b, tv, c.t)
		}
	}
}

func parseTime(s string) time.Time {
	const layout = "2006-01-02T15:04:05"
	t, err := time.ParseInLocation(layout, s, time.Local)
	if err != nil {
		panic(err)
	}
	return t
}
