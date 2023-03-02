// Package benchmarking implements a benchmarking logic for two test cases - System Model time complexity evaluation
// and Reliability model (ME-ERT-CORE) time complexity evaluation
package benchmarking

import (
	"gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/pkg/draw"
	"gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/pkg/systemmodel"
	"log"
	"math"
	"time"
)

var timer uint64

// BenchSystemModelNoParam function performs benchmarking of a Fractal MAS System Model and does not require input parameters
func BenchSystemModelNoParam() error {
	// FIXME - change this parameters before actual benchmarking!!
	// Setting number of iterations to perform on a single parameter set.
	// The more the number is, the less is an error due to system resources fluctuation
	numIterations := 100
	// setting maximum System Model depth
	maxDepth := 5
	// setting maximum number of applications in MAIS
	maxAppNumber := 51
	// setting a power of maximum number of instances per Application (e.g., 3 corresponds to 10^3, thus algorithm
	// would iterate over 10^0, 10^1, 10^2 and 10^3)
	maxNumInstancesPerApp := 2

	err := BenchSystemModel(maxDepth, maxAppNumber, maxNumInstancesPerApp, numIterations)
	if err != nil {
		return err
	}
	return nil
}

// BenchSystemModel function performs benchmarking of a Fractal MAS System Model
func BenchSystemModel(maxDepth int, maxAppNumber int, maxNumInstancesPerApp int, numIterations int) error {
	// initializing map
	//benchmarkedData = make(map[common.MapKey]float64, 0)
	benchmarkedData = make(map[int]map[int]map[int]float64, 0)
	// iterating over the depth of the system
	for depth := 1; depth <= maxDepth; depth++ {
		benchmarkedData[depth] = make(map[int]map[int]float64, 0)
		// iterating over the amount of apps in the system
		for appNumber := 1; appNumber < maxAppNumber; appNumber += 5 {
			benchmarkedData[depth][appNumber] = make(map[int]float64, 0)
			// iterating over the range of the minimum and maximum number of instances deployed by application
			for power := 0; power <= maxNumInstancesPerApp; power++ {
				maxNumInstances := int(math.Pow10(power))
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
					// we know that it's going to be positive number
					timer += uint64(duration.Nanoseconds()) // taking nanoseconds for better preciseness
				}
				log.Printf("Benchmarked time is %v ns, measured %v ns in %d operations\n", float64(timer)/float64(numIterations), timer, numIterations)
				benchmarkedData[depth][appNumber][maxNumInstances] = float64(timer) / float64(numIterations)
			}
		}
	}
	//log.Printf("Gathered data are:\n%v\n", benchmarkedData)
	err := draw.PlotTimeComplexities(benchmarkedData, maxDepth, maxAppNumber, maxNumInstancesPerApp)
	if err != nil {
		log.Panicf("Something went wrong during benchmarking... %v\n", err)
	}

	err = exportDataToJSON("data/", "test", benchmarkedData, "", " ")
	if err != nil {
		log.Panicf("Something went wrong during storing of the data in JSON file... %v\n", err)
		return err
	}

	err = exportDataToCSV("data/", "test", benchmarkedData, "Fractal MAS Depth [-]",
		"Application Number in Fractal MAS [-]", "Maximum Number of Instances Deployed by Application [-]",
		"Time [ns]")
	if err != nil {
		log.Panicf("Something went wrong during storing of the data in CSV file... %v\n", err)
		return err
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
