package dexcom

import (
	"log"
	"time"
)

func (cgm *Cgm) ReadEgvRecords(since time.Time) []EgvRecord {
	context := cgm.ReadPageRange(EGV_DATA)
	if cgm.Error() != nil {
		return nil
	}
	results := []EgvRecord{}
	proc := func(v []byte, context RecordContext) (bool, error) {
		r := EgvRecord{}
		err := r.Unmarshal(v)
		if err != nil {
			return true, err
		}
		t := r.Timestamp.DisplayTime
		if t.Before(since) {
			log.Printf("stopping at timestamp %s", t.Format(time.RFC3339))
			return true, nil
		}
		results = append(results, r)
		return false, nil
	}
	cgm.IterRecords(context, proc)
	return results
}

func (cgm *Cgm) ReadSensorRecords(since time.Time) []SensorRecord {
	context := cgm.ReadPageRange(SENSOR_DATA)
	if cgm.Error() != nil {
		return nil
	}
	results := []SensorRecord{}
	proc := func(v []byte, context RecordContext) (bool, error) {
		r := SensorRecord{}
		err := r.Unmarshal(v)
		if err != nil {
			return true, err
		}
		t := r.Timestamp.DisplayTime
		if t.Before(since) {
			log.Printf("stopping at timestamp %s", t.Format(time.RFC3339))
			return true, nil
		}
		results = append(results, r)
		return false, nil
	}
	cgm.IterRecords(context, proc)
	return results
}

func (cgm *Cgm) ReadCalibrationRecords(since time.Time) []CalibrationRecord {
	context := cgm.ReadPageRange(CAL_SET)
	if cgm.Error() != nil {
		return nil
	}
	results := []CalibrationRecord{}
	proc := func(v []byte, context RecordContext) (bool, error) {
		r := CalibrationRecord{}
		err := r.Unmarshal(v)
		if err != nil {
			return true, err
		}
		t := r.Timestamp.DisplayTime
		if t.Before(since) {
			log.Printf("stopping at timestamp %s", t.Format(time.RFC3339))
			return true, nil
		}
		results = append(results, r)
		return false, nil
	}
	cgm.IterRecords(context, proc)
	return results
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

func (cgm *Cgm) GlucoseReadings(since time.Time) []GlucoseReading {
	egv := cgm.ReadEgvRecords(since)
	if cgm.Error() != nil {
		return nil
	}
	numEgv := len(egv)
	sensor := cgm.ReadSensorRecords(since)
	if cgm.Error() != nil {
		return nil
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
	return readings
}
