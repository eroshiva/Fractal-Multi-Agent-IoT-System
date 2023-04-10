package meertcore

import (
	"fmt"
	"gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/pkg/systemmodel"
	"gotest.tools/assert"
	"testing"
)

func TestSingleInstanceFMAIS(t *testing.T) {
	systemModel := &systemmodel.SystemModel{}

	systemModel.InitializeSystemModel(10, 4)
	systemModel.InitializeRootLayer()

	err := systemModel.SetInstancePrioritiesRandom()
	assert.NilError(t, err)
	err = systemModel.SetInstanceReliabilitiesRandom()
	assert.NilError(t, err)

	// compute reliability
	me := MeErtCore{
		SystemModel: systemModel,
		Reliability: -1.23456789,
	}

	totalRel, err := me.ComputeReliabilityPerDefinition()
	assert.NilError(t, err)
	t.Logf("Computed reliability is %v", totalRel)

	cr, err := me.SystemModel.Layers[1].Instances[0].GetReliability()
	assert.NilError(t, err)
	t.Logf("Computed reliability is %v", cr)
}

func TestComputeReliabilityPerDefinition(t *testing.T) {
	// creating a sample System Model with two VIs and two Applications running
	systemModel := systemmodel.CreateExampleBasicFMAIS()
	t.Logf("System model is\n%v", systemModel)

	// Compute total reliability of the system
	me := &MeErtCore{
		SystemModel: systemModel,
		Reliability: -1.23456789,
	}
	totalRel, err := me.ComputeReliabilityPerDefinition()
	assert.NilError(t, err)
	assert.Assert(t, totalRel != -1.23456789)
	assert.Equal(t, fmt.Sprintf("%.12f", totalRel), "0.155589687500")

	systemModel.PrettyPrintApplications()
	systemModel.PrettyPrintLayers()
	t.Logf("Total reliability of the system is: %v\n", fmt.Sprintf("%.12f", totalRel))
}

func TestComputeReliabilityOptimized(t *testing.T) {
	// creating a sample System Model with two VIs and two Applications running
	systemModel := systemmodel.CreateExampleBasicFMAIS()
	err := systemModel.SetChainCoefficients()
	assert.NilError(t, err)
	t.Logf("System model is\n%v", systemModel)

	// Compute total reliability of the system
	me := &MeErtCore{
		SystemModel: systemModel,
		Reliability: -1.23456789,
	}
	totalRel, err := me.ComputeReliabilityPerDefinition()
	assert.NilError(t, err)
	assert.Assert(t, totalRel != -1.23456789)
	assert.Equal(t, fmt.Sprintf("%.12f", totalRel), "0.155589687500")
	t.Logf("Total reliability of the system is: %v\n", fmt.Sprintf("%.12f", totalRel))

	// computing reliability with optimized algorithm
	// some pre-processing
	// get reliabilities for all apps
	appRel, err := systemModel.GatherAllApplicationsReliabilities()
	assert.NilError(t, err)
	assert.Assert(t, len(appRel) > 0)
	t.Logf("Reliabilities for each application are\n%v\n", appRel)

	// computing reliability per optimized algorithm
	totalRel1, err1 := me.ComputeReliabilityOptimized()
	assert.NilError(t, err1)
	t.Logf("Computed (with optimized algo) reliability is %v\n", totalRel1)

	systemModel.PrettyPrintApplications()
	systemModel.PrettyPrintLayers()
}
