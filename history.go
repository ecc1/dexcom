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

type GlucoseReading struct {
	Egv    EgvRecord
	Sensor SensorRecord
}

const (
	// Time window within which EGV and sensor readings will be merged.
	glucoseReadingWindow = 2 * time.Second
)

func withinWindow(t1, t2 time.Time, window time.Duration) bool {
	d := t1.Sub(t2)
	return (0 <= d && d < window) || (0 <= -d && -d < window)
}

func GlucoseReadings(since time.Time) ([]GlucoseReading, error) {
	egv, err := ReadEgvRecords(since)
	if err != nil {
		return nil, err
	}
	numEgv := len(egv)
	sensor, err := ReadSensorRecords(since)
	if err != nil {
		return nil, err
	}
	numSensor := len(sensor)
	readings := []GlucoseReading{}
	i, j := 0, 0
	for {
		r := GlucoseReading{}
		if i < numEgv && j < numSensor {
			egvTime := egv[i].Timestamp.DisplayTime
			sensorTime := sensor[j].Timestamp.DisplayTime
			if withinWindow(egvTime, sensorTime, glucoseReadingWindow) {
				r = GlucoseReading{Egv: egv[i], Sensor: sensor[j]}
				i++
				j++
			} else if egvTime.After(sensorTime) {
				r = GlucoseReading{Egv: egv[i]}
				i++
			} else {
				r = GlucoseReading{Sensor: sensor[j]}
				j++
			}
		} else if i < numEgv {
			r = GlucoseReading{Egv: egv[i]}
			i++
		} else if j < numSensor {
			r = GlucoseReading{Sensor: sensor[j]}
			j++
		} else {
			break
		}
		readings = append(readings, r)
	}
	return readings, nil
}
