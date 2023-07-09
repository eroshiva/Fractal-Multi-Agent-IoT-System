// main package structures experiment which measures time complexity of the algorithms
package main

import (
	"github.com/spf13/cobra"
	"gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/internal/benchmarking"
	"gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/internal/measurement"
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
var genFig []string

// The main entry point
func main() {
	if err := fractalMAIS().Execute(); err != nil {
		println(err)
		os.Exit(1)
	}
}

// fractalMAIS implements a command line interface for Fractal MAIS project
func fractalMAIS() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fractal-mais",
		Short: "Fractal Multi-Agent IoT System benchmarker",
		Long: "Fractal Multi-Agent IoT System, or Fractal MAIS, implements a System Model based on means of Fractal theory" +
			"which covers scalable MAIS systems. This tool implements a benchmarker, drawer (to plot results of a benchmark)" +
			"and ME-ERT-CORE (with its predecessor ERT-CORE) packages.",
		RunE: runFractalMAIS,
	}
	// adding flags - default values are false
	cmd.PersistentFlags().Bool("example", false, "generates in a single run a random Fractal MAIS and plots a figure of it")
	cmd.PersistentFlags().Bool("benchmark", false, "performs a time complexity benchmarking of a Fractal MAIS system model algorithm and ME-ERT-CORE algorithms")
	cmd.PersistentFlags().Bool("hardcoded", false, "performs a hardcoded benchmarking (with hardcoded values")
	cmd.PersistentFlags().Bool("benchFMAIS", false, "performs a time complexity benchmarking of a Fractal MAIS system model algorithm")
	cmd.PersistentFlags().Bool("benchMeErtCORE", false, "performs a time complexity benchmarking of a ME-ERT-CORE algorithm")
	cmd.PersistentFlags().Bool("benchMeErtCOREoptimized", false, "performs a time complexity benchmarking of an optimized version of ME-ERT-CORE algorithm")
	// ToDo - this is to be done in the near future..
	//cmd.PersistentFlags().Bool("benchErtCORE", false, "performs a time complexity benchmarking of a ERT-CORE algorithms")
	cmd.PersistentFlags().IntVar(&iterations, "iterations", 25000, "sets a number of iterations per single parameter set to perform")
	cmd.PersistentFlags().IntVar(&depth, "depth", 4, "sets a depth of a system model")
	cmd.PersistentFlags().IntVar(&appNumber, "appNumber", 100, "number of applications to be deployed")
	cmd.PersistentFlags().IntVar(&maxNumInstances, "maxNumInstances", 100, "maximum number of instances to be deployed by application")
	cmd.PersistentFlags().StringArrayVar(&genFig, "generateFigures", nil, "generates figures based on the provided benchmarked data")
	cmd.PersistentFlags().Bool("docker", false, "indicates that the benchmarking is done in Docker container")
	cmd.PersistentFlags().Bool("greyScale", false, "indicates that the plotter should generate figures in grey scale")
	cmd.PersistentFlags().Bool("runMeasurement", false, "runs measurement for FMAIS of Depth 2, 3 and 4")
	return cmd
}

// runFractalMAIS performs main logic of the CLI - either showcases an example System Model or performs benchmarking
func runFractalMAIS(cmd *cobra.Command, _ []string) error {
	example, _ := cmd.Flags().GetBool("example")
	benchmark, _ := cmd.Flags().GetBool("benchmark")
	hardcoded, _ := cmd.Flags().GetBool("hardcoded")
	benchFMAIS, _ := cmd.Flags().GetBool("benchFMAIS")
	benchMeErtCORE, _ := cmd.Flags().GetBool("benchMeErtCORE")
	benchMeErtCOREoptimized, _ := cmd.Flags().GetBool("benchMeErtCOREoptimized")
	//benchErtCORE, _ := cmd.Flags().GetBool("benchErtCORE")
	iterations, _ = cmd.Flags().GetInt("iterations")
	genFig, _ = cmd.Flags().GetStringArray("generateFigures")
	docker, _ := cmd.Flags().GetBool("docker")
	greyScale, _ := cmd.Flags().GetBool("greyScale")
	runMeasurement, _ := cmd.Flags().GetBool("runMeasurement")

	log.Printf("Starting fractal-mais\nExample: %v\nBenchmarking: %v\n"+
		"Hardcoded: %v\nBenchmark Fractal MAIS: %v\nBenchmark ME-ERT-CORE: %v\n"+
		"Depth: %v\nNumber of applications: %v\nMaximum number of instances per application: %v\n"+
		"Data file(s) provided: %v\nBenchmarked in Docker: %v\n",
		example, benchmark, hardcoded, benchFMAIS, benchMeErtCORE,
		depth, appNumber, maxNumInstances, genFig, docker)

	if example {
		generateExampleSystemModel()
	}
	if benchmark && hardcoded {
		err := benchmarking.BenchSystemModelNoParam(docker, greyScale)
		if err != nil {
			return err
		}
		err = benchmarking.BenchMeErtCORENoParam(docker, greyScale)
		if err != nil {
			return err
		}
	}
	if benchmark && !hardcoded {
		err := benchmarking.BenchSystemModel(depth, appNumber, maxNumInstances, iterations, docker, greyScale)
		if err != nil {
			return err
		}
		err = benchmarking.BenchMeErtCORE(depth, appNumber, maxNumInstances, iterations, docker, greyScale)
		if err != nil {
			return err
		}
	}
	if benchFMAIS && hardcoded {
		err := benchmarking.BenchSystemModelNoParam(docker, greyScale)
		if err != nil {
			return err
		}
	}
	if benchFMAIS && !hardcoded {
		err := benchmarking.BenchSystemModel(depth, appNumber, maxNumInstances, iterations, docker, greyScale)
		if err != nil {
			return err
		}
	}
	if benchMeErtCORE && hardcoded {
		err := benchmarking.BenchMeErtCORENoParam(docker, greyScale)
		if err != nil {
			return err
		}
	}
	if benchMeErtCORE && !hardcoded {
		err := benchmarking.BenchMeErtCORE(depth, appNumber, maxNumInstances, iterations, docker, greyScale)
		if err != nil {
			return err
		}
	}

	if benchMeErtCOREoptimized {
		err := benchmarking.BenchMeErtCoreOptimized(appNumber, maxNumInstances, docker, greyScale)
		if err != nil {
			return err
		}
	}

	//if benchErtCORE {
	//	err := benchmarking.BenchErtCORE(docker)
	//	if err != nil {
	//		return err
	//	}
	//}

	if genFig != nil {
		err := draw.PlotFigures(greyScale, genFig...)
		if err != nil {
			return err
		}
	}

	if runMeasurement {
		log.Printf("Running measurement\n")
		err := measurement.RunMeasurement()
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
