package meertcore

import (
	"gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/pkg/systemmodel"
	"gotest.tools/assert"
	"testing"
)

func TestTotalReliability(t *testing.T) {
	// creating a sample System Model with two VIs and two Applications running
	systemModel := systemmodel.CreateExampleBasicFMAS()
	t.Logf("System model is\n%v", systemModel)

	systemModel.SetApplicationPrioritiesRandom()
	systemModel.PrettyPrintLayers()
	err := systemModel.SetInstancePrioritiesRandom()
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
