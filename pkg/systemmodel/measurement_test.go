package systemmodel

import (
	"gotest.tools/assert"
	"testing"
)

func TestCreateSystemModelDepth4(t *testing.T) {
	sm := CreateSystemModelDepth4()
	assert.Assert(t, sm != nil)

	// updating application instances' reliability values
	err := sm.UpdateApplicationReliability("App#4", map[int64]float64{
		1: 0.7,
	})
	assert.NilError(t, err)

	// getting instance back to check if the reliability was actually set
	app4, err := sm.GetInstance("App#4-4-1")
	assert.NilError(t, err)

	// retrieving reliability of an instance
	app4Rel, err := app4.GetReliability()
	assert.NilError(t, err)
	assert.Equal(t, app4Rel, 0.7)
}

func TestCreateSystemModelDepth3(t *testing.T) {
	sm := CreateSystemModelDepth3()
	assert.Assert(t, sm != nil)

	// updating application instances' reliability values
	err := sm.UpdateApplicationReliability("App#1", map[int64]float64{
		1: 0.69,
		2: 0.55,
		3: 0.9,
	})
	assert.NilError(t, err)

	// getting instance back to check if the reliability was actually set
	app11, err := sm.GetInstance("App#2-1-1")
	assert.NilError(t, err)

	// retrieving reliability of an instance
	app11Rel, err := app11.GetReliability()
	assert.NilError(t, err)
	assert.Equal(t, app11Rel, 0.69)

	// getting instance back to check if the reliability was actually set
	app12, err := sm.GetInstance("App#2-1-2")
	assert.NilError(t, err)

	// retrieving reliability of an instance
	app12Rel, err := app12.GetReliability()
	assert.NilError(t, err)
	assert.Equal(t, app12Rel, 0.55)

	// getting instance back to check if the reliability was actually set
	app13, err := sm.GetInstance("App#2-1-3")
	assert.NilError(t, err)

	// retrieving reliability of an instance
	app13Rel, err := app13.GetReliability()
	assert.NilError(t, err)
	assert.Equal(t, app13Rel, 0.9)
}

func TestCreateSystemModelDepth2(t *testing.T) {
	sm := CreateSystemModelDepth2()
	assert.Assert(t, sm != nil)

	// updating application instances' reliability values
	err := sm.UpdateApplicationReliability("VI#2-1", map[int64]float64{
		1: 0.0123,
	})
	assert.NilError(t, err)

	// getting instance back to check if the reliability was actually set
	vi1, err := sm.GetInstance("VI#2-1")
	assert.NilError(t, err)

	// retrieving reliability of an instance
	vi1Rel, err := vi1.GetReliability()
	assert.NilError(t, err)
	assert.Equal(t, vi1Rel, 0.0123)
}
