package dexcom

import (
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	cases := []struct {
		n int64
		t time.Time
	}{
		{0x00000000, parseTime("2009-01-01 00:00:00")},
		{0x01234567, parseTime("2009-08-09 22:25:43")},
		{0x04030201, parseTime("2011-02-19 00:06:25")},
		{0x0E339074, parseTime("2016-07-20 15:25:40")},
		{0x0EDCBA98, parseTime("2016-11-25 22:58:32")},
	}
	for _, c := range cases {
		tv := toTime(c.n)
		if !tv.Equal(c.t) {
			t.Errorf("toTime(%X) == %v, want %v", c.n, tv, c.t)
		}
		n := fromTime(c.t)
		if n != c.n {
			t.Errorf("fromTime(%v) == %X, want %X", c.t, n, c.n)
		}
	}
}

func parseTime(s string) time.Time {
	t, err := time.ParseInLocation(UserTimeLayout, s, time.Local)
	if err != nil {
		panic(err)
	}
	return t
}
