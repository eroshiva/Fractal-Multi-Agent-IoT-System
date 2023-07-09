package measurement

import (
	"gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/pkg/meertcore"
	"gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/pkg/systemmodel"
	"gotest.tools/assert"
	"math"
	"testing"
)

func TestUpdateReliabilities(t *testing.T) {
	// initializing input data
	app1, app2, vi := initializeInputDataDepth2()

	// initialising system model
	sm2 := systemmodel.CreateSystemModelDepth2()

	meErtCore := meertcore.MeErtCore{
		SystemModel: sm2,
		Reliability: 0.0,
	}

	// Application #1, instance 1
	inst, err := meErtCore.SystemModel.GetInstance("App#2-1-1")
	assert.NilError(t, err)
	rel, err := inst.GetReliability()
	assert.NilError(t, err)
	assert.Equal(t, rel, 0.77)

	// setting initial reliabilities (in step #1)
	err = UpdateReliabilities(meErtCore.SystemModel, 1, app1, app2, vi)
	assert.NilError(t, err)

	// verifying that they were set successfully
	// Application #1, instance 1
	inst, err = meErtCore.SystemModel.GetInstance("App#2-1-1")
	assert.NilError(t, err)
	rel, err = inst.GetReliability()
	assert.NilError(t, err)
	assert.Assert(t, (rel < app1.inst1[1]+deviation) && (rel > app1.inst1[1]-deviation))

	// Application #1, instance 2
	inst, err = meErtCore.SystemModel.GetInstance("App#2-1-2")
	assert.NilError(t, err)
	rel, err = inst.GetReliability()
	assert.NilError(t, err)
	assert.Assert(t, (rel < app1.inst2[1]+deviation) && (rel > app1.inst2[1]-deviation))

	// Application #1, instance 3
	inst, err = meErtCore.SystemModel.GetInstance("App#2-1-3")
	assert.NilError(t, err)
	rel, err = inst.GetReliability()
	assert.NilError(t, err)
	assert.Assert(t, (rel < app1.inst3[1]+deviation) && (rel > app1.inst3[1]-deviation))

	// Application #2, instance 1
	inst, err = meErtCore.SystemModel.GetInstance("App#2-2-1")
	assert.NilError(t, err)
	rel, err = inst.GetReliability()
	assert.NilError(t, err)
	assert.Assert(t, (rel < app2.inst1[1]+deviation) && (rel > app2.inst1[1]-deviation))

	// Application #2, instance 2
	inst, err = meErtCore.SystemModel.GetInstance("App#2-2-2")
	assert.NilError(t, err)
	rel, err = inst.GetReliability()
	assert.NilError(t, err)
	assert.Assert(t, (rel < app2.inst2[1]+deviation) && (rel > app2.inst2[1]-deviation))

	// iterating a bit more and checking critical points for us
	for i := 2; i <= 100; i++ {
		// setting reliabilities for each instance
		err := UpdateReliabilities(meErtCore.SystemModel, i, app1, app2, vi)
		assert.NilError(t, err)
		// Application #1, instance 1
		if i == 60 {
			inst, err := meErtCore.SystemModel.GetInstance("App#2-1-2")
			assert.NilError(t, err)
			rel, err := inst.GetReliability()
			assert.NilError(t, err)
			t.Logf("Reliability is %v", rel)
			assert.Assert(t, (rel < app1.inst2[60]+deviation) && (rel > app1.inst2[60]-deviation))
		}
		// Application #1, instance 2
		if i == 30 {
			inst, err := meErtCore.SystemModel.GetInstance("App#2-1-2")
			assert.NilError(t, err)
			rel, err := inst.GetReliability()
			assert.NilError(t, err)
			assert.Assert(t, (rel < 0.27+deviation) && (rel > 0.27-deviation))
		}

	}
}

func TestMeasurementDepth4(t *testing.T) {
	err := runMeasurementForDepth4(true)
	assert.NilError(t, err)
}

func TestMeasurementDepth3(t *testing.T) {
	err := runMeasurementForDepth3(true)
	assert.NilError(t, err)
}

func TestMeasurementDepth2(t *testing.T) {
	err := runMeasurementForDepth2(true)
	assert.NilError(t, err)
}

func TestGenerateRandomVectorOfLength(t *testing.T) {
	v := generateRandomVectorOfLength(0.5, 10)
	t.Logf("Generated vector is %v", v)
	t.Logf("Maximum float64 is %v", math.MaxFloat64)
}

func TestComputeMeErtCoreReliabilities(t *testing.T) {

	relMap := map[int]float64{
		1: 0.654876213,
		2: 0.1236854984,
		3: 0.46821546,
		4: 0.3569875654,
	}

	_, err := computeMeErtCoreCoefficients(relMap, 4)
	assert.ErrorContains(t, err, "obtained incomplete map")

	v := generateInputDataForInstance(app1inst1)
	t.Logf("Computed coefficients are: %v", v)
	coefs, err := computeMeErtCoreCoefficients(v, 4)
	assert.NilError(t, err)
	t.Logf("Computed coefficients are: %v", coefs)
}

func TestMeasurementWide(t *testing.T) {
	err := runMeasurementWide(10, 10, true)
	assert.NilError(t, err)
}

func TestMeasurementWide2(t *testing.T) {
	//deviation = 0.1 * deviation
	deviation = 0

	tc := make([]int, 0)
	tc = append(tc, 10)   // FMAIS of depth 3 with 10 Apps
	tc = append(tc, 1000) // FMAIS of depth 4 with 1000 Apps

	// initializing input data
	app, appFailed := InitializeInputDataWide()

	// iterating over test cases
	for _, val := range tc {

		// initialize FMAIS of depth 3 with 10 applications
		sm, err := systemmodel.CreateSystemModelWide(val)
		assert.NilError(t, err)

		meErtCore := meertcore.MeErtCore{
			SystemModel: sm,
			Reliability: 0.0,
		}

		///// step 1
		// setting reliabilities for each instance
		err = UpdateReliabilities(meErtCore.SystemModel, 1, appFailed, app)
		assert.NilError(t, err)

		_, err = meErtCore.SystemModel.GatherAllApplicationsReliabilities()
		assert.NilError(t, err)

		// computing reliability of the FMAIS per (optimized) ME-ERT-CORE
		rel, err := meErtCore.ComputeReliabilityOptimizedSimple()
		assert.NilError(t, err)

		t.Logf("Computed reliability is %v for FMAIS with %d Apps", rel, val)

		///// step 101
		// setting reliabilities for each instance
		err = UpdateReliabilities(meErtCore.SystemModel, 101, appFailed, app)
		assert.NilError(t, err)

		_, err = meErtCore.SystemModel.GatherAllApplicationsReliabilities()
		assert.NilError(t, err)

		// computing reliability of the FMAIS per (optimized) ME-ERT-CORE
		rel1, err := meErtCore.ComputeReliabilityOptimizedSimple()
		assert.NilError(t, err)

		t.Logf("Computed reliability (first failure) is %v for FMAIS with %d Apps", rel1, val)

		///// step 150
		// setting reliabilities for each instance
		err = UpdateReliabilities(meErtCore.SystemModel, 150, appFailed, app)
		assert.NilError(t, err)

		_, err = meErtCore.SystemModel.GatherAllApplicationsReliabilities()
		assert.NilError(t, err)

		// computing reliability of the FMAIS per (optimized) ME-ERT-CORE
		rel2, err := meErtCore.ComputeReliabilityOptimizedSimple()
		assert.NilError(t, err)

		t.Logf("Computed reliability is %v for FMAIS with %d Apps", rel2, val)

		///// step 170
		// setting reliabilities for each instance
		err = UpdateReliabilities(meErtCore.SystemModel, 170, appFailed, app)
		assert.NilError(t, err)

		_, err = meErtCore.SystemModel.GatherAllApplicationsReliabilities()
		assert.NilError(t, err)

		// computing reliability of the FMAIS per (optimized) ME-ERT-CORE
		rel3, err := meErtCore.ComputeReliabilityOptimizedSimple()
		assert.NilError(t, err)

		t.Logf("Computed reliability (second failure) is %v for FMAIS with %d Apps", rel3, val)

		///// step 200
		// setting reliabilities for each instance
		err = UpdateReliabilities(meErtCore.SystemModel, 200, appFailed, app)
		assert.NilError(t, err)

		_, err = meErtCore.SystemModel.GatherAllApplicationsReliabilities()
		assert.NilError(t, err)

		// computing reliability of the FMAIS per (optimized) ME-ERT-CORE
		rel4, err := meErtCore.ComputeReliabilityOptimizedSimple()
		assert.NilError(t, err)

		t.Logf("Computed reliability is %v for FMAIS with %d Apps", rel4, val)
	}
}

func TestSmWideBench(t *testing.T) {
	deviation = 0

	numApps := 100
	app, appFailed := InitializeInputDataWide()

	sm, err := systemmodel.CreateSystemModelWideBench(numApps, 2, 4)
	assert.NilError(t, err)
	assert.Equal(t, len(sm.Applications)-1, numApps)
	assert.Equal(t, sm.Depth, 4)
	//sm.PrettyPrintApplications().PrettyPrintLayers()

	err = UpdateReliabilities(sm, 101, appFailed, app)
	assert.NilError(t, err)

	_, err = sm.GatherAllApplicationsReliabilities()
	assert.NilError(t, err)
}

func TestSmWideBench2(t *testing.T) {
	deviation = 0.05

	numApps := 100
	app, appFailed := InitializeInputDataWide()

	sm, err := systemmodel.CreateSystemModelWideBench(numApps, 26, 4)
	assert.NilError(t, err)
	assert.Equal(t, len(sm.Applications)-1, numApps)

	err = UpdateReliabilities(sm, 101, appFailed, app)
	assert.NilError(t, err)

	_, err = sm.GatherAllApplicationsReliabilities()
	assert.NilError(t, err)
}
