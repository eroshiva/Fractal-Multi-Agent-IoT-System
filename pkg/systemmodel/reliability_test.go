package systemmodel

import (
	"gotest.tools/assert"
	"testing"
)

func TestGatherApplicationReliabilities(t *testing.T) {
	// creating (by hand) System Model with two VIs and two Applications running
	systemModel := &SystemModel{}

	//names := GenerateAppNames(2)
	systemModel.InitializeSystemModel(10, 4)
	systemModel.InitializeRootLayer()
	systemModel.CreateApplication(2, 0.4, "VI").
		CreateApplication(3, 0.3, "App#1").
		CreateApplication(5, 0.3, "App#2")

	// creating two VI instances and App#1 on layer 2
	// we do NOT care about relations. This is not important to compute reliability (yet)
	vi1 := &Instance{}
	vi1.CreateInstance("VI#1-2", CreateInstanceTypeVI())
	vi2 := &Instance{}
	vi2.CreateInstance("VI#2-2", CreateInstanceTypeVI())
	app11 := &Instance{}
	app11.CreateInstance("App#1-1-2", CreateInstanceTypeApp())
	app12 := &Instance{}
	app12.CreateInstance("App#1-2-2", CreateInstanceTypeApp())
	app13 := &Instance{}
	app13.CreateInstance("App#1-3-2", CreateInstanceTypeApp())

	// initializing layer 2
	layer2 := &Layer{
		VIwasDeployed: true,
	}
	layer2.InitializeLayer()
	layer2.AddInstanceToLayer(vi1).AddInstanceToLayer(vi2).AddInstanceToLayer(app11).
		AddInstanceToLayer(app12).AddInstanceToLayer(app13)
	systemModel.AddLayer(layer2, 2)

	// creating App#2 and one VI
	app21 := &Instance{}
	app21.CreateInstance("App#2-1-3", CreateInstanceTypeApp())
	app22 := &Instance{}
	app22.CreateInstance("App#2-2-3", CreateInstanceTypeApp())
	app23 := &Instance{}
	app23.CreateInstance("App#2-3-3", CreateInstanceTypeApp())
	app24 := &Instance{}
	app24.CreateInstance("App#2-4-3", CreateInstanceTypeApp())
	app25 := &Instance{}
	app25.CreateInstance("App#2-5-3", CreateInstanceTypeApp())
	vi3 := &Instance{}
	vi3.CreateInstance("VI#3-3", CreateInstanceTypeVI())
	vi4 := &Instance{}
	vi4.CreateInstance("VI#4-3", CreateInstanceTypeVI())

	// initializing layer 3
	layer3 := &Layer{
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

	systemModel.PrettyPrintApplications()
	systemModel.PrettyPrintLayers()
}
