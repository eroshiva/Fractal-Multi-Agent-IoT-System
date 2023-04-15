package storedata

import (
	"gotest.tools/assert"
	"testing"
)

func Test_exportDataToJSON(t *testing.T) {
	// defining a path and a filename to import data
	path := "../../data/"
	filename := "unittest"

	// creating dummy data
	data := make(map[int]map[int]map[int]float64, 0)
	data[100] = make(map[int]map[int]float64, 0)
	data[100][20] = make(map[int]float64, 0)
	data[100][20][31] = 3.1457

	// exporting data
	err := exportDataToJSON(path, filename, data, "", " ")
	assert.NilError(t, err)
}

func Test_importDataFromJSON(t *testing.T) {
	// defining a path and a filename to import data
	path := "../../data/"
	filename := "unittest"

	// importing data
	data, err := importDataFromJSON(path, filename)
	assert.NilError(t, err)
	assert.Assert(t, len(data) > 0)
	t.Logf("Imported data are:\n%v\n", data)

	value, ok := data[100][20][31]
	assert.Equal(t, ok, true)
	assert.Equal(t, value, 3.1457)
}

func Test_exportDataToCSV(t *testing.T) {
	// defining a path and a filename to import data
	path := "../../data/"
	filename := "unittest"

	// creating dummy data
	data := make(map[int]map[int]map[int]float64, 0)
	data[100] = make(map[int]map[int]float64, 0)
	data[100][20] = make(map[int]float64, 0)
	data[100][20][31] = 3.1457

	// exporting data
	err := exportDataToCSV(path, filename, data, "Fractal MAIS Depth [-]", "Application Number in Fractal MAIS [-]",
		"Maximum Number of Instances Deployed by Application [-]", "Time [us]")
	assert.NilError(t, err)
}

func Test_exportDataFromCSV(t *testing.T) {
	// defining a path and a filename to import data
	path := "../../data/"
	filename := "unittest"

	// importing data
	data, err := importDataFromCSV(path, filename)
	assert.NilError(t, err)
	assert.Assert(t, len(data) > 0)
	t.Logf("Imported data are:\n%v\n", data)

	value, ok := data[100][20][31]
	assert.Equal(t, ok, true)
	assert.Equal(t, value, 3.1457)
}
