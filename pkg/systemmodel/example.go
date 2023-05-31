// Package systemmodel implements basic functionality to compose a Fractal MAIS system. This file implements a utility
// function, which defines a very basic example of a Fractal MAIS system.
package systemmodel

// CreateExampleBasicFMAIS implements an utility function which creates a very basic Fractal MAIS system
// with 4 VIs and 2 Applications. Here is a detailed description
// - Second level deploys two VIs and Application#1 (of 3 instances)
// - Third level deploys:
//   - two VIs, which are connected to the second VI on the second level
//   - Application#2 (of 5 instances) which is connected to the first VI on the second level
func CreateExampleBasicFMAIS() *SystemModel {
	systemModel := &SystemModel{}
	systemModel.InitializeSystemModel(10, 4)
	systemModel.InitializeRootLayer()
	systemModel.CreateApplication(2, 0.4, "VI").
		CreateApplication(3, 0.3, "App#1").
		CreateApplication(5, 0.3, "App#2")

	// setting priorities for each application
	systemModel.Applications["VI"].SetPriority(0.35)
	systemModel.Applications["App#1"].SetPriority(0.25)
	systemModel.Applications["App#2"].SetPriority(0.4)

	// creating two VI instances and App#1 on layer 2
	vi1 := &Instance{}
	vi1.CreateInstance("VI#2-1", CreateInstanceTypeVI()).SetPriority(0.25) // we don't set reliability here, cause this instance deploys other instances
	vi2 := &Instance{}
	vi2.CreateInstance("VI#2-2", CreateInstanceTypeVI()).SetPriority(0.25) // we don't set reliability here, cause this instance deploys other instances
	app11 := &Instance{}
	app11.CreateInstance("App#2-1-1", CreateInstanceTypeApp()).SetPriority(0.2).SetReliability(0.77)
	app12 := &Instance{}
	app12.CreateInstance("App#2-1-2", CreateInstanceTypeApp()).SetPriority(0.5).SetReliability(0.34)
	app13 := &Instance{}
	app13.CreateInstance("App#2-1-3", CreateInstanceTypeApp()).SetPriority(0.3).SetReliability(0.62)

	// adding these new instances as a relation to the Root node
	systemModel.Layers[1].Instances[0].AddRelation(vi1).AddRelation(vi2).AddRelation(app11).
		AddRelation(app12).AddRelation(app13)

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
	app21.CreateInstance("App#3-2-1", CreateInstanceTypeApp()).SetPriority(0.2).SetReliability(0.47)
	app22 := &Instance{}
	app22.CreateInstance("App#3-2-2", CreateInstanceTypeApp()).SetPriority(0.2).SetReliability(0.39)
	app23 := &Instance{}
	app23.CreateInstance("App#3-2-3", CreateInstanceTypeApp()).SetPriority(0.2).SetReliability(0.53)
	app24 := &Instance{}
	app24.CreateInstance("App#3-2-4", CreateInstanceTypeApp()).SetPriority(0.2).SetReliability(0.45)
	app25 := &Instance{}
	app25.CreateInstance("App#3-2-5", CreateInstanceTypeApp()).SetPriority(0.2).SetReliability(0.74)
	vi3 := &Instance{}
	vi3.CreateInstance("VI#3-3", CreateInstanceTypeVI()).SetPriority(0.25).SetReliability(0.61)
	vi4 := &Instance{}
	vi4.CreateInstance("VI#3-4", CreateInstanceTypeVI()).SetPriority(0.25).SetReliability(0.7)

	// adding new instances as a relations
	vi1.AddRelation(app21).AddRelation(app22).AddRelation(app23).AddRelation(app24).AddRelation(app25)
	vi2.AddRelation(vi3).AddRelation(vi4)

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
	var viCount uint64 = 2
	systemModel.VIcount = &viCount

	return systemModel
}
