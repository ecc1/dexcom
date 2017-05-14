package dexcom

import (
	"log"
	"time"
)

// ReadHistory returns records since the specified time.
func (cgm *CGM) ReadHistory(pageType PageType, since time.Time) []Record {
	first, last := cgm.ReadPageRange(pageType)
	if cgm.Error() != nil {
		return nil
	}
	var results []Record
	proc := func(r Record) (bool, error) {
		t := r.Time()
		if t.Before(since) {
			log.Printf("stopping %v scan at %s", pageType, t.Format(UserTimeLayout))
			return true, nil
		}
		results = append(results, r)
		return false, nil
	}
	cgm.IterRecords(pageType, first, last, proc)
	return results
}

// ReadCount returns a specified number of most recent records.
func (cgm *CGM) ReadCount(pageType PageType, count int) []Record {
	first, last := cgm.ReadPageRange(pageType)
	if cgm.Error() != nil {
		return nil
	}
	var results []Record
	proc := func(r Record) (bool, error) {
		results = append(results, r)
		return len(results) == count, nil
	}
	cgm.IterRecords(pageType, first, last, proc)
	return results
}

// MergeHistory merges slices of records that are already
// in reverse chronological order into a single ordered slice.
func MergeHistory(slices ...[]Record) []Record {
	n := len(slices)
	if n == 0 {
		return nil
	}
	if n == 1 {
		return slices[0]
	}
	length := make([]int, n)
	total := 0
	for i, v := range slices {
		length[i] = len(v)
		total += len(v)
	}
	results := make([]Record, total)
	index := make([]int, n)
	for next := range results {
		// Find slice with latest current value.
		which := -1
		max := time.Time{}
		for i, v := range slices {
			if index[i] < len(v) {
				t := v[index[i]].Time()
				if t.After(max) {
					which = i
					max = t
				}
			}
		}
		results[next] = slices[which][index[which]]
		index[which]++
	}
	return results
}

const (
	// Time window within which EGV and sensor readings will be merged.
	glucoseReadingWindow = 10 * time.Second
)

// GlucoseReadings returns sensor and EGV records since the specified time.
func (cgm *CGM) GlucoseReadings(since time.Time) []Record {
	sensor := cgm.ReadHistory(SensorData, since)
	if cgm.Error() != nil {
		return nil
	}
	egv := cgm.ReadHistory(EGVData, since)
	if cgm.Error() != nil {
		return nil
	}
	readings := make([]Record, 0, len(sensor))
	i, j := 0, 0
	for {
		var r Record
		if i < len(sensor) && j < len(egv) {
			r = chooseRecord(sensor, egv, &i, &j)
		} else if i < len(sensor) {
			r = sensor[i]
			i++
		} else if j < len(egv) {
			r = egv[j]
			j++
		} else {
			break
		}
		readings = append(readings, r)
	}
	return readings
}

func chooseRecord(sensor, egv []Record, ip, jp *int) Record {
	i := *ip
	j := *jp
	sensorTime := sensor[i].Time()
	egvTime := egv[j].Time()
	delta := egvTime.Sub(sensorTime)
	var r Record
	if 0 <= delta && delta < glucoseReadingWindow {
		// Merge using sensor[i]'s slightly earlier time.
		r = sensor[i]
		r.EGV = egv[j].EGV
		i++
		j++
	} else if 0 <= -delta && -delta < glucoseReadingWindow {
		// Merge using egv[j]'s slightly earlier time.
		r = egv[j]
		r.Sensor = sensor[i].Sensor
		i++
		j++
	} else if sensorTime.After(egvTime) {
		r = sensor[i]
		i++
	} else {
		r = egv[j]
		j++
	}
	*ip = i
	*jp = j
	return r
}
