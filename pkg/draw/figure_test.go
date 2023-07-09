package draw

import (
	"gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/pkg/storedata"
	"gotest.tools/assert"
	"sort"
	"testing"
)

func TestGetLinesForInstances(t *testing.T) {

	path := "../../data/"
	fileName := "benchmark_fmais_2023-03-26_01-59-08.json"

	tc, err := storedata.ImportData(path, fileName)
	assert.NilError(t, err)

	lines := getLinesForInstances(tc, []int{2, 3, 4}, []int{26}, []int{26, 56, 96})

	keys := make([]string, 0, len(lines))
	for key := range lines {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	//check that it is in order now
	for _, key := range keys {
		line := lines[key]
		t.Logf("Line for key %v is %v", key, line)
	}
}
