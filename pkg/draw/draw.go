// Package draw implements a set of helper functions to draw SystemModel
package draw

import (
	"gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/pkg/systemmodel"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"strings"
)

// Draw structure holds all necessary information for plotting a figure of SystemModel
type Draw struct {
	Rendered       bool   // indicates whether a figure was rendered or not
	OutputFileName string // an output file name, where the figure would be saved
	FigureName     string // carries a figure name
	XaxisName      string // carries X-axis name
	YaxisName      string // carries Y-axis name
	xmin           *int64 // optional: sets a minimum boundary for X-axis
	xmax           *int64 // optional: sets a maximum boundary for X-axis
	ymin           *int64 // optional: sets a minimum boundary for Y-axis
	ymax           *int64 // optional: sets a maximum boundary for Y-axis
	gridOn         bool   // enable grid on the figure
}

// InitializeDrawStruct initializes a Draw structure
func (d *Draw) InitializeDrawStruct() {
	d.Rendered = false
	d.OutputFileName = "Default"
	d.FigureName = ""
	d.XaxisName = ""
	d.YaxisName = ""
	d.gridOn = true // enabling grid by default
}

// SetOutputFileName sets a filename of the rendered picture
func (d *Draw) SetOutputFileName(fileName string) {
	d.OutputFileName = fileName
}

// SetFigureName sets a filename of the rendered picture
func (d *Draw) SetFigureName(fn string) {
	d.FigureName = fn
}

// SetXaxisName sets a filename of the rendered picture
func (d *Draw) SetXaxisName(xname string) {
	d.XaxisName = xname
}

// SetYaxisName sets a filename of the rendered picture
func (d *Draw) SetYaxisName(yname string) {
	d.YaxisName = yname
}

// SetXmin sets a minimum bound for an X-axis
func (d *Draw) SetXmin(xmin int64) {
	*d.xmin = xmin
}

// SetXmax sets a maximum bound for an X-axis
func (d *Draw) SetXmax(xmax int64) {
	*d.xmax = xmax
}

// SetYmin sets a minimum bound for an Y-axis
func (d *Draw) SetYmin(ymin int64) {
	*d.ymin = ymin
}

// SetYmax sets a maximum bound for an Y-axis
func (d *Draw) SetYmax(ymax int64) {
	*d.ymax = ymax
}

// Coordinates is a structure that carries all information about the nodes and their coordinates in
// the systemmodel.SystemModel structure
type Coordinates struct {
	Points map[string]*Coordinate // this map contains as a key Name of the node and it's coordinates..
}

// Coordinate structure is a hybrid structure between systemmodel.SystemModel and the custom plotter.XYer structure.
// It serves as an intermediate layer to convert data from SystemModel to plotter-friendly data and
// ensure drawing of the graph
type Coordinate struct {
	Coordinates plotter.XYs // coordinates of the point
}

// ConvertSystemModelToDrawStruct converts SystemModel to a plotter-friendly structure, which holds information about
// coordinates of each node
func (ds *Coordinates) ConvertSystemModelToDrawStruct(sm *systemmodel.SystemModel) {
	ds.Points = make(map[string]*Coordinate, sm.GetTotalNumberOfInstances())
	for i := 1; i <= len(sm.Layers); i++ {
		layer := sm.Layers[int32(i)]
		j := 1
		for _, v := range layer.Instances {
			// FIXME: this assignment of coordinates here may be a potential source of issues in the graph..
			var data plotter.XYs
			if v.Type == systemmodel.VI && strings.HasPrefix(v.Name, "MAIS") {
				data = createRootNodePoints()
			} else {
				data = createPoints(i, len(sm.Layers), j, len(layer.Instances))
			}
			dp := &Coordinate{
				Coordinates: data,
			}
			ds.Points[v.Name] = dp
			j++
		}
	}
}

// ExtractPoints extracts an array of all points in graph
func (ds *Coordinates) ExtractPoints() []plotter.XYer {
	list := make([]plotter.XYer, 0)
	for _, v := range ds.Points {
		list = append(list, v.Coordinates)
	}

	return list
}

// DrawSystemModel draws a figure representing provided SystemModel. It is done in the following way:
//
//	Iterate over layers, over each node in the layer.
//	Once you've reached a node, create its coordinates (X, Y). After that iterate over each related instance (connection) and
//	create a coordinates (X, Y) for each relation. Then, add corresponding lines between the node and its relations.
//	Next, iterate over the rest of the nodes.
//	To avoid duplication of the nodes, create a custom structure, which carries coordinates (X, Y) information, name of the node (which is unique)
//	and the status (were coordinates (X, Y) already created?).
func (d *Draw) DrawSystemModel(sm *systemmodel.SystemModel) error {

	p := plot.New()

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
		layer := sm.Layers[int32(i)]
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

	// extracting data points to plot on a figure
	pts := ds.ExtractPoints()
	pnts, err := plotutil.NewErrorPoints(plotutil.MedianAndMinMax, pts...)
	if err != nil {
		return err
	}
	err = plotutil.AddScatters(p, d.FigureName, pnts)
	if err != nil {
		return err
	}
	// using plot.GlyphBoxer to draw nodes in a figure
	p.Add(plotter.NewGlyphBoxes())
	// Save the plot to a PNG file. // ToDo - customize dimensions..
	if err := p.Save(20*vg.Inch, 20*vg.Inch, "figures/"+d.OutputFileName+".png"); err != nil {
		return err
	}
	d.Rendered = true
	return nil
}

// createRootNodePoints creates a root node, MAIS
func createRootNodePoints() plotter.XYs {
	data := make(plotter.XYs, 1)
	data[0].X = 100.0 // we know that it's an only node by definition
	data[0].Y = 100.0

	return data
}

// createPoints creates points with regard to given input
func createPoints(currentLevel int, levels int, itemNumber int, density int) plotter.XYs {
	data := make(plotter.XYs, 1)
	data[0].X = float64(itemNumber) * float64(200) / float64(density+1)
	data[0].Y = float64(levels-currentLevel+1) * float64(100/levels)

	return data
}
