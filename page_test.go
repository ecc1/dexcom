package dexcom

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
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
	page, err := unmarshalPage(data)
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
	decoded, err := unmarshalRecords(c.pageType, page.Records)
	if err != nil {
		t.Errorf("unmarshalRecords(%v, % X) returned %v", c.pageType, page.Records, err)
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
	eq, msg := compareJSON(decoded, jsonFile)
	if !eq {
		t.Errorf("JSON is different:\n%s\n", msg)
	}
}

func readBytes(r io.Reader) ([]byte, error) {
	var data []byte
	for {
		var b byte
		n, err := fmt.Fscanf(r, "%02x", &b)
		if n == 0 {
			break
		}
		if err != nil {
			return data, err
		}
		data = append(data, b)
	}
	return data, nil
}

func compareJSON(data interface{}, jsonFile string) (bool, string) {
	// Write data in JSON format to temporary file.
	tmpfile, err := ioutil.TempFile("", "json")
	if err != nil {
		return false, err.Error()
	}
	defer func() { _ = os.Remove(tmpfile.Name()) }()
	e := json.NewEncoder(tmpfile)
	e.SetIndent("", "  ")
	err = e.Encode(data)
	_ = tmpfile.Close()
	if err != nil {
		return false, err.Error()
	}
	// Write JSON in canonical form for comparison.
	canon1 := canonicalJSON(jsonFile)
	canon2 := canonicalJSON(tmpfile.Name())
	// Find differences.
	cmd := exec.Command("diff", "-u", "--label", jsonFile, "--label", "decoded", canon1, canon2)
	diffs, err := cmd.Output()
	_ = os.Remove(canon1)
	_ = os.Remove(canon2)
	return err == nil, string(diffs)
}

// canonicalJSON reads the given file and creates a temporary file
// containing equivalent JSON in canonical form
// (using the "jq" command, which must be on the user's PATH).
// It returns the temporary file name; it is the caller's responsibility
// to remove it when done.
func canonicalJSON(file string) string {
	canon, err := exec.Command("jq", "-S", ".", file).Output()
	if err != nil {
		panic(err)
	}
	tmpfile, err := ioutil.TempFile("", "json")
	if err != nil {
		panic(err)
	}
	_, _ = tmpfile.Write(canon)
	_ = tmpfile.Close()
	return tmpfile.Name()
}
