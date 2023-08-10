package systemmodel

import (
	"gotest.tools/assert"
	"testing"
)

func TestGatherApplicationReliabilities(t *testing.T) {
	// creating a sample System Model with two VIs and two Applications running
	systemModel := CreateExampleBasicFMAIS()
	t.Logf("System model is\n%v", systemModel)

	systemModel.SetApplicationPrioritiesRandom()
	systemModel.PrettyPrintApplications().PrettyPrintLayers()
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

	systemModel.PrettyPrintApplications()
	systemModel.PrettyPrintLayers()
}
