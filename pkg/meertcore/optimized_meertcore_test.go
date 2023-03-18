package meertcore

import (
	"gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/pkg/systemmodel"
	"gotest.tools/assert"
	"testing"
)

func TestTotalReliability(t *testing.T) {

	// creating (by hand) System Model with two VIs and two Applications running
	systemModel := &systemmodel.SystemModel{}

	//names := GenerateAppNames(2)
	systemModel.InitializeSystemModel(10, 4)
	systemModel.InitializeRootLayer()
	systemModel.CreateApplication(2, 0.4, "VI").
		CreateApplication(3, 0.3, "App#1").
		CreateApplication(5, 0.3, "App#2")

	// creating two VI instances and App#1 on layer 2
	// we do NOT care about relations. This is not important to compute reliability (yet)
	vi1 := &systemmodel.Instance{}
	vi1.CreateInstance("VI#1-2", systemmodel.CreateInstanceTypeVI())
	vi2 := &systemmodel.Instance{}
	vi2.CreateInstance("VI#2-2", systemmodel.CreateInstanceTypeVI())
	app11 := &systemmodel.Instance{}
	app11.CreateInstance("App#1-1-2", systemmodel.CreateInstanceTypeApp())
	app12 := &systemmodel.Instance{}
	app12.CreateInstance("App#1-2-2", systemmodel.CreateInstanceTypeApp())
	app13 := &systemmodel.Instance{}
	app13.CreateInstance("App#1-3-2", systemmodel.CreateInstanceTypeApp())

	// initializing layer 2
	layer2 := &systemmodel.Layer{
		VIwasDeployed: true,
	}
	layer2.InitializeLayer()
	layer2.AddInstanceToLayer(vi1).AddInstanceToLayer(vi2).AddInstanceToLayer(app11).
		AddInstanceToLayer(app12).AddInstanceToLayer(app13)
	systemModel.AddLayer(layer2, 2)

	// creating App#2 and one VI
	app21 := &systemmodel.Instance{}
	app21.CreateInstance("App#2-1-3", systemmodel.CreateInstanceTypeApp())
	app22 := &systemmodel.Instance{}
	app22.CreateInstance("App#2-2-3", systemmodel.CreateInstanceTypeApp())
	app23 := &systemmodel.Instance{}
	app23.CreateInstance("App#2-3-3", systemmodel.CreateInstanceTypeApp())
	app24 := &systemmodel.Instance{}
	app24.CreateInstance("App#2-4-3", systemmodel.CreateInstanceTypeApp())
	app25 := &systemmodel.Instance{}
	app25.CreateInstance("App#2-5-3", systemmodel.CreateInstanceTypeApp())
	vi3 := &systemmodel.Instance{}
	vi3.CreateInstance("VI#3-3", systemmodel.CreateInstanceTypeVI())
	vi4 := &systemmodel.Instance{}
	vi4.CreateInstance("VI#4-3", systemmodel.CreateInstanceTypeVI())

	// initializing layer 3
	layer3 := &systemmodel.Layer{
		VIwasDeployed: true,
	}
	layer3.InitializeLayer()
	layer3.AddInstanceToLayer(app21).AddInstanceToLayer(app22).AddInstanceToLayer(app23).
		AddInstanceToLayer(app24).AddInstanceToLayer(app25).AddInstanceToLayer(vi3).AddInstanceToLayer(vi4)
	systemModel.AddLayer(layer3, 3)

	// setting deployed state of all Apps and VI to true
	systemModel.Applications["VI"].State = true
	systemModel.Applications["App#1"].State = true
	systemModel.Applications["App#2"].State = true
	var viCount uint64 = 4
	systemModel.VIcount = &viCount

	t.Logf("System model is\n%v", systemModel)

	systemModel.SetApplicationPrioritiesRandom()
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
		Reliability: 0,
	}
	totalRel, err := me.ComputeReliabilityOptimized()
	assert.NilError(t, err)
	assert.Assert(t, totalRel != 0.0)
	t.Logf("Total reliability of the system is: %v\n", totalRel)

	systemModel.PrettyPrintApplications()
	systemModel.PrettyPrintLayers()
}
