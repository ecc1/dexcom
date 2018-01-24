package dexcom

import (
	"log"
	"time"
)

// ReadHistory returns records since the specified time.
func (cgm *CGM) ReadHistory(pageType PageType, since time.Time) Records {
	first, last := cgm.ReadPageRange(pageType)
	if cgm.Error() != nil {
		return nil
	}
	var results Records
	proc := func(r Record) error {
		t := r.Time()
		if !t.After(since) {
			log.Printf("stopping %v scan at %s", pageType, t.Format(UserTimeLayout))
			return IterationDone
		}
		results = append(results, r)
		return nil
	}
	cgm.IterRecords(pageType, first, last, proc)
	return results
}

// ReadCount returns a specified number of most recent records.
func (cgm *CGM) ReadCount(pageType PageType, count int) Records {
	first, last := cgm.ReadPageRange(pageType)
	if cgm.Error() != nil {
		return nil
	}
	results := make(Records, 0, count)
	proc := func(r Record) error {
		results = append(results, r)
		if len(results) == count {
			return IterationDone
		}
		return nil
	}
	cgm.IterRecords(pageType, first, last, proc)
	return results
}

// MergeHistory merges slices of records that are already
// in reverse chronological order into a single ordered slice.
func MergeHistory(slices ...Records) Records {
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
	results := make(Records, total)
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
