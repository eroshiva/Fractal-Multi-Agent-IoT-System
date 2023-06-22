// Package measurement provides a measurement logic and all helper functions.
package measurement

import (
	"fmt"
	"gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/pkg/draw"
	"gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/pkg/meertcore"
	"gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/pkg/storedata"
	"gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/pkg/systemmodel"
	"log"
	"math/rand"
	"strconv"
)

var maxAppNumWide = 1000
var stepWide = 10

// inputData is a structure to define input data values in a defined range
type inputData struct {
	from  int     // from which step
	to    int     // to which step
	value float64 // mean value in the defined interval
}

// app1 structure holds input data for Application #1
type app1 struct {
	inst1 map[int]float64
	inst2 map[int]float64
	inst3 map[int]float64
}

// app2 structure holds input data for Application #2
type app2 struct {
	inst1 map[int]float64
	inst2 map[int]float64
}

// app1 structure holds input data for Application #3
type app3 struct {
	inst1 map[int]float64
	inst2 map[int]float64
}

// app4 structure holds input data for Application #4
type app4 struct {
	inst1 map[int]float64
}

// vi structure holds input data for VIaaS
type vi struct {
	inst map[int]float64
}

// appNoFail structure holds input data for Application (which does NOT fail)
type appNoFail struct {
	inst1 map[int]float64
	inst2 map[int]float64
}

// appFail structure holds input data for Application (which FAILS)
type appFail struct {
	inst1 map[int]float64
	inst2 map[int]float64
}

// RunMeasurement function initializes and runs measurement for all FMAIS depths
func RunMeasurement() error {
	// run measurement for FMAIS of depth 4
	err := runMeasurementForDepth4(false)
	if err != nil {
		return err
	}

	// run measurement for FMAIS of depth 3
	err = runMeasurementForDepth3(false)
	if err != nil {
		return err
	}

	// run measurement for FMAIS of depth 2
	err = runMeasurementForDepth2(false)
	if err != nil {
		return err
	}

	// re-assigning a deviation in order to generate smoother results in large-scale measurement
	deviation = 0.1 * deviation
	// run measurement for FMAIS of depth 4 with large number of applications
	err = runMeasurementWide(maxAppNumWide, stepWide, false)
	if err != nil {
		return err
	}

	return nil
}

// generateRandomVectorOfLength function generates a random vector of given length with random values with Normal distribution
// and a mean value meanVal
func generateRandomVectorOfLength(meanVal float64, length int) map[int]float64 {
	arr := make(map[int]float64, length)

	for i := 0; i < length; i++ {
		// an elegant solution
		rnd := (rand.Float64()*2-1)*deviation + meanVal
		arr[i] = rnd
	}

	return arr
}

// generateInputDataForInstance generates input data for Application #1
func generateInputDataForInstance(id []inputData) map[int]float64 {
	arr := make(map[int]float64, 0)

	for _, inData := range id {
		newMap := generateRandomVectorOfLength(inData.value, inData.to-inData.from+1)
		// updating map
		for k, v := range newMap {
			arr[inData.from+k] = v
		}
	}

	return arr
}

// initializeInputDataDepth4 function initializes input data for the measurement with Depth 4
func initializeInputDataDepth4() (app1, app2, app3, app4, vi) {

	app1 := app1{
		inst1: generateInputDataForInstance(app1inst1),
		inst2: generateInputDataForInstance(app1inst2),
		inst3: generateInputDataForInstance(app1inst3),
	}

	app2 := app2{
		inst1: generateInputDataForInstance(app2inst1),
		inst2: generateInputDataForInstance(app2inst2),
	}

	app3 := app3{
		inst1: generateInputDataForInstance(app3inst1),
		inst2: generateInputDataForInstance(app3inst2),
	}

	app4 := app4{
		inst1: generateInputDataForInstance(app4inst1),
	}

	viaas := vi{
		inst: generateInputDataForInstance(viaas),
	}

	return app1, app2, app3, app4, viaas
}

// initializeInputDataDepth3 function initializes input data for the measurement with Depth 3
func initializeInputDataDepth3() (app1, app2, app3, vi) {

	app1 := app1{
		inst1: generateInputDataForInstance(app1inst1),
		inst2: generateInputDataForInstance(app1inst2),
		inst3: generateInputDataForInstance(app1inst3),
	}

	app2 := app2{
		inst1: generateInputDataForInstance(app2inst1),
		inst2: generateInputDataForInstance(app2inst2),
	}

	app3 := app3{
		inst1: generateInputDataForInstance(app3inst1),
		inst2: generateInputDataForInstance(app3inst2),
	}

	viaas := vi{
		inst: generateInputDataForInstance(viaas),
	}

	return app1, app2, app3, viaas
}

// initializeInputDataDepth2 function initializes input data for the measurement with Depth 2
func initializeInputDataDepth2() (app1, app2, vi) {

	app1 := app1{
		inst1: generateInputDataForInstance(app1inst1),
		inst2: generateInputDataForInstance(app1inst2),
		inst3: generateInputDataForInstance(app1inst3),
	}

	app2 := app2{
		inst1: generateInputDataForInstance(app2inst1),
		inst2: generateInputDataForInstance(app2inst2),
	}

	viaas := vi{
		inst: generateInputDataForInstance(viaas),
	}

	return app1, app2, viaas
}

// initializeInputDataWide function initializes input data for the large-scale measurement with Depth 4
func initializeInputDataWide() (appNoFail, appFail) {

	app := appNoFail{
		inst1: generateInputDataForInstance(appInst1),
		inst2: generateInputDataForInstance(appInst2),
	}

	appF := appFail{
		inst1: generateInputDataForInstance(appFailInst1),
		inst2: generateInputDataForInstance(appFailInst2),
	}

	return app, appF
}

func updateReliabilities(sm *systemmodel.SystemModel, step int, i ...interface{}) error {

	for _, item := range i {
		switch t := item.(type) {
		case app1:
			err := sm.UpdateApplicationReliability(app1Name, map[int64]float64{
				1: t.inst1[step],
				2: t.inst2[step],
				3: t.inst3[step],
			})
			if err != nil {
				return fmt.Errorf("something went wrong during update of Application #1 instances reliabilities: %w", err)
			}
		case app2:
			err := sm.UpdateApplicationReliability(app2Name, map[int64]float64{
				1: t.inst1[step],
				2: t.inst2[step],
			})
			if err != nil {
				return fmt.Errorf("something went wrong during update of Application #2 instances reliabilities: %w", err)
			}
		case app3:
			err := sm.UpdateApplicationReliability(app3Name, map[int64]float64{
				1: t.inst1[step],
				2: t.inst2[step],
			})
			if err != nil {
				return fmt.Errorf("something went wrong during update of Application #3 instances reliabilities: %w", err)
			}
		case app4:
			err := sm.UpdateApplicationReliability(app4Name, map[int64]float64{
				1: t.inst1[step],
			})
			if err != nil {
				return fmt.Errorf("something went wrong during update of Application #4 instance reliabilities: %w", err)
			}
		case vi:
			err := sm.UpdateApplicationReliability(viName, map[int64]float64{
				1: t.inst[step],
			})
			if err != nil {
				return fmt.Errorf("something went wrong during update of VI instance reliabilities: %w", err)
			}
		case appFail:
			err := sm.UpdateApplicationReliability(appFailName, map[int64]float64{
				1: t.inst1[step],
				2: t.inst2[step],
			})
			if err != nil {
				return fmt.Errorf("something went wrong during update of %s (Failed App) instance reliabilities: %w", appFailName, err)
			}
		case appNoFail:
			appName := "App#"
			for a := 2; a <= len(sm.Applications)-1; a++ {
				updVal := appNoFail{
					inst1: generateInputDataForInstance(appInst1),
					inst2: generateInputDataForInstance(appInst2),
				}
				err := sm.UpdateApplicationReliability(appName+strconv.Itoa(a), map[int64]float64{
					1: updVal.inst1[step],
					2: updVal.inst2[step],
				})
				if err != nil {
					return fmt.Errorf("something went wrong during update of %s instance reliabilities: %w", appName+strconv.Itoa(a), err)
				}
			}
		default:
			return fmt.Errorf("received an unexpected type: %v", t)
		}
	}
	return nil
}

// runMeasurementForDepth4 function runs a measurement for FMAIS of Depth 4
func runMeasurementForDepth4(test bool) error {
	log.Printf("Running measurement for FMAIS of depth 4\n")
	// initializing a reliability map
	relArr := make(map[int]float64, 0)

	// initializing input data
	app1, app2, app3, app4, vi := initializeInputDataDepth4()

	// initialising system model
	sm4 := systemmodel.CreateSystemModelDepth4()

	meErtCore := meertcore.MeErtCore{
		SystemModel: sm4,
		Reliability: 0.0,
	}

	// running the measurement itself
	for i := 1; i <= 300; i++ {
		// setting reliabilities for each instance
		err := updateReliabilities(meErtCore.SystemModel, i, app1, app2, app3, app4, vi)
		if err != nil {
			return fmt.Errorf("something went wrong during updating of Application/VI reliabilities: %w", err)
		}

		_, err = meErtCore.SystemModel.GatherAllApplicationsReliabilities()
		if err != nil {
			return fmt.Errorf("something went wrong during gathering of all application reliabilities: %w", err)
		}

		// computing reliability of the FMAIS per (optimized) ME-ERT-CORE
		rel, err := meErtCore.ComputeReliabilityOptimizedSimple()
		if err != nil {
			return fmt.Errorf("something went wrong during the reliability computation (per optimized method): %w", err)
		}

		// updating reliability map
		relArr[i] = rel
	}

	if !test {
		// exporting data to JSON
		err := storedata.ExportDataToJSON("data/", "me-ert-core_fmais_depth_"+strconv.Itoa(sm4.Depth),
			relArr, "", " ")
		if err != nil {
			log.Panicf("Something went wrong during storing of the data in JSON file... %v\n", err)
			return err
		}
		// plotting a graph for measured reliability
		err = draw.PlotMeasuredReliability(relArr, len(sm4.Applications)-1, sm4.Depth, false, false)
		if err != nil {
			return fmt.Errorf("something went wrong during plotting of a reliability values: %w", err)
		}

		// compute ME-ERT-CORE coefficients and plot it
		coefs, err := computeMeErtCoreCoefficients(relArr, len(sm4.Applications)-1)
		if err != nil {
			return err
		}

		// plotting a graph for measured reliability
		err = draw.PlotMeErtCoreCoefficients(coefs, len(sm4.Applications)-1, sm4.Depth, false, false)
		if err != nil {
			return fmt.Errorf("something went wrong during plotting of a ME-ERT-CORE coefficients: %w", err)
		}
	}

	return nil
}

// runMeasurementForDepth3 function runs a measurement for FMAIS of Depth 3
func runMeasurementForDepth3(test bool) error {
	log.Printf("Running measurement for FMAIS of depth 3\n")

	// initializing a reliability map
	relArr := make(map[int]float64, 0)

	// initializing input data
	app1, app2, app3, vi := initializeInputDataDepth3()

	// initialising system model
	sm3 := systemmodel.CreateSystemModelDepth3()

	meErtCore := meertcore.MeErtCore{
		SystemModel: sm3,
		Reliability: 0.0,
	}

	// running the measurement itself
	for i := 1; i <= 300; i++ {
		// setting reliabilities for each instance
		err := updateReliabilities(meErtCore.SystemModel, i, app1, app2, app3, vi)
		if err != nil {
			return fmt.Errorf("something went wrong during updating of Application/VI reliabilities: %w", err)
		}

		_, err = meErtCore.SystemModel.GatherAllApplicationsReliabilities()
		if err != nil {
			return fmt.Errorf("something went wrong during gathering of all application reliabilities: %w", err)
		}

		// computing reliability of the FMAIS per (optimized) ME-ERT-CORE
		rel, err := meErtCore.ComputeReliabilityOptimizedSimple()
		if err != nil {
			return fmt.Errorf("something went wrong during the reliability computation (per optimized method): %w", err)
		}

		// updating reliability map
		relArr[i] = rel
	}

	if !test {
		// exporting data to JSON
		err := storedata.ExportDataToJSON("data/", "me-ert-core_fmais_depth_"+strconv.Itoa(sm3.Depth),
			relArr, "", " ")
		if err != nil {
			log.Panicf("Something went wrong during storing of the data in JSON file... %v\n", err)
			return err
		}
		// plotting a graph for measured reliability
		err = draw.PlotMeasuredReliability(relArr, len(sm3.Applications)-1, sm3.Depth, false, false)
		if err != nil {
			return fmt.Errorf("something went wrong during plotting of a reliability values: %w", err)
		}

		// compute ME-ERT-CORE coefficients and plot it
		coefs, err := computeMeErtCoreCoefficients(relArr, len(sm3.Applications)-1)
		if err != nil {
			return err
		}

		// plotting a graph for measured reliability
		err = draw.PlotMeErtCoreCoefficients(coefs, len(sm3.Applications)-1, sm3.Depth, false, false)
		if err != nil {
			return fmt.Errorf("something went wrong during plotting of a ME-ERT-CORE coefficients: %w", err)
		}
	}

	return nil
}

// runMeasurementForDepth2 function runs a measurement for FMAIS of Depth 2
func runMeasurementForDepth2(test bool) error {
	log.Printf("Running measurement for FMAIS of depth 2\n")

	// initializing a reliability map
	relArr := make(map[int]float64, 0)

	// initializing input data
	app1, app2, vi := initializeInputDataDepth2()

	// initialising system model
	sm2 := systemmodel.CreateSystemModelDepth2()

	meErtCore := meertcore.MeErtCore{
		SystemModel: sm2,
		Reliability: 0.0,
	}

	// running the measurement itself
	for i := 1; i <= 300; i++ {
		// setting reliabilities for each instance
		err := updateReliabilities(meErtCore.SystemModel, i, app1, app2, vi)
		if err != nil {
			return fmt.Errorf("something went wrong during updating of Application/VI reliabilities: %w", err)
		}

		_, err = meErtCore.SystemModel.GatherAllApplicationsReliabilities()
		if err != nil {
			return fmt.Errorf("something went wrong during gathering of all application reliabilities: %w", err)
		}

		// computing reliability of the FMAIS per (optimized) ME-ERT-CORE
		rel, err := meErtCore.ComputeReliabilityOptimizedSimple()
		if err != nil {
			return fmt.Errorf("something went wrong during the reliability computation (per optimized method): %w", err)
		}

		// updating reliability map
		relArr[i] = rel
	}

	if !test {
		// exporting data to JSON
		err := storedata.ExportDataToJSON("data/", "me-ert-core_fmais_depth_"+strconv.Itoa(sm2.Depth),
			relArr, "", " ")
		if err != nil {
			log.Panicf("Something went wrong during storing of the data in JSON file... %v\n", err)
			return err
		}
		// plotting a graph for measured reliability
		err = draw.PlotMeasuredReliability(relArr, len(sm2.Applications)-1, sm2.Depth, false, false)
		if err != nil {
			return fmt.Errorf("something went wrong during plotting of a reliability values: %w", err)
		}

		// compute ME-ERT-CORE coefficients and plot it
		coefs, err := computeMeErtCoreCoefficients(relArr, len(sm2.Applications)-1)
		if err != nil {
			return err
		}

		// plotting a graph for measured reliability
		err = draw.PlotMeErtCoreCoefficients(coefs, len(sm2.Applications)-1, sm2.Depth, false, false)
		if err != nil {
			return fmt.Errorf("something went wrong during plotting of a ME-ERT-CORE coefficients: %w", err)
		}
	}

	return nil
}

func computeMeErtCoreCoefficients(relMap map[int]float64, appNum int) (map[int]float64, error) {

	if len(relMap) != 300 {
		return nil, fmt.Errorf("obtained incomplete map with %d elements in it: %v", len(relMap), relMap)
	}

	res := make(map[int]float64, 0)
	for i := 1; i <= 300; i++ {
		rel, ok := relMap[i]
		if !ok {
			return nil, fmt.Errorf("map entry for key %d does NOT exist, relMap is: %v", i, relMap)
		}
		coef, err := meertcore.ComputeMeErtCoreCoefficient(rel, appNum)
		if err != nil {
			return nil, err
		}
		res[i] = coef
	}

	return res, nil
}

// runMeasurementWide function runs a measurement for a large-scale FMAIS, which contains up to 100 Applications,
// which are identical in terms of component priority distributions. This measurement is intended to showcase the
// application of a ME-ERT-CORE coefficient
func runMeasurementWide(maxAppNum, step int, test bool) error {
	log.Printf("Running measurement for large-scale FMAIS\n")

	// initializing an overall reliability map
	relArr := make(map[int]map[int]float64, 0)

	// initializing an overall ME-ERT-CORE coefficient map
	meErtCoreCoefs := make(map[int]map[int]float64, 0)

	// initializing input data
	app, appFailed := initializeInputDataWide()

	for a := 10; a <= maxAppNum; a = a + step {
		log.Printf("Running large-scale measurement for FMAIS with %d Apps\n", a)

		// initializing a reliability map for a current experiment
		relMeasured := make(map[int]float64, 0)

		// initialize FMAIS of depth 4 with up to 100 applications (5 applications per VI)
		sm, err := systemmodel.CreateSystemModelWide(a)
		if err != nil {
			return err
		}

		meErtCore := meertcore.MeErtCore{
			SystemModel: sm,
			Reliability: 0.0,
		}

		// running the measurement itself
		for i := 1; i <= 300; i++ {
			// setting reliabilities for each instance
			err = updateReliabilities(meErtCore.SystemModel, i, appFailed, app)
			if err != nil {
				return fmt.Errorf("something went wrong during updating of Application/VI reliabilities: %w", err)
			}

			_, err = meErtCore.SystemModel.GatherAllApplicationsReliabilities()
			if err != nil {
				return fmt.Errorf("something went wrong during gathering of all application reliabilities: %w", err)
			}

			// computing reliability of the FMAIS per (optimized) ME-ERT-CORE
			rel, err := meErtCore.ComputeReliabilityOptimizedSimple()
			if err != nil {
				return fmt.Errorf("something went wrong during the reliability computation (per optimized method): %w", err)
			}

			// updating reliability map
			relMeasured[i] = rel
		}

		// updating reliability map
		relArr[a] = relMeasured

		// computing ME-ERT-CORE coefficients
		meErtCoreCoefs[a], err = computeMeErtCoreCoefficients(relMeasured, len(sm.Applications)-1)
		if err != nil {
			return fmt.Errorf("something went wrong while computing ME-ERT-CORE coefficients for %d Apps with "+
				"following input data: %v\n%w", a, relMeasured, err)
		}

		// plot figure
		if !test {
			// plotting a graph for measured reliability
			err = draw.PlotMeasuredReliability(relMeasured, len(sm.Applications)-1, sm.Depth, false, true)
			if err != nil {
				return fmt.Errorf("something went wrong during plotting of a reliability values: %w", err)
			}

			// plotting a graph for measured reliability
			err = draw.PlotMeErtCoreCoefficients(meErtCoreCoefs[a], len(sm.Applications)-1, sm.Depth, false, true)
			if err != nil {
				return fmt.Errorf("something went wrong during plotting of a ME-ERT-CORE coefficients: %w", err)
			}
		}
	}
	log.Printf("Measurement for large-scale FMAIS has finished. Storing results...\n")

	if !test {
		// exporting data to JSON
		err := storedata.ExportDataToJSON("data/", "me-ert-core-wide_fmais_depth_"+strconv.Itoa(4),
			relArr, "", " ")
		if err != nil {
			log.Panicf("Something went wrong during storing of the data in JSON file... %v\n", err)
			return err
		}

		err = storedata.ExportDataToJSON("data/", "me-ert-core-wide-coefs_fmais_depth_"+strconv.Itoa(4),
			meErtCoreCoefs, "", " ")
		if err != nil {
			log.Panicf("Something went wrong during storing of the data in JSON file... %v\n", err)
			return err
		}
	}
	log.Printf("Results are stored. Measurement is finished.\n")

	return nil
}
