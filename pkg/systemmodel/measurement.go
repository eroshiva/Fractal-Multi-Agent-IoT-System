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

	count := len(rls)

	counter := 0
	for d := len(sm.Layers); d > 1; d-- {
		layer, ok := sm.Layers[d]
		if !ok {
			sm.PrettyPrintApplications().PrettyPrintApplications()
			return fmt.Errorf("layer %d was not found in FMAIS", d)
		}
		for _, inst := range layer.Instances {
			if counter == count {
				break
			}
			instAppName, err := inst.GetAppName()
			if err != nil {
				return err
			}
			instNumber, err := inst.GetInstanceNumber()
			if err != nil {
				return err
			}
			if strings.EqualFold(appName, instAppName) {
				counter++
				// extract correct instance reliability
				instRel, xtrctd := rls[instNumber]
				if !xtrctd && !inst.IsVI() {
					return fmt.Errorf("can't extract reliability for instance %d from %v", instNumber, rls)
				}
				// re-try in case we hit VIaaS
				if inst.IsVI() {
					instRel, xtrctd = rls[1]
					if !xtrctd {
						return fmt.Errorf("can't extract reliability for VIaaS instance %s from %v", inst.Name, rls)
					}
				}
				// update instance's reliability
				inst.SetReliability(instRel)
			}
		}
	}

	if counter != count && !strings.HasPrefix(appName, "VI") {
		return fmt.Errorf("couldn't find all instances of the Application %s, found only %d, but expected %d", appName, counter, count)
	}

	return nil
}

// CreateSystemModelWideBench creates a FMAIS system model of depth 4 with 100 Applications
func CreateSystemModelWideBench(numApps, numAppInst, depth int) (*SystemModel, error) {
	viNumInst := 2      // number of instances per VI
	viNum := 0          // defining number of VIs on 2nd or 3rd layer (last layer before the layer with Applications only)
	viNotLastLayer := 0 // number of VIs residing not on a 2nd layer

	// defining an instance priority
	instancePriority := 1 / float64(numAppInst)

	// some pre-computations
	viPrior := 0.1                               // VIaaS priority
	appPrior := (1 - viPrior) / float64(numApps) // priority for each application

	// how many VIs do we need?
	chunkVI := 5 // defining a maximum number of Apps hosted by one VI
	viLastLayer := numApps / chunkVI

	// attaching all apps to the last, 4th, layer
	viNum = viLastLayer
	viNotLastLayer = viNum / viNumInst

	systemModel := &SystemModel{}
	systemModel.InitializeSystemModel(numApps, depth)
	systemModel.InitializeRootLayer()
	systemModel.CreateApplication(viNumInst, 1.0, "VI")
	systemModel.Applications["VI"].SetPriority(viPrior).Deploy()

	switch depth {
	case 2:
		// initializing layer 2
		layer2 := &Layer{
			VIwasDeployed: true,
		}
		layer2.InitializeLayer()
		// all applications are attached to the Root instance, FMAIS
		for i := 1; i <= numApps; i++ {
			appName := "App#" + fmt.Sprintf("%d", i)
			systemModel.CreateApplication(numAppInst, 1.0, appName)
			// setting priorities for each application
			systemModel.Applications[appName].SetPriority(appPrior).Deploy()
			for instance := 1; instance <= numAppInst; instance++ {
				appInst := &Instance{}
				appInst.CreateInstance(fmt.Sprintf("App#2-%d-%d", i, instance), CreateInstanceTypeApp()).SetPriority(instancePriority).SetReliability(0.77)
				// adding this new instance as a relation to the Root node
				systemModel.Layers[1].Instances[0].AddRelation(appInst)
				layer2.AddInstanceToLayer(appInst)
			}
		}
		systemModel.AddLayer(layer2, 2)
	case 3:
		// creating 2nd layer (VIs only) and creating 3rd layer (Apps only)
		// initializing layer 2 and layer 3
		layer2 := &Layer{
			VIwasDeployed: true,
		}
		layer2.InitializeLayer()

		layer3 := &Layer{
			VIwasDeployed: false,
		}
		layer3.InitializeLayer()

		for j := 1; j <= viNum; j++ {
			vi := &Instance{} // this is VIaaS
			vi.CreateInstance(fmt.Sprintf("VI#2-%d", j), CreateInstanceTypeVI()).SetPriority(1)
			// adding these new instances as a relation to the Root node
			systemModel.Layers[1].Instances[0].AddRelation(vi)

			// all applications are attached to the Root instance, FMAIS
			for i := (j-1)*chunkVI + 1; i <= chunkVI+(j-1)*chunkVI; i++ {
				appName := "App#" + fmt.Sprintf("%d", i)
				systemModel.CreateApplication(numAppInst, 1.0, appName)
				// setting priorities for each application
				systemModel.Applications[appName].SetPriority(appPrior).Deploy()
				for instance := 1; instance <= numAppInst; instance++ {
					appInst := &Instance{}
					appInst.CreateInstance(fmt.Sprintf("App#3-%d-%d", i, instance), CreateInstanceTypeApp()).SetPriority(instancePriority).SetReliability(0.77)
					// adding these new instances as a relation to the previously created VI
					vi.AddRelation(appInst)
					// adding instances to the 3rd layer
					layer3.AddInstanceToLayer(appInst)
				}
			}

			// adding instances to the 2nd layer
			layer2.AddInstanceToLayer(vi)
		}
		systemModel.AddLayer(layer2, 2)
		systemModel.AddLayer(layer3, 3)
	case 4:
		// creating 2nd and 3rd layer (VIs only), and creating 4th layer (Apps only)
		// initializing all layers
		layer2 := &Layer{
			VIwasDeployed: true,
		}
		layer2.InitializeLayer()

		layer3 := &Layer{
			VIwasDeployed: true,
		}
		layer3.InitializeLayer()

		layer4 := &Layer{
			VIwasDeployed: false,
		}
		layer4.InitializeLayer()

		for k := 1; k <= viNotLastLayer; k++ {
			viL2 := &Instance{} // this is VIaaS
			viL2.CreateInstance(fmt.Sprintf("VI#2-%d", k), CreateInstanceTypeVI()).SetPriority(1)
			// adding these new instances as a relation to the Root node
			systemModel.Layers[1].Instances[0].AddRelation(viL2)

			for j := (k-1)*viNumInst + 1; j <= viNumInst+(k-1)*viNumInst; j++ {
				vi := &Instance{} // this is VIaaS
				vi.CreateInstance(fmt.Sprintf("VI#3-%d", j), CreateInstanceTypeVI()).SetPriority(1)
				// adding these new instances as a relation to the Root node
				viL2.AddRelation(vi)

				// all applications are attached to the Root instance, FMAIS
				for i := (j-1)*chunkVI + 1; i <= chunkVI+(j-1)*chunkVI; i++ {
					appName := "App#" + fmt.Sprintf("%d", i)
					systemModel.CreateApplication(numAppInst, 1.0, appName)
					// setting priorities for each application
					systemModel.Applications[appName].SetPriority(appPrior).Deploy()
					for instance := 1; instance <= numAppInst; instance++ {
						appInst := &Instance{}
						appInst.CreateInstance(fmt.Sprintf("App#4-%d-%d", i, instance), CreateInstanceTypeApp()).SetPriority(instancePriority).SetReliability(0.77)
						// adding these new instances as a relation to the previously created VI
						vi.AddRelation(appInst)
						// adding instances to the 4th layer
						layer4.AddInstanceToLayer(appInst)
					}
				}

				// adding instances to the 2nd layer
				layer3.AddInstanceToLayer(vi)
			}
			layer2.AddInstanceToLayer(viL2)
		}
		systemModel.AddLayer(layer2, 2)
		systemModel.AddLayer(layer3, 3)
		systemModel.AddLayer(layer4, 4)
	default:
		return nil, fmt.Errorf("obtained unspecified depth: %d", depth)
	}

	return systemModel, nil
}

// CreateSystemModelWide creates a FMAIS system model of depth 4 with 100 Applications
func CreateSystemModelWide(numApps int) (*SystemModel, error) {
	depth := 4          // system model depth
	viNumInst := 2      // number of instances per VI
	viNum := 0          // defining number of VIs on 2nd or 3rd layer (last layer before the layer with Applications only)
	viNotLastLayer := 0 // number of VIs residing not on a 2nd layer

	numAppInst := 2 // number of instances per application

	// some pre-computations
	viPrior := 0.1                               // VIaaS priority
	appPrior := (1 - viPrior) / float64(numApps) // priority for each application

	// how many VIs do we need?
	chunkVI := 5 // defining a maximum number of Apps hosted by one VI
	viLastLayer := numApps / chunkVI

	// updating number of VIs on the last layer before applications and number of VIs on the layer above
	if viLastLayer < 2 {
		// attaching all apps to the root layer
		depth = 2
	} else if viLastLayer < 10 {
		// attaching all apps to the VIs deployed on the second layer
		viNum = viLastLayer
		depth = 3
	} else {
		// attaching all apps to the last, 4th, layer
		viNum = viLastLayer
		viNotLastLayer = viNum / viNumInst
	}

	systemModel := &SystemModel{}
	systemModel.InitializeSystemModel(numApps, depth)
	systemModel.InitializeRootLayer()
	systemModel.CreateApplication(viNumInst, 1.0, "VI")
	systemModel.Applications["VI"].SetPriority(viPrior).Deploy()

	switch depth {
	case 2:
		// initializing layer 2
		layer2 := &Layer{
			VIwasDeployed: true,
		}
		layer2.InitializeLayer()
		// all applications are attached to the Root instance, FMAIS
		for i := 1; i <= numApps; i++ {
			appName := "App#" + fmt.Sprintf("%d", i)
			systemModel.CreateApplication(numAppInst, 1.0, appName)
			// setting priorities for each application
			systemModel.Applications[appName].SetPriority(appPrior).Deploy()
			app11 := &Instance{}
			app11.CreateInstance(fmt.Sprintf("App#2-%d-1", i), CreateInstanceTypeApp()).SetPriority(0.41).SetReliability(0.77)
			app12 := &Instance{}
			app12.CreateInstance(fmt.Sprintf("App#2-%d-2", i), CreateInstanceTypeApp()).SetPriority(0.59).SetReliability(0.34)

			// adding these new instances as a relation to the Root node
			systemModel.Layers[1].Instances[0].AddRelation(app11).AddRelation(app12)

			// adding instances to the 2nd layer
			layer2.AddInstanceToLayer(app11).AddInstanceToLayer(app12)
		}
		systemModel.AddLayer(layer2, 2)
	case 3:
		// creating 2nd layer (VIs only) and creating 3rd layer (Apps only)
		// initializing layer 2 and layer 3
		layer2 := &Layer{
			VIwasDeployed: true,
		}
		layer2.InitializeLayer()

		layer3 := &Layer{
			VIwasDeployed: false,
		}
		layer3.InitializeLayer()

		for j := 1; j <= viNum; j++ {
			vi := &Instance{} // this is VIaaS
			vi.CreateInstance(fmt.Sprintf("VI#2-%d", j), CreateInstanceTypeVI()).SetPriority(1)
			// adding these new instances as a relation to the Root node
			systemModel.Layers[1].Instances[0].AddRelation(vi)

			// all applications are attached to the Root instance, FMAIS
			for i := (j-1)*chunkVI + 1; i <= chunkVI+(j-1)*chunkVI; i++ {
				appName := "App#" + fmt.Sprintf("%d", i)
				systemModel.CreateApplication(numAppInst, 1.0, appName)
				// setting priorities for each application
				systemModel.Applications[appName].SetPriority(appPrior).Deploy()
				app11 := &Instance{}
				app11.CreateInstance(fmt.Sprintf("App#3-%d-1", i), CreateInstanceTypeApp()).SetPriority(0.41).SetReliability(0.77)
				app12 := &Instance{}
				app12.CreateInstance(fmt.Sprintf("App#3-%d-2", i), CreateInstanceTypeApp()).SetPriority(0.59).SetReliability(0.34)

				// adding these new instances as a relation to the Root node
				vi.AddRelation(app11).AddRelation(app12)

				// adding instances to the 2nd layer
				layer3.AddInstanceToLayer(app11).AddInstanceToLayer(app12)
			}

			// adding instances to the 2nd layer
			layer2.AddInstanceToLayer(vi)
		}
		systemModel.AddLayer(layer2, 2)
		systemModel.AddLayer(layer3, 3)
	case 4:
		// creating 2nd and 3rd layer (VIs only), and creating 4th layer (Apps only)
		// initializing all layers
		layer2 := &Layer{
			VIwasDeployed: true,
		}
		layer2.InitializeLayer()

		layer3 := &Layer{
			VIwasDeployed: true,
		}
		layer3.InitializeLayer()

		layer4 := &Layer{
			VIwasDeployed: false,
		}
		layer4.InitializeLayer()

		for k := 1; k <= viNotLastLayer; k++ {
			viL2 := &Instance{} // this is VIaaS
			viL2.CreateInstance(fmt.Sprintf("VI#2-%d", k), CreateInstanceTypeVI()).SetPriority(1)
			// adding these new instances as a relation to the Root node
			systemModel.Layers[1].Instances[0].AddRelation(viL2)

			for j := (k-1)*viNumInst + 1; j <= viNumInst+(k-1)*viNumInst; j++ {
				vi := &Instance{} // this is VIaaS
				vi.CreateInstance(fmt.Sprintf("VI#3-%d", j), CreateInstanceTypeVI()).SetPriority(1)
				// adding these new instances as a relation to the Root node
				viL2.AddRelation(vi)

				// all applications are attached to the Root instance, FMAIS
				for i := (j-1)*chunkVI + 1; i <= chunkVI+(j-1)*chunkVI; i++ {
					appName := "App#" + fmt.Sprintf("%d", i)
					systemModel.CreateApplication(numAppInst, 1.0, appName)
					// setting priorities for each application
					systemModel.Applications[appName].SetPriority(appPrior).Deploy()
					app11 := &Instance{}
					app11.CreateInstance(fmt.Sprintf("App#4-%d-1", i), CreateInstanceTypeApp()).SetPriority(0.41).SetReliability(0.77)
					app12 := &Instance{}
					app12.CreateInstance(fmt.Sprintf("App#4-%d-2", i), CreateInstanceTypeApp()).SetPriority(0.59).SetReliability(0.34)

					// adding these new instances as a relation to the Root node
					vi.AddRelation(app11).AddRelation(app12)

					// adding instances to the 2nd layer
					layer4.AddInstanceToLayer(app11).AddInstanceToLayer(app12)
				}

				// adding instances to the 2nd layer
				layer3.AddInstanceToLayer(vi)
			}
			layer2.AddInstanceToLayer(viL2)
		}
		systemModel.AddLayer(layer2, 2)
		systemModel.AddLayer(layer3, 3)
		systemModel.AddLayer(layer4, 4)
	default:
		return nil, fmt.Errorf("obtained unspecified depth: %d", depth)
	}

	return systemModel, nil
}
