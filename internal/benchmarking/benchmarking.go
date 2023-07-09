// Package benchmarking implements a benchmarking logic for two test cases - System Model time complexity evaluation
// and Reliability model (ME-ERT-CORE) time complexity evaluation
package benchmarking

import (
	"fmt"
	"gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/internal/measurement"
	"gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/pkg/draw"
	"gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/pkg/meertcore"
	"gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/pkg/storedata"
	"gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/pkg/systemmodel"
	"log"
	"math/rand"
	"runtime"
	"strings"
	"time"
)

var timer uint64

// Setting number of iterations to perform on a single parameter set.
// The more the number is, the less is an error due to system resources fluctuation
var numIterations = 25000

// setting maximum System Model depth
var maxDepth = 4

// setting maximum number of applications in MAIS
var maxAppNumber = 100

// setting a maximum number of instances per Application
var maxNumInstancesPerApp = 100

// This map stores time needed to generate a System Model. Notation is:
// map[key1]map[key3]map[key3]time
// key1 is a system model depth
// key2 is a number of the applications within a system
// key3 is a maximum number of instances which one application can deploy
var benchmarkedData map[int]map[int]map[int]float64
var benchmarkedAvRel map[int]map[int]map[int]float64

// BenchSystemModelNoParam function performs benchmarking of a Fractal MAIS System Model and does not require input parameters
func BenchSystemModelNoParam(docker, greyScale bool) error {
	err := BenchSystemModel(maxDepth, maxAppNumber, maxNumInstancesPerApp, numIterations, docker, greyScale)
	if err != nil {
		return err
	}
	return nil
}

// BenchSystemModel function performs benchmarking of a Fractal MAIS System Model
func BenchSystemModel(maxDepth int, maxAppNumber int, maxNumInstancesPerApp int, numIterations int, docker, greyScale bool) error {
	// initializing some variables to gather statistics
	var maxNumIncs int64 = -1
	var appMax, depthMax, instMax int
	var allocBytes uint64
	var allocAppMax, allocDepthMax, allocInstMax, allocTotalInstances int64
	// initializing map
	benchmarkedData = make(map[int]map[int]map[int]float64, 0)
	maxNumInst := make(map[int]map[int]map[int]float64, 0)
	// iterating over the depth of the system
	for depth := 1; depth <= maxDepth; depth++ {
		benchmarkedData[depth] = make(map[int]map[int]float64, 0)
		// iterating over the amount of apps in the system
		for appNumber := 1; appNumber <= maxAppNumber+1; appNumber += 5 {
			benchmarkedData[depth][appNumber] = make(map[int]float64, 0)
			// iterating over the range of the minimum and maximum number of instances deployed by application
			for maxNumInstances := 1; maxNumInstances <= maxNumInstancesPerApp+1; maxNumInstances += 5 {
				log.Printf("Fractal MAIS benchmarking: %d iterations over Depth %v, App number %v, Number of instances %v\n", numIterations, depth, appNumber, maxNumInstances)
				// setting timer to 0
				timer = 0
				for iteration := 0; iteration < numIterations; iteration++ {
					// Generating a system Model
					sm := systemmodel.SystemModel{}
					// defining list of application names
					names := systemmodel.GenerateAppNames(appNumber)
					sm.InitializeSystemModel(appNumber, depth)
					sm.CreateRandomApplications(names, 1, maxNumInstances)
					start := time.Now()
					sm.GenerateSystemModel() // generates FMAIS System Model without any parameters (requires additional parsing = some code refactoring, complexity stays the same)
					duration := time.Since(start)
					// we know that it's going to be a positive number
					timer += uint64(duration.Microseconds()) // taking microseconds

					// gather some statistics
					// allocated bytes per this run
					bb := gatherAllocatedBytesSizeInMb()
					if allocBytes < bb {
						allocBytes = bb
						allocAppMax = int64(appNumber)
						allocDepthMax = int64(depth)
						allocInstMax = int64(maxNumInstances)
						allocTotalInstances = sm.GetTotalNumberOfInstances()
					}
					incs := sm.GetTotalNumberOfInstances()
					if maxNumIncs < incs {
						maxNumIncs = incs
						appMax = appNumber
						depthMax = depth
						instMax = maxNumInstances
					}
				}
				log.Printf("Fractal MAIS benchmarking: Benchmarked time is %v us, measured %v us in %d operations\n", float64(timer)/float64(numIterations), timer, numIterations)
				benchmarkedData[depth][appNumber][maxNumInstances] = float64(timer) / float64(numIterations)
			}
		}
	}
	log.Printf("Fractal MAIS benchmarking: Maximum number of instances is %v. It was for depth %v, number applications %v, instances per app %v.\n",
		maxNumIncs, depthMax, appMax, instMax)
	// storing data in a map
	maxNumInst[depthMax] = make(map[int]map[int]float64, 0)
	maxNumInst[depthMax][appMax] = make(map[int]float64, 0)
	maxNumInst[depthMax][appMax][instMax] = float64(maxNumIncs)
	log.Printf("Fractal MAIS benchmarking: Maximum amount of allocated memory is %v MB. It was for depth %v, number applications %v, instances per app %v."+
		" Number of instances at this point was %v.\n",
		allocBytes, allocDepthMax, allocAppMax, allocInstMax, allocTotalInstances)

	// get current time to format a filename
	ct := time.Now()
	// making a string with timestamp
	ts := ct.Format(time.DateOnly) + "_" + ct.Format(time.TimeOnly)
	ts = strings.ReplaceAll(ts, ":", "-")
	if docker {
		ts = "docker_" + ts
	}
	err := storedata.SaveData(benchmarkedData, "benchmark_fmais_"+ts)
	if err != nil {
		log.Panicf("Fractal MAIS benchmarking: Something went wrong when storing bechmarked data... %v\n", err)
	}
	err = storedata.SaveData(maxNumInst, "maxNumInstances_fmais_"+ts)
	if err != nil {
		log.Panicf("Fractal MAIS benchmarking: Something went wrong when storing bechmarked data... %v\n", err)
	}

	prefix := "FMAIS"
	if docker {
		prefix = "Docker_" + prefix
	}
	err = draw.PlotTimeComplexities(benchmarkedData, maxDepth, maxAppNumber, maxNumInstancesPerApp, prefix, greyScale, false)
	if err != nil {
		log.Panicf("Fractal MAIS benchmarking: Something went wrong during plotting of the results of benchmarking... %v\n", err)
	}

	return nil
}

// BenchMeErtCORENoParam function performs benchmarking of a ME-ERT-CORE Reliability Model and does not require input parameters
func BenchMeErtCORENoParam(docker, greyScale bool) error {
	err := BenchMeErtCORE(maxDepth, maxAppNumber, maxNumInstancesPerApp, numIterations, docker, greyScale)
	if err != nil {
		return err
	}
	return nil
}

// BenchMeErtCORE function performs benchmarking of a ME-ERT-CORE reliability model
func BenchMeErtCORE(maxDepth int, maxAppNumber int, maxNumInstancesPerApp int, numIterations int, docker, greyScale bool) error {
	// initializing some variables to gather statistics
	var maxNumIncs int64 = -1
	var appMaxInst, depthMaxInst, instMaxInst int
	var avRel float64
	var maxRel float64 = -1
	var appMax, depthMax, instMax int
	var minRel float64 = 100000000
	var appMin, depthMin, instMin int
	// initializing map
	benchmarkedData = make(map[int]map[int]map[int]float64, 0)
	benchmarkedAvRel = make(map[int]map[int]map[int]float64, 0)
	benchmarkedMaxRel := make(map[int]map[int]map[int]float64, 0)
	benchmarkedMinRel := make(map[int]map[int]map[int]float64, 0)
	maxNumInst := make(map[int]map[int]map[int]float64, 0)
	// iterating over the depth of the system
	for depth := 1; depth <= maxDepth; depth++ {
		benchmarkedData[depth] = make(map[int]map[int]float64, 0)
		benchmarkedAvRel[depth] = make(map[int]map[int]float64, 0)
		// iterating over the amount of apps in the system
		for appNumber := 1; appNumber <= maxAppNumber+1; appNumber += 5 {
			benchmarkedData[depth][appNumber] = make(map[int]float64, 0)
			benchmarkedAvRel[depth][appNumber] = make(map[int]float64, 0)
			// iterating over the range of the minimum and maximum number of instances deployed by application
			for maxNumInstances := 1; maxNumInstances <= maxNumInstancesPerApp+1; maxNumInstances += 5 {
				log.Printf("ME-ERT-CORE benchmarking: %d iterations over Depth %v, App number %v, Number of instances %v\n", numIterations, depth, appNumber, maxNumInstances)
				// setting timer to 0
				timer = 0
				avRel = 0.0
				for iteration := 0; iteration < numIterations; iteration++ {
					// Generating a system Model
					sm := &systemmodel.SystemModel{}
					// defining list of application names
					names := systemmodel.GenerateAppNames(appNumber)
					sm.InitializeSystemModel(appNumber, depth)
					sm.CreateRandomApplications(names, 1, maxNumInstances)
					sm.GenerateSystemModel()
					sm.SetApplicationPrioritiesRandom()

					err := sm.SetInstancePrioritiesRandom()
					if err != nil {
						sm.PrettyPrintApplications().PrettyPrintLayers()
						log.Panicf("Something went wrong when setting instance Priorities: %v\n", err)
					}
					err = sm.SetInstanceReliabilitiesRandom()
					if err != nil {
						sm.PrettyPrintApplications().PrettyPrintLayers()
						log.Panicf("Something went wrong when setting instance Reliabilities: %v\n", err)
					}

					me := meertcore.MeErtCore{
						SystemModel: sm,
						Reliability: -1.23456789,
					}

					// actual measurement
					start := time.Now()
					totalRel, err := me.ComputeReliabilityPerDefinition()
					duration := time.Since(start)
					if err != nil {
						sm.PrettyPrintApplications().PrettyPrintLayers()
						log.Panicf("ME-ERT-CORE benchmarking: an error during reliability computation occurred: %v", err)
					}
					// we know that it's going to be a positive number
					timer += uint64(duration.Microseconds()) // taking microseconds
					avRel += totalRel / float64(numIterations)

					// gathering some statistics
					if maxRel < totalRel {
						maxRel = totalRel
						depthMax = depth
						appMax = appNumber
						instMax = maxNumInstances
					}
					if minRel > totalRel {
						minRel = totalRel
						depthMin = depth
						appMin = appNumber
						instMin = maxNumInstances
					}
					incs := sm.GetTotalNumberOfInstances()
					if maxNumIncs < incs {
						maxNumIncs = incs
						appMaxInst = appNumber
						depthMaxInst = depth
						instMaxInst = maxNumInstances
					}

				}
				log.Printf("ME-ERT-CORE benchmarking: Benchmarked time is %v us, measured %v us in %d operations\n", float64(timer)/float64(numIterations), timer, numIterations)
				benchmarkedData[depth][appNumber][maxNumInstances] = float64(timer) / float64(numIterations)
				benchmarkedAvRel[depth][appNumber][maxNumInstances] = avRel
			}
		}
	}
	log.Printf("Fractal MAIS benchmarking: Maximum number of instances is %v. It was for depth %v, number applications %v, instances per app %v.\n",
		maxNumIncs, depthMaxInst, appMaxInst, instMaxInst)
	// storing data in a map
	maxNumInst[depthMaxInst] = make(map[int]map[int]float64, 0)
	maxNumInst[depthMaxInst][appMaxInst] = make(map[int]float64, 0)
	maxNumInst[depthMaxInst][appMaxInst][instMaxInst] = float64(maxNumIncs)
	log.Printf("ME-ERT-CORE benchmarking: Maximum measured Reliability is %v. It was for depth %v, number applications %v, instances per app %v.\n",
		maxRel, depthMax, appMax, instMax)
	benchmarkedMaxRel[depthMax] = make(map[int]map[int]float64, 0)
	benchmarkedMaxRel[depthMax][appMax] = make(map[int]float64, 0)
	benchmarkedMaxRel[depthMax][appMax][instMax] = maxRel
	log.Printf("ME-ERT-CORE benchmarking: Minimum measured Reliability is %v. It was for depth %v, number applications %v, instances per app %v.\n",
		minRel, depthMin, appMin, instMin)
	benchmarkedMinRel[depthMin] = make(map[int]map[int]float64, 0)
	benchmarkedMinRel[depthMin][appMin] = make(map[int]float64, 0)
	benchmarkedMinRel[depthMin][appMin][instMin] = minRel

	// get current time to format a filename
	ct := time.Now()
	// making a string with timestamp
	ts := ct.Format(time.DateOnly) + "_" + ct.Format(time.TimeOnly)
	ts = strings.ReplaceAll(ts, ":", "-")
	if docker {
		ts = "docker_" + ts
	}
	err := storedata.SaveData(benchmarkedData, "benchmark_meertcore_"+ts)
	if err != nil {
		log.Panicf("ME-ERT-CORE benchmarking: Something went wrong when storing bechmarked data... %v\n", err)
	}
	err = storedata.SaveData(maxNumInst, "maxNumInstances_meertcore_"+ts)
	if err != nil {
		log.Panicf("ME-ERT-CORE benchmarking: Something went wrong when storing maximum number of instnace... %v\n", err)
	}
	err = storedata.SaveData(benchmarkedAvRel, "benchmark_average_reliability_"+ts)
	if err != nil {
		log.Panicf("ME-ERT-CORE benchmarking: Something went wrong when storing bechmarked average reliabilities data... %v\n", err)
	}
	err = storedata.SaveData(benchmarkedMaxRel, "benchmark_maximum_reliability_"+ts)
	if err != nil {
		log.Panicf("ME-ERT-CORE benchmarking: Something went wrong when storing bechmarked maximum reliabilities data... %v\n", err)
	}
	err = storedata.SaveData(benchmarkedMinRel, "benchmark_minimum_reliability_"+ts)
	if err != nil {
		log.Panicf("ME-ERT-CORE benchmarking: Something went wrong when storing bechmarked minimum reliabilities data... %v\n", err)
	}

	prefix := "MeErtCore"
	if docker {
		prefix = "Docker_" + prefix
	}
	err = draw.PlotTimeComplexities(benchmarkedData, maxDepth, maxAppNumber, maxNumInstancesPerApp, prefix, greyScale, true)
	if err != nil {
		log.Panicf("ME-ERT-CORE benchmarking: Something went wrong during plotting of the results of benchmarking... %v\n", err)
	}

	return nil
}

// BenchErtCore is a placeholder for future implementation of ErtCore reliability model in Go
//func BenchErtCore () {
//
//}

func gatherAllocatedBytesSizeInMb() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return m.Alloc / 1024 / 1024
}

// BenchMeErtCoreOptimized function benchmarks optimized version of ME-ERT-CORE
func BenchMeErtCoreOptimized(maxInstNum, appNum int, docker, greyScale bool) error {
	log.Printf("Running benchmarking for large-scale FMAIS with Optimized ME-ERT-CORE\n")

	benchmarkedData = make(map[int]map[int]map[int]float64, 0)

	// initializing input data
	app, appFailed := measurement.InitializeInputDataWide()

	for depth := 2; depth <= maxDepth; depth++ {
		benchmarkedData[depth] = make(map[int]map[int]float64, 0)
		benchmarkedData[depth][appNum] = make(map[int]float64, 0)
		for inst := 1; inst <= maxInstNum; inst += 5 {
			sm, err := systemmodel.CreateSystemModelWideBench(appNum, inst, depth)
			if err != nil {
				return err
			}

			log.Printf("ME-ERT-CORE (optimized) benchmarking: %d iterations over FMAIS of depth %v, with %d Apps and %d instances per App\n", numIterations, sm.Depth, appNum, inst)
			// setting timer to 0
			timer = 0
			for iteration := 0; iteration < numIterations; iteration++ {

				meErtCore := meertcore.MeErtCore{
					SystemModel: sm,
					Reliability: 0.0,
				}

				rnd := rand.Intn(300)
				// setting reliabilities for each instance
				err = measurement.UpdateReliabilities(meErtCore.SystemModel, rnd, appFailed, app)
				if err != nil {
					sm.PrettyPrintApplications().PrettyPrintLayers()
					return fmt.Errorf("something went wrong during updating of Application/VI reliabilities: %w", err)
				}

				_, err = meErtCore.SystemModel.GatherAllApplicationsReliabilities()
				if err != nil {
					sm.PrettyPrintApplications().PrettyPrintLayers()
					return fmt.Errorf("something went wrong during gathering of all application reliabilities: %w", err)
				}

				// computing reliability of the FMAIS per (optimized) ME-ERT-CORE
				start := time.Now()
				_, err = meErtCore.ComputeReliabilityOptimizedSimple()
				duration := time.Since(start)
				if err != nil {
					sm.PrettyPrintApplications().PrettyPrintLayers()
					return fmt.Errorf("something went wrong during the reliability computation (per optimized method): %w", err)
				}
				// we know that it's going to be a positive number
				timer += uint64(duration.Microseconds()) // taking microseconds
			}
			log.Printf("ME-ERT-CORE (optimized) benchmarking: Benchmarked time is %v us, measured %v us in %d operations\n", float64(timer)/float64(numIterations), timer, numIterations)
			benchmarkedData[depth][appNum][inst] = float64(timer) / float64(numIterations)
		}
	}

	// get current time to format a filename
	ct := time.Now()
	// making a string with timestamp
	ts := ct.Format(time.DateOnly) + "_" + ct.Format(time.TimeOnly)
	ts = strings.ReplaceAll(ts, ":", "-")
	if docker {
		ts = "docker_" + ts
	}
	err := storedata.SaveData(benchmarkedData, "benchmark_meertcore_optimized_"+ts)
	if err != nil {
		log.Panicf("ME-ERT-CORE (optimized) benchmarking: Something went wrong when storing bechmarked data... %v\n", err)
	}

	prefix := "MeErtCore_Optimized"
	if docker {
		prefix = "Docker_" + prefix
	}
	err = draw.PlotTimeComplexities(benchmarkedData, maxDepth, maxAppNumber, maxNumInstancesPerApp, prefix, greyScale, true)
	if err != nil {
		log.Panicf("ME-ERT-CORE (optimized) benchmarking: Something went wrong during plotting of the results of benchmarking... %v\n", err)
	}

	return nil
}
