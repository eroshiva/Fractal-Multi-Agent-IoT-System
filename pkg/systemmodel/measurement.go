// Package systemmodel holds means of FMAIS system model generation. This file in particular carries all necessary functionality
// for measurement
package systemmodel

import (
	"fmt"
	"strings"
)

// CreateSystemModelDepth4 creates a FMAIS system model of depth 4 with 4 Applications
// with maximum 3 number of instances per Application
func CreateSystemModelDepth4() *SystemModel {
	systemModel := &SystemModel{}
	systemModel.InitializeSystemModel(4, 4)
	systemModel.InitializeRootLayer()
	systemModel.CreateApplication(2, 1.0, "VI").
		CreateApplication(3, 1.0, "App#1").
		CreateApplication(2, 1.0, "App#2").
		CreateApplication(2, 1.0, "App#3").
		CreateApplication(1, 1.0, "App#4")

	// setting priorities for each application
	systemModel.Applications["VI"].SetPriority(0.35).Deploy()
	systemModel.Applications["App#1"].SetPriority(0.19).Deploy()
	systemModel.Applications["App#2"].SetPriority(0.15).Deploy()
	systemModel.Applications["App#3"].SetPriority(0.21).Deploy()
	systemModel.Applications["App#4"].SetPriority(0.1).Deploy()

	// creating two VI instances and App#1 on layer 2
	// we DO care about relations. This is not important to compute reliability (yet)
	vi1 := &Instance{}
	vi1.CreateInstance("VI#2-1", CreateInstanceTypeVI()).SetPriority(0.25) // we don't set reliability here, cause this instance deploys other instances
	vi2 := &Instance{}
	vi2.CreateInstance("VI#2-2", CreateInstanceTypeVI()).SetPriority(0.25) // we don't set reliability here, cause this instance deploys other instances
	app11 := &Instance{}
	app11.CreateInstance("App#2-1-1", CreateInstanceTypeApp()).SetPriority(0.41).SetReliability(0.77)
	app12 := &Instance{}
	app12.CreateInstance("App#2-1-2", CreateInstanceTypeApp()).SetPriority(0.28).SetReliability(0.34)
	app13 := &Instance{}
	app13.CreateInstance("App#2-1-3", CreateInstanceTypeApp()).SetPriority(0.31).SetReliability(0.62)

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

	// creating layer 3 with App#2, App#3, one VIaaS and one other VI
	app21 := &Instance{}
	app21.CreateInstance("App#3-2-1", CreateInstanceTypeApp()).SetPriority(0.35).SetReliability(0.77)
	app22 := &Instance{}
	app22.CreateInstance("App#3-2-2", CreateInstanceTypeApp()).SetPriority(0.65).SetReliability(0.34)
	app31 := &Instance{}
	app31.CreateInstance("App#3-3-1", CreateInstanceTypeApp()).SetPriority(0.7).SetReliability(0.77)
	app32 := &Instance{}
	app32.CreateInstance("App#3-3-2", CreateInstanceTypeApp()).SetPriority(0.3).SetReliability(0.34)
	vi3 := &Instance{}
	vi3.CreateInstance("VI#3-3", CreateInstanceTypeVI()).SetPriority(0.25) // we don't set reliability here, cause this instance deploys other instances
	vi4 := &Instance{}                                                     // this is VIaaS
	vi4.CreateInstance("VI#3-4", CreateInstanceTypeVI()).SetPriority(1).SetReliability(0.45)

	// adding these new instances as a relation to the instances from the level above
	vi1.AddRelation(app21).AddRelation(app22)
	vi2.AddRelation(app31).AddRelation(app32).AddRelation(vi3).AddRelation(vi4)

	// initializing layer 3
	layer3 := &Layer{
		VIwasDeployed: true,
	}
	layer3.InitializeLayer()
	layer3.AddInstanceToLayer(vi3).AddInstanceToLayer(vi4).AddInstanceToLayer(app21).
		AddInstanceToLayer(app22).AddInstanceToLayer(app31).AddInstanceToLayer(app32)
	systemModel.AddLayer(layer3, 3)

	// creating layer 4 with App#4 and no other instance
	app4 := &Instance{}
	app4.CreateInstance("App#4-4-1", CreateInstanceTypeApp()).SetPriority(1).SetReliability(0.77)
	// adding this new instance to the relation of VI#3
	vi3.AddRelation(app4)

	// initializing layer 3
	layer4 := &Layer{
		VIwasDeployed: false,
	}
	layer4.InitializeLayer()
	layer4.AddInstanceToLayer(app4)
	systemModel.AddLayer(layer4, 4)

	return systemModel
}

// CreateSystemModelDepth3 creates a FMAIS system model of depth 3 with 3 Applications
// with maximum 3 number of instances per Application
func CreateSystemModelDepth3() *SystemModel {
	systemModel := &SystemModel{}
	systemModel.InitializeSystemModel(3, 3)
	systemModel.InitializeRootLayer()
	systemModel.CreateApplication(2, 1.0, "VI").
		CreateApplication(3, 1.0, "App#1").
		CreateApplication(2, 1.0, "App#2").
		CreateApplication(2, 1.0, "App#3")

	// setting priorities for each application
	systemModel.Applications["VI"].SetPriority(0.35).Deploy()
	systemModel.Applications["App#1"].SetPriority(0.27).Deploy()
	systemModel.Applications["App#2"].SetPriority(0.18).Deploy()
	systemModel.Applications["App#3"].SetPriority(0.2).Deploy()

	// creating two VI instances and App#1 on layer 2
	// we DO care about relations. This is not important to compute reliability (yet)
	vi1 := &Instance{}
	vi1.CreateInstance("VI#2-1", CreateInstanceTypeVI()).SetPriority(0.25) // we don't set reliability here, cause this instance deploys other instances
	vi2 := &Instance{}
	vi2.CreateInstance("VI#2-2", CreateInstanceTypeVI()).SetPriority(0.25) // we don't set reliability here, cause this instance deploys other instances
	app11 := &Instance{}
	app11.CreateInstance("App#2-1-1", CreateInstanceTypeApp()).SetPriority(0.41).SetReliability(0.77)
	app12 := &Instance{}
	app12.CreateInstance("App#2-1-2", CreateInstanceTypeApp()).SetPriority(0.28).SetReliability(0.34)
	app13 := &Instance{}
	app13.CreateInstance("App#2-1-3", CreateInstanceTypeApp()).SetPriority(0.31).SetReliability(0.62)

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

	// creating layer 3 with App#2 and App#3, and one VIaaS
	app21 := &Instance{}
	app21.CreateInstance("App#3-2-1", CreateInstanceTypeApp()).SetPriority(0.35).SetReliability(0.77)
	app22 := &Instance{}
	app22.CreateInstance("App#3-2-2", CreateInstanceTypeApp()).SetPriority(0.65).SetReliability(0.34)
	app31 := &Instance{}
	app31.CreateInstance("App#3-3-1", CreateInstanceTypeApp()).SetPriority(0.7).SetReliability(0.77)
	app32 := &Instance{}
	app32.CreateInstance("App#3-3-2", CreateInstanceTypeApp()).SetPriority(0.3).SetReliability(0.34)
	vi3 := &Instance{} // this is VIaaS
	vi3.CreateInstance("VI#3-3", CreateInstanceTypeVI()).SetPriority(1).SetReliability(0.45)

	// adding these new instances as a relation to the instances from the level above
	vi1.AddRelation(app21).AddRelation(app22)
	vi2.AddRelation(app31).AddRelation(app32).AddRelation(vi3)

	// initializing layer 3
	layer3 := &Layer{
		VIwasDeployed: true,
	}
	layer3.InitializeLayer()
	layer3.AddInstanceToLayer(vi3).AddInstanceToLayer(app21).
		AddInstanceToLayer(app22).AddInstanceToLayer(app31).AddInstanceToLayer(app32)
	systemModel.AddLayer(layer3, 3)

	return systemModel
}

// CreateSystemModelDepth2 creates a FMAIS system model of depth 2 with 2 Applications
// with maximum 3 number of instances per Application
func CreateSystemModelDepth2() *SystemModel {
	systemModel := &SystemModel{}
	systemModel.InitializeSystemModel(2, 2)
	systemModel.InitializeRootLayer()
	systemModel.CreateApplication(2, 1.0, "VI").
		CreateApplication(3, 1.0, "App#1").
		CreateApplication(2, 1.0, "App#2")

	// setting priorities for each application
	systemModel.Applications["VI"].SetPriority(0.35).Deploy()
	systemModel.Applications["App#1"].SetPriority(0.27).Deploy()
	systemModel.Applications["App#2"].SetPriority(0.38).Deploy()

	// creating VIaaS, App#1 and App#2 on layer 2
	// we DO care about relations. This is not important to compute reliability (yet)
	app11 := &Instance{}
	app11.CreateInstance("App#2-1-1", CreateInstanceTypeApp()).SetPriority(0.41).SetReliability(0.77)
	app12 := &Instance{}
	app12.CreateInstance("App#2-1-2", CreateInstanceTypeApp()).SetPriority(0.28).SetReliability(0.34)
	app13 := &Instance{}
	app13.CreateInstance("App#2-1-3", CreateInstanceTypeApp()).SetPriority(0.31).SetReliability(0.62)
	app21 := &Instance{}
	app21.CreateInstance("App#2-2-1", CreateInstanceTypeApp()).SetPriority(0.35).SetReliability(0.77)
	app22 := &Instance{}
	app22.CreateInstance("App#2-2-2", CreateInstanceTypeApp()).SetPriority(0.65).SetReliability(0.34)
	vi1 := &Instance{} // this is VIaaS
	vi1.CreateInstance("VI#2-1", CreateInstanceTypeVI()).SetPriority(1).SetReliability(0.45)

	// adding these new instances as a relation to the Root node
	systemModel.Layers[1].Instances[0].AddRelation(vi1).AddRelation(app11).
		AddRelation(app12).AddRelation(app13).AddRelation(app21).AddRelation(app22)

	// initializing layer 2
	layer2 := &Layer{
		VIwasDeployed: true,
	}
	layer2.InitializeLayer()
	layer2.AddInstanceToLayer(vi1).AddInstanceToLayer(app11).AddInstanceToLayer(app12).AddInstanceToLayer(app13).
		AddInstanceToLayer(app21).AddInstanceToLayer(app22)
	systemModel.AddLayer(layer2, 2)

	return systemModel
}

// UpdateApplicationReliability updates reliability values of all instances of a provided Application
func (sm *SystemModel) UpdateApplicationReliability(appName string, rls map[int64]float64) error {
	err := sm.SetApplicationReliability(appName, rls)
	if err != nil {
		return err
	}

	return nil
}

// SetApplicationReliability sets reliability values of all instances of a provided Application
func (sm *SystemModel) SetApplicationReliability(appName string, rls map[int64]float64) error {

	for d := len(sm.Layers); d > 1; d-- {
		layer, ok := sm.Layers[d]
		if !ok {
			sm.PrettyPrintApplications().PrettyPrintApplications()
			return fmt.Errorf("layer %d was not found in FMAIS", d)
		}
		for _, inst := range layer.Instances {
			instAppName, err := inst.GetAppName()
			if err != nil {
				return err
			}
			instNumber, err := inst.GetInstanceNumber()
			if err != nil {
				return err
			}
			if strings.EqualFold(appName, instAppName) {
				// extract correct instance reliability
				instRel, xtrctd := rls[instNumber]
				if !xtrctd {
					return fmt.Errorf("can't extract reliability for instance %d from %v", instNumber, rls)
				}
				// update instance's reliability
				inst.SetReliability(instRel)
			}
		}
	}

	return nil
}
