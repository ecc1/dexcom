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
	proc := func(v []byte, context RecordContext) (bool, error) {
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
	proc := func(v []byte, context RecordContext) (bool, error) {
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

func ReadCalibrationRecords(since time.Time) ([]CalibrationRecord, error) {
	context, err := ReadPageRange(CAL_SET)
	if err != nil {
		return nil, err
	}
	results := []CalibrationRecord{}
	proc := func(v []byte, context RecordContext) (bool, error) {
		r := CalibrationRecord{}
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
