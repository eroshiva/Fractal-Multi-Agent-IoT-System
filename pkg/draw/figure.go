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
const unitLabels = vg.Centimeter * 1.25
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
func PlotTimeComplexities(tc map[int]map[int]map[int]float64, maxDepth int, maxAppNumber int, maxNumInstancesPerApp int, prefix string, greyScale, meertcore bool) error {
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
	//		- Take into account system with 26, 51, 76, and 96 apps
	//		- This will produce 12 curves as well

	figureName := prefix + " Time Complexity\nDependency "

	// plotting time complexity dependency based on depth
	depthFigure := Draw{}
	depthFigure.InitializeDrawStruct()
	depthFigure.SetOutputFileName(strings.ToLower(prefix) + "_time-complexity-depth-26-inst").SetFigureName(figureName + "on the number of layers").
		SetYaxisName("Time [ms]").SetXaxisName("Layers [-]")
	var depthArr []int
	for i := 1; i <= maxDepth; i++ {
		depthArr = append(depthArr, i)
	}
	lines := getLinesForDepth(tc, depthArr, []int{26, 51, 76, 96}, []int{26})
	err := depthFigure.plotTimeComplexity(lines, greyScale, meertcore, false, true, false)
	if err != nil {
		return err
	}

	depthFigure.SetOutputFileName(strings.ToLower(prefix) + "_time-complexity-depth-56-inst").SetFigureName(figureName + "on the number of layers")
	lines = getLinesForDepth(tc, depthArr, []int{26, 51, 76, 96}, []int{56})
	err = depthFigure.plotTimeComplexity(lines, greyScale, meertcore, false, true, false)
	if err != nil {
		return err
	}

	depthFigure.SetOutputFileName(strings.ToLower(prefix) + "_time-complexity-depth-96-inst").SetFigureName(figureName + "on the number of layers")
	lines = getLinesForDepth(tc, depthArr, []int{26, 51, 76, 96}, []int{96})
	err = depthFigure.plotTimeComplexity(lines, greyScale, meertcore, false, true, false)
	if err != nil {
		return err
	}

	//////// plotting dependencies for applications number
	depthFigure.SetOutputFileName(strings.ToLower(prefix) + "_time-complexity-apps-number-26-inst").SetFigureName(figureName + "on the App number").
		SetYaxisName("Time [ms]").SetXaxisName("Number of Applications [-]")
	// iterating over the amount of apps in the system
	var appArr []int
	for appNumber := 1; appNumber <= maxAppNumber; appNumber += 5 {
		appArr = append(appArr, appNumber)
	}
	lines = getLinesForAppNumber(tc, []int{2, 3, 4}, appArr, []int{26})
	err = depthFigure.plotTimeComplexity(lines, greyScale, meertcore, true, false, false)
	if err != nil {
		return err
	}

	depthFigure.SetOutputFileName(strings.ToLower(prefix) + "_time-complexity-apps-number-56-inst").SetFigureName(figureName + "on the App number")
	lines = getLinesForAppNumber(tc, []int{2, 3, 4}, appArr, []int{56})
	err = depthFigure.plotTimeComplexity(lines, greyScale, meertcore, true, false, false)
	if err != nil {
		return err
	}

	depthFigure.SetOutputFileName(strings.ToLower(prefix) + "_time-complexity-apps-number-96-inst").SetFigureName(figureName + "on the App number")
	lines = getLinesForAppNumber(tc, []int{2, 3, 4}, appArr, []int{96})
	err = depthFigure.plotTimeComplexity(lines, greyScale, meertcore, true, false, false)
	if err != nil {
		return err
	}

	//////// plotting dependencies for instances per application
	var instArr []int
	for instNumber := 1; instNumber <= maxNumInstancesPerApp; instNumber += 5 {
		instArr = append(instArr, instNumber)
	}

	depthFigure.SetOutputFileName(strings.ToLower(prefix) + "_time-complexity-instances-per-app-1-apps").SetFigureName(figureName + "on instances per App")
	lines = getLinesForInstances(tc, []int{2, 3, 4}, []int{1}, instArr)
	err = depthFigure.plotTimeComplexity(lines, greyScale, meertcore, false, false, true)
	if err != nil {
		return err
	}

	depthFigure.SetOutputFileName(strings.ToLower(prefix) + "_time-complexity-instances-per-app-26-apps").SetFigureName(figureName + "on instances per App").
		SetYaxisName("Time [ms]").SetXaxisName("Instances (per App) [-]")
	lines = getLinesForInstances(tc, []int{2, 3, 4}, []int{26}, instArr)
	err = depthFigure.plotTimeComplexity(lines, greyScale, meertcore, false, false, true)
	if err != nil {
		return err
	}

	depthFigure.SetOutputFileName(strings.ToLower(prefix) + "_time-complexity-instances-per-app-56-apps").SetFigureName(figureName + "on instances per App")
	lines = getLinesForInstances(tc, []int{2, 3, 4}, []int{56}, instArr)
	err = depthFigure.plotTimeComplexity(lines, greyScale, meertcore, false, false, true)
	if err != nil {
		return err
	}

	depthFigure.SetOutputFileName(strings.ToLower(prefix) + "_time-complexity-instances-per-app-96-apps").SetFigureName(figureName + "on instances per App")
	lines = getLinesForInstances(tc, []int{2, 3, 4}, []int{96}, instArr)
	err = depthFigure.plotTimeComplexity(lines, greyScale, meertcore, false, false, true)
	if err != nil {
		return err
	}

	return nil
}

// plotTimeComplexity produces single figure representing a certain case
func (d *Draw) plotTimeComplexity(lines map[string]plotter.XYs, greyScale, meertcore, appsNumberDep, layerDep, instancesDep bool) error {
	p := d.initializeAndSetPlotter(meertcore, appsNumberDep, layerDep, instancesDep)

	// adding grid to the figure
	if d.gridOn {
		p.Add(plotter.NewGrid())
	}

	// adding plotters for gathered lines to the figure
	err := AddScattersAndLines(p, greyScale, lines)
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

// getLinesForDepth returns (X,Y) coordinates for all provided data where X-axis is fixed to depth
func getLinesForDepth(tc map[int]map[int]map[int]float64, depth []int, appNumber []int, instances []int) map[string]plotter.XYs {
	lines := make(map[string]plotter.XYs, 0)

	for _, a := range appNumber {
		for _, i := range instances {
			line := make(plotter.XYs, 0)
			for _, d := range depth {
				xy := plotter.XY{
					X: float64(d),
					Y: tc[d][a][i] / 1000, // converting to milliseconds
				}
				line = append(line, xy)
			}
			// this is to store graphs legend..
			key := "FMAIS; " + strconv.Itoa(a) + " Apps with " + strconv.Itoa(i) + " instances (per App)"
			lines[key] = line
		}
	}

	return lines
}

// getLinesForAppNumber returns (X,Y) coordinates for all provided data where X-axis is fixed to Apps number
func getLinesForAppNumber(tc map[int]map[int]map[int]float64, depth []int, appNumber []int, instances []int) map[string]plotter.XYs {
	lines := make(map[string]plotter.XYs, 0)

	for _, d := range depth {
		for _, i := range instances {
			line := make(plotter.XYs, 0)
			for _, a := range appNumber {
				xy := plotter.XY{
					X: float64(a),
					Y: tc[d][a][i] / 1000, // converting to milliseconds
				}
				line = append(line, xy)
			}
			// this is to store graphs legend..
			key := "FMAIS; " + strconv.Itoa(d) + " layers with " + strconv.Itoa(i) + " instances (per App)"
			lines[key] = line
		}
	}

	return lines
}

// getLinesForInstances returns (X,Y) coordinates for all provided data where X-axis is fixed to Instances per App number
func getLinesForInstances(tc map[int]map[int]map[int]float64, depth []int, appNumber []int, instances []int) map[string]plotter.XYs {
	lines := make(map[string]plotter.XYs, 0)

	for _, d := range depth {
		for _, a := range appNumber {
			line := make(plotter.XYs, 0)
			for _, i := range instances {
				xy := plotter.XY{
					X: float64(i),
					Y: tc[d][a][i] / 1000, // converting to milliseconds
				}
				line = append(line, xy)
			}
			// this is to store graphs legend..
			key := "FMAIS; " + strconv.Itoa(d) + " layers with " + strconv.Itoa(a) + " Apps"
			lines[key] = line
		}
	}

	return lines
}

// getLinesForReliability converts measure ME-ERT-CORE reliability to plotter-friendly data
func getLinesForReliability(tc map[int]float64, apps, depth int) (map[string]plotter.XYs, error) {
	lines := make(map[string]plotter.XYs, 0)
	line := make(plotter.XYs, 0)

	for i := 1; i <= 300; i++ {
		val, ok := tc[i]
		if !ok {
			return nil, fmt.Errorf("couldn't extract key %d from map %v", i, tc)
		}
		xy := plotter.XY{
			X: float64(i),
			Y: val,
		}
		line = append(line, xy)
	}

	key := fmt.Sprintf("FMAIS; %d layers with %d Apps", depth, apps)
	lines[key] = line

	return lines, nil
}

// PlotFigures function plots a figures for SystemModel for a provided data in file called fileName
func PlotFigures(greyScale bool, fileNames ...string) error {
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
		prefix := fmt.Sprintf("Generated %s", name)
		meertcore := false
		if strings.Contains(strings.ToLower(name), "meertcore") || strings.Contains(strings.ToLower(name), "me-ert-core") {
			meertcore = true
			prefix = "ME-ERT-CORE"
		} else if strings.Contains(strings.ToLower(name), "fmais") {
			prefix = "FMAIS"
		}

		// plot figures
		err = PlotTimeComplexities(data, depth, apps, instances, prefix, greyScale, meertcore)
		if err != nil {
			return err
		}
	}

	return nil
}

// PlotMeasuredReliability plots reliability values obtained from measurement
func PlotMeasuredReliability(rels map[int]float64, apps, depth int, greyScale, wide bool) error {
	// setting figure name
	figureName := "Measured ME-ERT-CORE Reliability values"
	fileName := fmt.Sprintf("measurement_meertcore_depth_%d_apps_%d", depth, apps)
	if wide {
		fileName = fmt.Sprintf("measurement_meertcore_wide_depth_%d_apps_%d", depth, apps)
	}

	// initializing structure for the Figure
	figure := Draw{}
	figure.InitializeDrawStruct().SetFigureName(figureName).SetOutputFileName(fileName).
		SetXaxisName("Time [s]").SetYaxisName("Reliability [-]").
		SetYmin(0).SetYmax(1)

	// converting measured reliability to XY data
	line, err := getLinesForReliability(rels, apps, depth)
	if err != nil {
		return err
	}

	// plotting data
	err = figure.plotMeasuredReliability(line, greyScale)
	if err != nil {
		return err
	}

	return nil
}

// PlotMeasuredReliabilityJoint plots multiple reliability values obtained from measurement
func PlotMeasuredReliabilityJoint(tc map[string]map[int]float64, apps, depth []int, greyScale bool) error {
	linesFMAIS := make(map[string]plotter.XYs, 0)

	// setting figure name
	figureName := "Measured ME-ERT-CORE Reliability values"
	fileName := "measurement_meertcore_joint"

	// initializing structure for the Figure
	figure := Draw{}
	figure.InitializeDrawStruct().SetFigureName(figureName).SetOutputFileName(fileName).
		SetXaxisName("Time [s]").SetYaxisName("Reliability [-]").
		SetYmin(0).SetYmax(1)

	i := 0
	for key, rels := range tc {
		// converting measured reliability to XY data
		line, err := getLinesForReliability(rels, apps[i], depth[i])
		if err != nil {
			return err
		}
		linesFMAIS[key] = line[key]
		i++
	}

	// plotting data
	err := figure.plotMeasuredReliability(linesFMAIS, greyScale)
	if err != nil {
		return err
	}

	return nil
}

// plotMeasuredReliability function plots measured ME-ERT-CORE reliability values
func (d *Draw) plotMeasuredReliability(lines map[string]plotter.XYs, greyScale bool) error {
	p := d.initializeAndSetPlotter(true, false, false, false)

	// adding grid to the figure
	if d.gridOn {
		p.Add(plotter.NewGrid())
	}

	// adding plotters for gathered lines to the figure
	err := AddScattersAndLines(p, greyScale, lines)
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

// initializeAndSetPlotter function initializes and sets default parameters to the figure
func (d *Draw) initializeAndSetPlotter(meertcore, appsNumberDep, layersDep, instancesDep bool) *plot.Plot {
	p := plot.New()

	// setting Figure name and its parameters
	p.Title.Text = d.FigureName
	p.Title.TextStyle.Font.Size = fontsizeLabels // set the size of Figure name

	// setting X-Axis name and its parameters
	p.X.Label.Text = d.XaxisName
	p.X.Label.TextStyle.Font.Size = 1.25 * fontsizeLabels // set the size of the X label
	p.X.Tick.Label.Font.Size = 1.45 * fontsizeLegend      // set the size of the X-axis numbers
	if d.xmin != nil {
		p.X.Min = *d.xmin
	}
	if d.xmax != nil {
		p.X.Max = *d.xmax
	}

	// setting Y-axis name and its parameters
	p.Y.Label.Text = d.YaxisName
	p.Y.Label.TextStyle.Font.Size = 1.25 * fontsizeLabels // set the size of the Y label
	p.Y.Tick.Label.Font.Size = 1.45 * fontsizeLegend      // set the size of the Y-axis numbers
	if d.ymin != nil {
		p.Y.Min = *d.ymin
	}
	if d.ymax != nil {
		p.Y.Max = *d.ymax
	}

	// setting Legend parameters
	p.Legend.TextStyle.Font.Size = 1.25 * fontsizeLegend // setting size of a Legend
	// locating legend on the top left of the figure
	p.Legend.Left = true
	p.Legend.Top = true
	p.Legend.YOffs = -vg.Inch            // place a legend a bit down
	p.Legend.XOffs = 1.5 * vg.Centimeter // place a legend a bit to the right
	if meertcore && appsNumberDep {
		p.Legend.YOffs = -5 * vg.Inch // place a legend a bit down
	}
	if !meertcore && appsNumberDep {
		p.Legend.YOffs = -10 * vg.Inch // place a legend a bit down
		p.Legend.XOffs = 6 * vg.Centimeter
	}
	if meertcore && layersDep {
		p.Legend.XOffs = 0.25 * vg.Centimeter
		p.Legend.YOffs = -0.5 * vg.Inch
	}
	if meertcore && appsNumberDep && layersDep && instancesDep {
		p.Legend.YOffs = -0.5 * vg.Inch
	}
	return p
}

// PlotMeErtCoreCoefficients plots reliability values obtained from measurement
func PlotMeErtCoreCoefficients(rels map[int]float64, apps, depth int, greyScale, wide bool) error {
	// setting figure name
	figureName := "Computed ME-ERT-CORE coefficients"
	fileName := fmt.Sprintf("measurement_meertcore_coef_depth_%d_apps_%d", depth, apps)
	if wide {
		fileName = fmt.Sprintf("measurement_meertcore_wide_coef_depth_%d_apps_%d", depth, apps)
	}

	// initializing structure for the Figure
	figure := Draw{}
	figure.InitializeDrawStruct().SetFigureName(figureName).SetOutputFileName(fileName).
		SetXaxisName("Time [s]").SetYaxisName("ME-ERT-CORE coefficient [-]").
		SetYmin(0).SetYmax(1)

	// converting measured reliability to XY data
	line, err := getLinesForReliability(rels, apps, depth)
	if err != nil {
		return err
	}

	// plotting data
	err = figure.plotMeasuredReliability(line, greyScale)
	if err != nil {
		return err
	}

	return nil
}

// PlotJointFigure function plots a figures for SystemModel for a provided data in file called fileName
func PlotJointFigure(greyScale, meertcore bool, fileNames ...string) error {
	if !meertcore {
		dataArr := make(map[string]map[int]map[int]map[int]float64, len(fileNames))
		var lastPrefix string

		for _, fileName := range fileNames {
			// ToDo - make a workaround with relative path..
			// read the data first
			data, err := storedata.ImportData("data/", fileName)
			if err != nil {
				return err
			}

			// cut out the file extension
			name := strings.ReplaceAll(fileName, ".json", "")
			name = strings.ReplaceAll(name, ".csv", "")
			prefix := fmt.Sprintf("Generated %s", name)
			if strings.Contains(strings.ToLower(name), "optimized") {
				prefix = "(O);"
			} else if strings.Contains(strings.ToLower(name), "definition") {
				prefix = "(PD);"
			}
			dataArr[prefix] = data
			lastPrefix = prefix
		}

		// parse SystemModel data
		_, apps, instances, err := systemmodel.GetSystemModelParameters(dataArr[lastPrefix])
		if err != nil {
			return err
		}

		// plot figures
		err = PlotTimeComplexitiesJoint(dataArr, apps, instances, greyScale)
		if err != nil {
			return err
		}
	} else {
		dataArr := make(map[string]map[int]float64, len(fileNames))
		apps := make([]int, 0)
		apps = append(apps, 2, 4)
		// assuming that passing at input in exact same order like in the Makefile
		i := 0
		for _, fileName := range fileNames {
			// read the data first
			// assuming only JSON files at input
			data, err := storedata.ImportDataMeErtCore("data/", fileName)
			if err != nil {
				return err
			}
			key := fmt.Sprintf("FMAIS; %d layers with %d Apps", apps[i], apps[i])
			dataArr[key] = data
			i++
		}

		err := PlotMeasuredReliabilityJoint(dataArr, apps, apps, greyScale)
		if err != nil {
			return err
		}
	}

	return nil
}

// PlotTimeComplexitiesJoint function plots time complexity curves for all input data files
func PlotTimeComplexitiesJoint(tc map[string]map[int]map[int]map[int]float64, maxAppNumber int,
	maxNumInstancesPerApp int, greyScale bool,
) error {
	// we need to extract and plot time complexity dependency on the number of apps and number of instances per apps
	linesApp := make(map[string]plotter.XYs, 0)
	linesInstApp := make(map[string]plotter.XYs, 0)

	// gathering plotters first
	for key, value := range tc {
		// iterating over the amount of apps in the system
		var appArr []int
		for appNumber := 6; appNumber <= maxAppNumber; appNumber += 5 {
			appArr = append(appArr, appNumber)
		}
		lines := getLinesForAppNumber(value, []int{4}, appArr, []int{26})
		linesApp[key+" FMAIS; 4 layers with 26 instances per App"] = lines["FMAIS; "+strconv.Itoa(4)+" layers with "+strconv.Itoa(26)+" instances (per App)"]

		lines = getLinesForAppNumber(value, []int{4}, appArr, []int{96})
		linesApp[key+" FMAIS; 4 layers with 96 instances per App"] = lines["FMAIS; "+strconv.Itoa(4)+" layers with "+strconv.Itoa(96)+" instances (per App)"]

		//////// plotting dependencies for instances per application
		var instArr []int
		for instNumber := 1; instNumber <= maxNumInstancesPerApp; instNumber += 5 {
			instArr = append(instArr, instNumber)
		}

		lines = getLinesForInstances(value, []int{4}, []int{26}, instArr)
		linesInstApp[key+"\nFMAIS; 4 layers with 26 Apps"] = lines["FMAIS; "+strconv.Itoa(4)+" layers with "+strconv.Itoa(26)+" Apps"]

		lines = getLinesForInstances(value, []int{4}, []int{96}, instArr)
		linesInstApp[key+"\nFMAIS; 4 layers with 96 Apps"] = lines["FMAIS; "+strconv.Itoa(4)+" layers with "+strconv.Itoa(96)+" Apps"]
	}

	// plotting time complexity dependency based on depth
	depthFigure := Draw{}
	depthFigure.InitializeDrawStruct()
	depthFigure.SetOutputFileName("me-ert-core_time-complexity-comparison-apps-number").SetFigureName("ME-ERT-CORE Time Complexity\nDependency on the Apps number").
		SetYaxisName("Time [ms]").SetXaxisName("Number of Applications [-]").SetYmax(12)
	err := depthFigure.plotTimeComplexity(linesApp, greyScale, true, true, true, true) // last 3 false is on purpose: to do not move the legend
	if err != nil {
		return err
	}

	depthFigure1 := Draw{}
	depthFigure1.InitializeDrawStruct()
	depthFigure1.SetOutputFileName("me-ert-core_time-complexity-comparison-app-inst-number").SetFigureName("ME-ERT-CORE Time Complexity\nDependency on the instances per App").
		SetYaxisName("Time [ms]").SetXaxisName("Instances (per App) [-]")
	err = depthFigure1.plotTimeComplexity(linesInstApp, greyScale, false, false, false, false) // last 3 false is on purpose: to do not move the legend
	if err != nil {
		return err
	}

	return nil
}
