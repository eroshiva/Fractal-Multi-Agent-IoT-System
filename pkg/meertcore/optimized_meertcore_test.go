package meertcore

import (
	"gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/pkg/systemmodel"
	"gotest.tools/assert"
	"testing"
)

func TestTotalReliability(t *testing.T) {
	// creating a sample System Model with two VIs and two Applications running
	systemModel := systemmodel.CreateExampleBasicFMAIS()
	err := systemModel.SetChainCoefficients()
	assert.NilError(t, err)
	t.Logf("System model is\n%v", systemModel)

	systemModel.SetApplicationPrioritiesRandom()
	systemModel.PrettyPrintApplications().PrettyPrintLayers()
	err = systemModel.SetInstancePrioritiesRandom()
	assert.NilError(t, err)
	err = systemModel.SetInstanceReliabilitiesRandom()
	assert.NilError(t, err)

	// get reliabilities for instances within app
	appName := "App#1"
	reliabilities, err := systemModel.GatherApplicationInstanceReliabilities(appName)
	assert.NilError(t, err)
	assert.Assert(t, len(reliabilities) > 0)
	t.Logf("Reliabilities for application %v are\n%v\n", appName, reliabilities)

	// get reliabilities for all apps
	appRel, err := systemModel.GatherAllApplicationsReliabilities()
	assert.NilError(t, err)
	assert.Assert(t, len(appRel) > 0)
	t.Logf("Reliabilities for each application are\n%v\n", appRel)

	// Compute total reliability of the system
	me := &MeErtCore{
		SystemModel: systemModel,
		Reliability: -1.23456789,
	}
	totalRel, err := me.ComputeReliabilityOptimized()
	assert.NilError(t, err)
	assert.Assert(t, totalRel != -1.23456789)
	t.Logf("Total reliability of the system is: %v\n", totalRel)

	systemModel.PrettyPrintApplications()
	systemModel.PrettyPrintLayers()
}

func TestComputeReliabilityOptimizedSimpleDepth4(t *testing.T) {
	// initialising system model
	sm4 := systemmodel.CreateSystemModelDepth4()

	// initializing ME-ERT-CORE
	meErtCore := MeErtCore{
		SystemModel: sm4,
		Reliability: 0.0,
	}

	// get reliabilities for all apps
	appRel, err := meErtCore.SystemModel.GatherAllApplicationsReliabilities()
	assert.NilError(t, err)
	assert.Assert(t, len(appRel) > 0)
	t.Logf("Reliabilities for each application are\n%v\n", appRel)

	rel, err := meErtCore.ComputeReliabilityOptimizedSimple()
	assert.NilError(t, err)
	t.Logf("Computed reliability is %v/%v", rel, meErtCore.Reliability)
	assert.Assert(t, rel-0.557274 < 0.000000001)
}

func TestComputeReliabilityOptimizedSimpleDepth3(t *testing.T) {
	// initialising system model
	sm3 := systemmodel.CreateSystemModelDepth3()

	// initializing ME-ERT-CORE
	meErtCore := MeErtCore{
		SystemModel: sm3,
		Reliability: 0.0,
	}

	// get reliabilities for all apps
	appRel, err := meErtCore.SystemModel.GatherAllApplicationsReliabilities()
	assert.NilError(t, err)
	assert.Assert(t, len(appRel) > 0)
	t.Logf("Reliabilities for each application are\n%v\n", appRel)

	rel, err := meErtCore.ComputeReliabilityOptimizedSimple()
	assert.NilError(t, err)
	t.Logf("Computed reliability is %v/%v", rel, meErtCore.Reliability)
	assert.Assert(t, rel-0.536827 < 0.000000001)
}

func TestComputeReliabilityOptimizedSimpleDepth2(t *testing.T) {
	// initialising system model
	sm2 := systemmodel.CreateSystemModelDepth2()

	// initializing ME-ERT-CORE
	meErtCore := MeErtCore{
		SystemModel: sm2,
		Reliability: 0.0,
	}

	// get reliabilities for all apps
	appRel, err := meErtCore.SystemModel.GatherAllApplicationsReliabilities()
	assert.NilError(t, err)
	assert.Assert(t, len(appRel) > 0)
	t.Logf("Reliabilities for each application are\n%v\n", appRel)

	rel, err := meErtCore.ComputeReliabilityOptimizedSimple()
	assert.NilError(t, err)
	t.Logf("Computed reliability is %v/%v", rel, meErtCore.Reliability)
	assert.Assert(t, rel-0.506727 < 0.000000001)
}

func TestComputeMeErtCoreCoefficient(t *testing.T) {
	relVal := 0.54893654512
	coef, err := ComputeMeErtCoreCoefficient(relVal, 9)
	assert.NilError(t, err)
	assert.Equal(t, coef, relVal)

	coef, err = ComputeMeErtCoreCoefficient(relVal, 11)
	assert.NilError(t, err)
	assert.Equal(t, coef, 0.4893654512)

	coef, err = ComputeMeErtCoreCoefficient(relVal, 99)
	assert.NilError(t, err)
	assert.Equal(t, coef, 0.4893654512)

	coef, err = ComputeMeErtCoreCoefficient(relVal, 999)
	assert.NilError(t, err)
	assert.Equal(t, coef, 0.893654512)

	coef, err = ComputeMeErtCoreCoefficient(relVal, 9999)
	assert.NilError(t, err)
	assert.Equal(t, coef, 0.93654512)

	coef, err = ComputeMeErtCoreCoefficient(relVal, 99999)
	assert.NilError(t, err)
	assert.Equal(t, coef, relVal)
}
