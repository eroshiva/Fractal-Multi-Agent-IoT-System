// Package benchmarking implements a benchmarking logic for two test cases - System Model time complexity evaluation
// and Reliability model (ME-ERT-CORE) time complexity evaluation
package benchmarking

import (
	"fmt"
	"gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/pkg/draw"
	"gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/pkg/storedata"
	"gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/pkg/systemmodel"
	"log"
	"runtime"
	"time"
)

var timer uint64

// This map stores time needed to generate a System Model. Notation is:
// map[key1]map[key3]map[key3]time
// key1 is a system model depth
// key2 is a number of the applications within a system
// key3 is a maximum number of instances which one application can deploy
var benchmarkedData map[int]map[int]map[int]float64

// BenchSystemModelNoParam function performs benchmarking of a Fractal MAS System Model and does not require input parameters
func BenchSystemModelNoParam() error {
	// Setting number of iterations to perform on a single parameter set.
	// The more the number is, the less is an error due to system resources fluctuation
	numIterations := 25000
	// setting maximum System Model depth
	maxDepth := 4
	// setting maximum number of applications in MAIS
	maxAppNumber := 100
	// setting a maximum number of instances per Application
	maxNumInstancesPerApp := 100

	err := BenchSystemModel(maxDepth, maxAppNumber, maxNumInstancesPerApp, numIterations)
	if err != nil {
		return err
	}
	return nil
}

// BenchSystemModel function performs benchmarking of a Fractal MAS System Model
func BenchSystemModel(maxDepth int, maxAppNumber int, maxNumInstancesPerApp int, numIterations int) error {
	// initializing some variables to gather statistics
	var maxNumIncs int64 = -1
	var appMax, depthMax, instMax int64
	var allocBytes uint64
	var allocAppMax, allocDepthMax, allocInstMax, allocTotalInstances int64
	// initializing map
	benchmarkedData = make(map[int]map[int]map[int]float64, 0)
	// iterating over the depth of the system
	for depth := 1; depth <= maxDepth; depth++ {
		benchmarkedData[depth] = make(map[int]map[int]float64, 0)
		// iterating over the amount of apps in the system
		for appNumber := 1; appNumber <= maxAppNumber+1; appNumber += 5 {
			benchmarkedData[depth][appNumber] = make(map[int]float64, 0)
			// iterating over the range of the minimum and maximum number of instances deployed by application
			for maxNumInstances := 1; maxNumInstances <= maxNumInstancesPerApp+1; maxNumInstances += 5 {
				log.Printf("%d iterations over Depth %v, App number %v, Number of instances %v\n", numIterations, depth, appNumber, maxNumInstances)
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
					sm.GenerateSystemModel()
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
						appMax = int64(appNumber)
						depthMax = int64(depth)
						instMax = int64(maxNumInstances)
					}
				}
				log.Printf("Benchmarked time is %v us, measured %v us in %d operations\n", float64(timer)/float64(numIterations), timer, numIterations)
				benchmarkedData[depth][appNumber][maxNumInstances] = float64(timer) / float64(numIterations)
			}
		}
	}
	log.Printf("Maximum number of instances is %v. It was for depth %v, number applications %v, instances per app %v.\n",
		maxNumIncs, depthMax, appMax, instMax)
	log.Printf("Maximum amount of allocated memory is %v MB. It was for depth %v, number applications %v, instances per app %v."+
		" Number of instances at this point was %v.\n",
		allocBytes, allocDepthMax, allocAppMax, allocInstMax, allocTotalInstances)

	// get current time to format a filename
	ct := time.Now()
	err := storedata.SaveData(benchmarkedData, "benchmark_"+ct.Format(time.DateOnly)+"_"+ct.Format(time.TimeOnly))
	if err != nil {
		log.Panicf("Something went wrong when storing bechmarked data... %v\n", err)
	}

	err = draw.PlotTimeComplexities(benchmarkedData, maxDepth, maxAppNumber, maxNumInstancesPerApp)
	if err != nil {
		log.Panicf("Something went wrong during plotting of the results of benchmarking... %v\n", err)
	}

	return nil
}

// BenchMeErtCORE function performs benchmarking of a ME-ERT-CORE reliability model
func BenchMeErtCORE() error {

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

// PrintMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
