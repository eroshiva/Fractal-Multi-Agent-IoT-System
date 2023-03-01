// Package benchmarking implements a benchmarking logic for two test cases - System Model time complexity evaluation
// and Reliability model (ME-ERT-CORE) time complexity evaluation
package benchmarking

import (
	"gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/pkg/common"
	"gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/pkg/draw"
	"gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/pkg/systemmodel"
	"log"
	"math"
	"time"
)

var timer uint64

// This map stores time needed to generate a System Model. Notation is:
// map[key1]map[key3]map[key3]time
// key1 is a system model depth
// key2 is a number of the applications within a system
// key3 is a maximum number of instances which one application can deploy
var benchmarkedData map[common.MapKey]float64

// BenchSystemModelNoParam function performs benchmarking of a Fractal MAS System Model and does not require input parameters
func BenchSystemModelNoParam() {
	// FIXME - change this parameters before actual benchmarking!!
	// Setting number of iterations to perform on a single parameter set.
	// The more the number is, the less is an error due to system resources fluctuation
	numIterations := 1000
	// setting maximum System Model depth
	maxDepth := 5
	// setting maximum number of applications in MAIS
	maxAppNumber := 10
	// setting a power of maximum number of instances per Application (e.g., 3 corresponds to 10^3, thus algorithm
	// would iterate over 10^0, 10^1, 10^2 and 10^3)
	maxNumInstancesPerApp := 1

	BenchSystemModel(maxDepth, maxAppNumber, maxNumInstancesPerApp, numIterations)
}

// BenchSystemModel function performs benchmarking of a Fractal MAS System Model
func BenchSystemModel(maxDepth int, maxAppNumber int, maxNumInstancesPerApp int, numIterations int) {
	// initializing map
	benchmarkedData = make(map[common.MapKey]float64, 0)
	// iterating over the depth of the system
	for depth := 1; depth <= maxDepth; depth++ {
		// iterating over the amount of apps in the system
		for appNumber := 1; appNumber < maxAppNumber; appNumber += 5 {
			// iterating over the range of the minimum and maximum number of instances deployed by application
			for power := 0; power <= maxNumInstancesPerApp; power++ {
				maxNumInstances := int(math.Pow10(power))
				//log.Printf("10'000 iterations over Depth %v, App number %v, Number of instances %v\n", depth, appNumber, maxNumInstances)
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
				//log.Printf("Benchmarked time is %v ns, measured %v ns in 10'000 operations\n", float64(timer)/float64(numIterations), timer)
				key := common.MapKey{
					Depth:     depth,
					AppNumber: appNumber,
					Instances: maxNumInstances,
				}
				benchmarkedData[key] = float64(timer) / float64(numIterations)
			}
		}
	}
	//log.Printf("Gathered data are:\n%v\n", benchmarkedData)
	err := draw.PlotTimeComplexities(benchmarkedData, maxDepth, maxAppNumber, maxNumInstancesPerApp)
	if err != nil {
		log.Panicf("Something went wrong during benchmarking... %v\n", err)
	}
}

// BenchMeErtCORE function performs benchmarking of a ME-ERT-CORE reliability model
func BenchMeErtCORE() {

}

// BenchErtCore is a placeholder for future implementation of ErtCore reliability model in Go
//func BenchErtCore () {
//
//}
