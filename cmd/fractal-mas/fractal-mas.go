// main package structures experiment which measures time complexity of the algorithms
package main

import (
	"fmt"
	"gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/pkg/draw"
	"gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/pkg/systemmodel"
	"strconv"
	"time"
)

func main() {
	start := time.Now()
	fmt.Println("Hello, World!")
	duration := time.Since(start)
	fmt.Printf("It took %d us to print previous message\n", duration.Microseconds())

	// Generating a system Model
	sm := systemmodel.SystemModel{}
	var l int32 = 4
	// defining a maximum number of instances per application
	var maxNumInstances int32 = 15 // 15 instances per app
	// defining list of application names
	names := []string{"VI", "App#1", "App#2", "App#3", "App#4", "App#5", "App#6", "App#7",
		"App#8", "App#9"}
	sm.InitializeSystemModel(maxNumInstances, l)
	sm.CreateRandomApplications(names, maxNumInstances)
	start = time.Now()
	sm.GenerateSystemModel()
	duration = time.Since(start)
	//sm.PrettyPrintLayers()
	fmt.Printf("It took %d us to generate a random System Model\n", duration.Microseconds())

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
	fmt.Printf("It took %d us to draw a System Model Figure\n", duration.Microseconds())
}
