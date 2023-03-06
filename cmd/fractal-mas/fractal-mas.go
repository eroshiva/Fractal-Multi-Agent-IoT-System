// main package structures experiment which measures time complexity of the algorithms
package main

import (
	"github.com/spf13/cobra"
	"gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/pkg/benchmarking"
	"gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/pkg/draw"
	"gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/pkg/systemmodel"
	"log"
	"os"
	"strconv"
	"time"
)

var depth int
var appNumber int
var iterations int
var maxNumInstances int
var benchFiles string

// The main entry point
func main() {
	if err := fractalMAIS().Execute(); err != nil {
		println(err)
		os.Exit(1)
	}
}

// fractalMAIS implements a command line interface for Fractal MAS project
func fractalMAIS() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fractal-mas",
		Short: "Fractal Multi-Agent System benchmarker",
		RunE:  runFractalMAIS,
	}
	// adding flags - default values are false
	cmd.PersistentFlags().Bool("example", false, "generates in a single run a random Fractal MAS and plots a figure of it")
	cmd.PersistentFlags().Bool("benchmark", false, "performs a time complexity benchmarking of a Fractal MAS system model algorithm and ME-ERT-CORE algorithms")
	cmd.PersistentFlags().Bool("hardcoded", false, "performs a hardcoded benchmarking (with hardcoded values")
	cmd.PersistentFlags().Bool("benchFMAS", false, "performs a time complexity benchmarking of a Fractal MAS system model algorithm")
	// ToDo - this is to be done in the near future..
	cmd.PersistentFlags().Bool("benchMeErtCORE", false, "performs a time complexity benchmarking of a ME-ERT-CORE algorithms")
	// ToDo - this is to be done in the near future..
	//cmd.PersistentFlags().Bool("benchErtCORE", false, "performs a time complexity benchmarking of a ERT-CORE algorithms")
	cmd.PersistentFlags().IntVar(&iterations, "iterations", 100, "sets a number of iterations per single parameter set to perform")
	cmd.PersistentFlags().IntVar(&depth, "depth", 4, "sets a depth of a system model")
	cmd.PersistentFlags().IntVar(&appNumber, "appNumber", 100, "number of applications to be deployed")
	cmd.PersistentFlags().IntVar(&maxNumInstances, "maxNumInstances", 100, "maximum number of instances to be deployed by application")
	cmd.PersistentFlags().StringVar(&benchFiles, "generateFigures", "benchmarked_*.json", "generates figures based on the provided benchmarked data")
	return cmd
}

// runFractalMAIS performs main logic of the CLI - either showcases an example System Model or performs benchmarking
func runFractalMAIS(cmd *cobra.Command, args []string) error {
	example, _ := cmd.Flags().GetBool("example")
	benchmark, _ := cmd.Flags().GetBool("benchmark")
	hardcoded, _ := cmd.Flags().GetBool("hardcoded")
	benchFMAS, _ := cmd.Flags().GetBool("benchFMAS")
	benchMeErtCORE, _ := cmd.Flags().GetBool("benchMeErtCORE")
	//benchErtCORE, _ := cmd.Flags().GetBool("benchErtCORE")
	// ToDo - do I need to read this flag or it is automatically read from the CLI??? I guess, the latter
	iterations, _ = cmd.Flags().GetInt("iterations")
	benchFiles, _ = cmd.Flags().GetString("generateFigures")

	log.Printf("Starting fractal-mas\nExample: %v\nBenchmarking: %v\n"+
		"Hardcoded: %v\nBenchmark Fractal MAS: %v\nBenchmark ME-ERT-CORE: %v\n"+
		"Depth: %v\nNumber of applications: %v\nMaximum number of instances per application: %v\n"+
		"Data file provided: %v\n",
		example, benchmark, hardcoded, benchFMAS, benchMeErtCORE,
		depth, appNumber, maxNumInstances, benchFiles)

	if example {
		generateExampleSystemModel()
	}
	if benchmark {
		// ToDo - parse SystemModel flags first..
		err := benchmarking.BenchSystemModel(depth, appNumber, maxNumInstances, iterations)
		if err != nil {
			return err
		}
		err = benchmarking.BenchMeErtCORE()
		if err != nil {
			return err
		}
	}
	if benchFMAS && hardcoded {
		err := benchmarking.BenchSystemModelNoParam()
		if err != nil {
			return err
		}
	}
	if benchFMAS && !hardcoded {
		err := benchmarking.BenchSystemModel(depth, appNumber, maxNumInstances, iterations)
		if err != nil {
			return err
		}
	}
	if benchMeErtCORE {
		err := benchmarking.BenchMeErtCORE()
		if err != nil {
			return err
		}
	}
	//if benchErtCORE {
	//	err := benchmarking.BenchErtCORE()
	//	if err != nil {
	//		return err
	//	}
	//}

	if benchFiles != "" {
		err := draw.PlotFigures(benchFiles)
		if err != nil {
			return err
		}
	}

	return nil
}

// generateExampleSystemModel generates System Model example
func generateExampleSystemModel() {
	// Generating a system Model
	sm := systemmodel.SystemModel{}
	// defining list of application names
	names := systemmodel.GenerateAppNames(appNumber)
	sm.InitializeSystemModel(appNumber, depth)
	sm.CreateRandomApplications(names, 1, maxNumInstances)
	start := time.Now()
	sm.GenerateSystemModel()
	duration := time.Since(start)
	log.Printf("It took %d us to generate a random System Model\n", duration.Microseconds())

	// Drawing a figure of System Model
	d := draw.Draw{}
	d.InitializeDrawStruct()
	d.FigureName = "Random System Model with " + strconv.FormatInt(sm.GetTotalNumberOfInstances(), 10) + " instances"
	start = time.Now()
	err := d.DrawSystemModel(&sm)
	if err != nil {
		panic(err)
	}
	duration = time.Since(start)
	log.Printf("It took %d us to draw a System Model Figure\n", duration.Microseconds())
}
