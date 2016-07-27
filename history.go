package dexcom

import (
	"log"
	"time"
)

func (cgm *Cgm) ReadHistory(pageType PageType, since time.Time) []Record {
	first, last := cgm.ReadPageRange(pageType)
	if cgm.Error() != nil {
		return nil
	}
	results := []Record{}
	proc := func(r Record) (bool, error) {
		t := r.Timestamp.DisplayTime
		if t.Before(since) {
			log.Printf("stopping at timestamp %s", t.Format(time.RFC3339))
			return true, nil
		}
		results = append(results, r)
		return false, nil
	}
	cgm.IterRecords(pageType, first, last, proc)
	return results
}

func (cgm *Cgm) ReadCount(pageType PageType, count int) []Record {
	first, last := cgm.ReadPageRange(pageType)
	if cgm.Error() != nil {
		return nil
	}
	results := []Record{}
	proc := func(r Record) (bool, error) {
		results = append(results, r)
		return len(results) == count, nil
	}
	cgm.IterRecords(pageType, first, last, proc)
	return results
}

// Merge slices of records that are already in reverse chronological order
// into a single ordered slice.
func MergeHistory(slices ...[]Record) []Record {
	n := len(slices)
	if n == 0 {
		return nil
	}
	length := make([]int, n)
	total := 0
	for i, v := range slices {
		length[i] = len(v)
		total += len(v)
	}
	results := make([]Record, total)
	index := make([]int, n)
	for next, _ := range results {
		// Find slice with latest current value.
		which := -1
		max := time.Time{}
		for i, v := range slices {
			if index[i] < len(v) {
				t := v[index[i]].Timestamp.DisplayTime
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
	glucoseReadingWindow = 2 * time.Second
)

func (cgm *Cgm) GlucoseReadings(since time.Time) []Record {
	sensor := cgm.ReadHistory(SENSOR_DATA, since)
	if cgm.Error() != nil {
		return nil
	}
	numSensor := len(sensor)
	egv := cgm.ReadHistory(EGV_DATA, since)
	if cgm.Error() != nil {
		return nil
	}
	numEgv := len(egv)
	readings := []Record{}
	i, j := 0, 0
	for {
		r := Record{}
		if i < numSensor && j < numEgv {
			sensorTime := sensor[i].Timestamp.DisplayTime
			egvTime := egv[j].Timestamp.DisplayTime
			delta := egvTime.Sub(sensorTime)
			if 0 <= delta && delta < glucoseReadingWindow {
				// Merge using sensor[i]'s slightly earlier time.
				r = sensor[i]
				r.Egv = egv[j].Egv
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
		} else if i < numSensor {
			r = sensor[i]
			i++
		} else if j < numEgv {
			r = egv[j]
			j++
		} else {
			break
		}
		readings = append(readings, r)
	}
	return readings
}
