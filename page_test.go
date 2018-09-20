package dexcom

import (
	"fmt"
	"os"
	"testing"
)

const testDataDir = "testdata"

type testCase struct {
	pageType    PageType
	pageNumber  int
	alternative int
}

func TestPage(t *testing.T) {
	cases := []testCase{
		{ManufacturingData, 0, 0},
		{SensorData, 469, 0},
		{EGVData, 312, 0},
		{CalibrationData, 1432, 0},
		// The one record in this page has CRC = 63 FF, followed by FF bytes for padding.
		{CalibrationData, 1432, 1},
		// These pages contain 148-byte rev 2 calibration records.
		{CalibrationData, 252, 0},
		{CalibrationData, 845, 0},
	}
	for _, c := range cases {
		t.Run(c.pageType.String(), func(t *testing.T) {
			pageTest(t, c)
		})
	}
}

func pageTest(t *testing.T, c testCase) {
	testFile := testFileName(c)
	f, err := os.Open(testFile + ".data")
	if err != nil {
		t.Error(err)
		return
	}
	data, err := readBytes(f)
	_ = f.Close()
	if err != nil {
		t.Error(err)
		return
	}
	page, err := UnmarshalPage(data)
	if err != nil {
		t.Error(err)
		return
	}
	if page.Type != c.pageType {
		panic("page type mismatch")
	}
	if page.Number != c.pageNumber {
		panic("page number mismatch")
	}
	decoded, err := UnmarshalRecords(c.pageType, page.Records)
	if err != nil {
		t.Errorf("UnmarshalRecords(%v, % X) returned %v", c.pageType, page.Records, err)
		return
	}
	checkRecords(t, decoded, testFile+".json")
}

func testFileName(c testCase) string {
	s := fmt.Sprintf("%s/%d.%d", testDataDir, c.pageType, c.pageNumber)
	if c.alternative != 0 {
		s += fmt.Sprintf("-%d", c.alternative)
	}
	return s
}

func checkRecords(t *testing.T, decoded Records, jsonFile string) {
	eq, msg := compareDataToJSON(decoded, jsonFile)
	if !eq {
		t.Errorf("JSON is different:\n%s\n", msg)
	}
}
