package dexcom

import (
	"log"
	"time"
)

func ReadEgvRecords(since time.Time) ([]EgvRecord, error) {
	context, err := ReadPageRange(EGV_DATA)
	if err != nil {
		return nil, err
	}
	results := []EgvRecord{}
	prevPage := -1
	proc := func(v []byte, context RecordContext) (bool, error) {
		if context.PageNumber != prevPage {
			log.Printf("scanning page %d", context.PageNumber)
			prevPage = context.PageNumber
		}
		r := EgvRecord{}
		err := r.Unmarshal(v)
		if err != nil {
			return false, err
		}
		t := r.Timestamp.DisplayTime
		if t.Before(since) {
			log.Printf("stopping at timestamp %s", t.Format(time.RFC3339))
			return false, nil
		}
		results = append(results, r)
		return true, nil
	}
	err = IterRecords(context, proc)
	return results, err
}

func ReadSensorRecords(since time.Time) ([]SensorRecord, error) {
	context, err := ReadPageRange(SENSOR_DATA)
	if err != nil {
		return nil, err
	}
	results := []SensorRecord{}
	prevPage := -1
	proc := func(v []byte, context RecordContext) (bool, error) {
		if context.PageNumber != prevPage {
			log.Printf("scanning page %d", context.PageNumber)
			prevPage = context.PageNumber
		}
		r := SensorRecord{}
		err := r.Unmarshal(v)
		if err != nil {
			return false, err
		}
		t := r.Timestamp.DisplayTime
		if t.Before(since) {
			log.Printf("stopping at timestamp %s", t.Format(time.RFC3339))
			return false, nil
		}
		results = append(results, r)
		return true, nil
	}
	err = IterRecords(context, proc)
	return results, err
}
