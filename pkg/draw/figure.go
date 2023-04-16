// Package draw - here are stored mainly functions which plot a figure with regard to certain scenario..
package draw

import (
	"fmt"
	"gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/pkg/storedata"
	"gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/pkg/systemmodel"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"strconv"
	"strings"
)

const refScale = 1.75 * 1.5
const unitLabels = vg.Centimeter
const fontsizeLabels = refScale * unitLabels * 2 / 3
const fontsizeLegend = refScale * unitLabels / 2

// DrawSystemModel draws a figure representing provided SystemModel. It is done in the following way:
//
//	Iterate over layers, over each node in the layer.
//	Once you've reached a node, create its coordinates (X, Y). After that iterate over each related instance (connection) and
//	create a coordinates (X, Y) for each relation. Then, add corresponding lines between the node and its relations.
//	Next, iterate over the rest of the nodes.
//	To avoid duplication of the nodes, create a custom structure, which carries coordinates (X, Y) information, name of the node (which is unique)
//	and the status (were coordinates (X, Y) already created?).
func (d *Draw) DrawSystemModel(sm *systemmodel.SystemModel) error {

	// Adjusting length of an x and Y axis
	maxItemNumber := sm.GetTheGreatestNumberOfInstancesPerLayer()
	d.XLength = 0.25 * vg.Centimeter * (4 / 3) * vg.Length(maxItemNumber)
	if d.XLength < 20*vg.Inch {
		d.XLength = 20 * vg.Inch
	}

	// creating new figure
	p := plot.New()

	// adding basic data to figure
	p.Title.Text = d.FigureName
	p.X.Label.Text = d.XaxisName
	p.Y.Label.Text = d.YaxisName

	// enabling grid
	if d.gridOn {
		p.Add(plotter.NewGrid())
	}

	// converting SystemModel to plotter-friendly structure
	ds := Coordinates{}
	ds.ConvertSystemModelToDrawStruct(sm)

	// adding lines between nodes
	for i := 1; i <= len(sm.Layers); i++ {
		layer := sm.Layers[i]
		for _, v := range layer.Instances {
			// creating a placeholder for a line
			line := make(plotter.XYs, 2)
			// extracting coordinates of the originate node
			line[0].X, line[0].Y = ds.Points[v.Name].Coordinates.XY(0)
			// iterating over relations
			for _, val := range v.Relations {
				// extracting coordinate of the child node
				line[1].X, line[1].Y = ds.Points[val.Name].Coordinates.XY(0)
				// adding a line to the graph
				err := plotutil.AddLines(p, line)
				if err != nil {
					return err
				}
			}
		}
	}

	// Adding labels to figure
	labels, err := plotter.NewLabels(ds.Labels)
	if err != nil {
		return err
	}
	p.Add(labels)

	// adding points to figure
	err = AddScattersSquare(p, ds.Points)
	if err != nil {
		return err
	}
	// Save the plot to a PNG file
	// ToDo - implement a relative path to enable execution out of everywhere in the system..
	if err := p.Save(d.XLength, d.YLength, "figures/"+d.OutputFileName+".png"); err != nil {
		return err
	}
	d.Rendered = true
	return nil
}

// PlotTimeComplexities plots all measured data (i.e., produces various figures)
func PlotTimeComplexities(tc map[int]map[int]map[int]float64, maxDepth int, maxAppNumber int, maxNumInstancesPerApp int, prefix string) error {
	// Firstly, convert the data into simple (X,Y) thing.
	// We want to plot and showcase following dependencies:
	// 1) Time complexity of the System Model based on its depth
	//		- Take into account system with 26, 56 and 96 apps
	//  	- Also, take into account 26, 56 and 96 instances per application
	//  	- This figure will have 9 curves/graphs
	// 2) Time complexity of the System Model based on the number of applications
	//  	- Fix the level of System Model to 2, 3 and 4
	//  	- Fix the number of instances per application to 26, 56 and 96
	//		- This will produce 9 curves
	// 3) Time complexity of the System Model based on the number of instances deployed per application
	//		- Fix the level of System Model to 2, 3 and 4
	//		- Take into account system with 1, 26, 56 and 96 apps
	//		- This will produce 12 curves as well

	// plotting time complexity dependency based on depth
	depthFigure := Draw{}
	depthFigure.InitializeDrawStruct()
	depthFigure.SetOutputFileName(strings.ToLower(prefix) + "_time-complexity-depth-26-inst").SetFigureName(prefix + ": Time Complexity of Fractal MAIS for 26 instances per App").
		SetYaxisName("Time [us]").SetXaxisName("Depth [-]")
	var depthArr []int
	for i := 1; i <= maxDepth; i++ {
		depthArr = append(depthArr, i)
	}
	lines := GetLinesForDepth(tc, depthArr, []int{1, 26, 51, 76, 96}, []int{26})
	err := depthFigure.PlotTimeComplexity(lines)
	if err != nil {
		return err
	}

	depthFigure.SetOutputFileName(strings.ToLower(prefix) + "_time-complexity-depth-56-inst").SetFigureName(prefix + ": Time Complexity of Fractal MAIS for 56 instances per App")
	lines = GetLinesForDepth(tc, depthArr, []int{1, 26, 51, 76, 96}, []int{56})
	err = depthFigure.PlotTimeComplexity(lines)
	if err != nil {
		return err
	}

	depthFigure.SetOutputFileName(strings.ToLower(prefix) + "_time-complexity-depth-96-inst").SetFigureName(prefix + ": Time Complexity of Fractal MAIS for 96 instances per App")
	lines = GetLinesForDepth(tc, depthArr, []int{1, 26, 51, 76, 96}, []int{96})
	err = depthFigure.PlotTimeComplexity(lines)
	if err != nil {
		return err
	}

	//////// plotting dependencies for applications number
	depthFigure.SetOutputFileName(strings.ToLower(prefix) + "_time-complexity-apps-number-26-inst").SetFigureName(prefix + ": Time Complexity of Fractal MAIS for 26 instances per App").
		SetYaxisName("Time [us]").SetXaxisName("Apps number [-]")
	// iterating over the amount of apps in the system
	var appArr []int
	for appNumber := 1; appNumber <= maxAppNumber; appNumber += 5 {
		appArr = append(appArr, appNumber)
	}
	lines = GetLinesForAppNumber(tc, []int{2, 3, 4}, appArr, []int{26})
	err = depthFigure.PlotTimeComplexity(lines)
	if err != nil {
		return err
	}

	depthFigure.SetOutputFileName(strings.ToLower(prefix) + "_time-complexity-apps-number-56-inst").SetFigureName(prefix + ": Time Complexity of Fractal MAIS for 56 instances per App")
	lines = GetLinesForAppNumber(tc, []int{2, 3, 4}, appArr, []int{56})
	err = depthFigure.PlotTimeComplexity(lines)
	if err != nil {
		return err
	}

	depthFigure.SetOutputFileName(strings.ToLower(prefix) + "_time-complexity-apps-number-96-inst").SetFigureName(prefix + ": Time Complexity of Fractal MAIS for 96 instances per App")
	lines = GetLinesForAppNumber(tc, []int{2, 3, 4}, appArr, []int{96})
	err = depthFigure.PlotTimeComplexity(lines)
	if err != nil {
		return err
	}

	//////// plotting dependencies for instances per application
	var instArr []int
	for instNumber := 1; instNumber <= maxNumInstancesPerApp; instNumber += 5 {
		instArr = append(instArr, instNumber)
	}
	depthFigure.SetOutputFileName(strings.ToLower(prefix) + "_time-complexity-instances-per-app-26-apps").SetFigureName(prefix + ": Time Complexity of Fractal MAIS for 26 Apps in MAIS").
		SetYaxisName("Time [us]").SetXaxisName("Instances (per App) [-]")
	lines = GetLinesForInstances(tc, []int{2, 3, 4}, []int{26}, instArr)
	err = depthFigure.PlotTimeComplexity(lines)
	if err != nil {
		return err
	}

	depthFigure.SetOutputFileName(strings.ToLower(prefix) + "_time-complexity-instances-per-app-56-apps").SetFigureName(prefix + ": Time Complexity of Fractal MAIS for 56 Apps in MAIS")
	lines = GetLinesForInstances(tc, []int{2, 3, 4}, []int{56}, instArr)
	err = depthFigure.PlotTimeComplexity(lines)
	if err != nil {
		return err
	}

	depthFigure.SetOutputFileName(strings.ToLower(prefix) + "_time-complexity-instances-per-app-96-apps").SetFigureName(prefix + ": Time Complexity of Fractal MAIS for 96 Apps in MAIS")
	lines = GetLinesForInstances(tc, []int{2, 3, 4}, []int{96}, instArr)
	err = depthFigure.PlotTimeComplexity(lines)
	if err != nil {
		return err
	}

	depthFigure.SetOutputFileName(strings.ToLower(prefix) + "_time-complexity-instances-per-app-1-apps").SetFigureName(prefix + ": Time Complexity of Fractal MAIS for 1 App in MAIS")
	lines = GetLinesForInstances(tc, []int{2, 3, 4}, []int{1}, instArr)
	err = depthFigure.PlotTimeComplexity(lines)
	if err != nil {
		return err
	}

	return nil
}

// PlotTimeComplexity produces single figure representing a certain case
func (d *Draw) PlotTimeComplexity(lines map[string]plotter.XYs) error {

	p := plot.New()
	// setting Figure name and its parameters
	p.Title.Text = d.FigureName
	p.Title.TextStyle.Font.Size = fontsizeLabels // set the size of Figure name

	// setting X-Axis name and its parameters
	p.X.Label.Text = d.XaxisName
	p.X.Label.TextStyle.Font.Size = fontsizeLabels // set the size of the X label
	p.X.Tick.Label.Font.Size = fontsizeLegend      // set the size of the X-axis numbers

	// setting Y-axis name and its parameters
	p.Y.Label.Text = d.YaxisName
	p.Y.Label.TextStyle.Font.Size = fontsizeLabels // set the size of the Y label
	p.Y.Tick.Label.Font.Size = fontsizeLegend      // set the size of the Y-axis numbers

	// setting Legend parameters
	p.Legend.YOffs = vg.Inch                      // place a legend a bit up
	p.Legend.TextStyle.Font.Size = fontsizeLegend // setting size of a Legend

	// adding grid to the figure
	p.Add(plotter.NewGrid())

	// adding plotters for gathered lines to the figure
	err := AddScattersAndLines(p, lines)
	if err != nil {
		return err
	}

	// Save the plot to a PNG file
	// ToDo - implement a relative path to enable execution out of everywhere in the system..
	if err := p.Save(d.XLength, d.YLength, "figures/"+d.OutputFileName+".eps"); err != nil {
		return err
	}
	if err := p.Save(d.XLength, d.YLength, "figures/"+d.OutputFileName+".png"); err != nil {
		return err
	}
	d.Rendered = true

	return nil
}

// GetLinesForDepth returns (X,Y) coordinates for all provided data where X-axis is fixed to depth
func GetLinesForDepth(tc map[int]map[int]map[int]float64, depth []int, appNumber []int, instances []int) map[string]plotter.XYs {
	lines := make(map[string]plotter.XYs, 0)

	for _, a := range appNumber {
		for _, i := range instances {
			line := make(plotter.XYs, 0)
			for _, d := range depth {
				xy := plotter.XY{
					X: float64(d),
					Y: tc[d][a][i],
				}
				line = append(line, xy)
			}
			// this is to store graphs legend..
			key := "FMAIS with " + strconv.Itoa(a) + " Apps and " + strconv.Itoa(i) + " instances (per App)"
			lines[key] = line
		}
	}

	return lines
}

// GetLinesForAppNumber returns (X,Y) coordinates for all provided data where X-axis is fixed to Apps number
func GetLinesForAppNumber(tc map[int]map[int]map[int]float64, depth []int, appNumber []int, instances []int) map[string]plotter.XYs {
	lines := make(map[string]plotter.XYs, 0)

	for _, d := range depth {
		for _, i := range instances {
			line := make(plotter.XYs, 0)
			for _, a := range appNumber {
				xy := plotter.XY{
					X: float64(a),
					Y: tc[d][a][i],
				}
				line = append(line, xy)
			}
			// this is to store graphs legend..
			key := "FMAIS of depth " + strconv.Itoa(d) + " and " + strconv.Itoa(i) + " instances (per App)"
			lines[key] = line
		}
	}

	return lines
}

// GetLinesForInstances returns (X,Y) coordinates for all provided data where X-axis is fixed to Instances per App number
func GetLinesForInstances(tc map[int]map[int]map[int]float64, depth []int, appNumber []int, instances []int) map[string]plotter.XYs {
	lines := make(map[string]plotter.XYs, 0)

	for _, d := range depth {
		for _, a := range appNumber {
			line := make(plotter.XYs, 0)
			for _, i := range instances {
				xy := plotter.XY{
					X: float64(i),
					Y: tc[d][a][i],
				}
				line = append(line, xy)
			}
			// this is to store graphs legend..
			key := "FMAIS of depth " + strconv.Itoa(d) + " and " + strconv.Itoa(a) + " Apps"
			lines[key] = line
		}
	}

	return lines
}

// PlotFigures function plots a figures for SystemModel for a provided data in file called fileName
func PlotFigures(fileNames ...string) error {

	for _, fileName := range fileNames {
		// ToDo - make a workaround with relative path..
		// read the data first
		data, err := storedata.ImportData("data/", fileName)
		if err != nil {
			return err
		}
		// parse SystemModel data
		depth, apps, instances, err := systemmodel.GetSystemModelParameters(data)
		if err != nil {
			return err
		}

		// cut out the file extension
		name := strings.ReplaceAll(fileName, ".json", "")
		name = strings.ReplaceAll(name, ".csv", "")
		// plot figures
		err = PlotTimeComplexities(data, depth, apps, instances, fmt.Sprintf("Generated %s", name))
		if err != nil {
			return err
		}
	}

	return nil
}
