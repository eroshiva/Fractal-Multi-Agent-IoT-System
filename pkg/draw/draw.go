// Package draw implements a set of helper functions to draw SystemModel
// Here are stored functions which manipulate with the internal structures and help to plot a figure
package draw

import (
	"fmt"
	"gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/pkg/systemmodel"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"strings"
)

// Draw structure holds all necessary information for plotting a figure of SystemModel
type Draw struct {
	Rendered       bool      // indicates whether a figure was rendered or not
	OutputFileName string    // an output file name, where the figure would be saved
	FigureName     string    // carries a figure name
	XaxisName      string    // carries X-axis name
	YaxisName      string    // carries Y-axis name
	xmin           *int64    // optional: sets a minimum boundary for X-axis
	xmax           *int64    // optional: sets a maximum boundary for X-axis
	ymin           *int64    // optional: sets a minimum boundary for Y-axis
	ymax           *int64    // optional: sets a maximum boundary for Y-axis
	gridOn         bool      // enable grid on the figure
	XLength        vg.Length // sets length of an X-axis (in Inches, Cm or mm)
	YLength        vg.Length // sets length of an Y-axis (in Inches, Cm or mm)
}

// InitializeDrawStruct initializes a Draw structure
func (d *Draw) InitializeDrawStruct() *Draw {
	d.Rendered = false
	d.gridOn = true // enabling grid by default
	d.SetOutputFileName("Default").SetFigureName("").SetXaxisName("").
		SetYaxisName("").SetXLength(20 * vg.Inch).SetYLength(20 * vg.Inch)
	return d
}

// SetOutputFileName sets a filename of the rendered picture
func (d *Draw) SetOutputFileName(fileName string) *Draw {
	d.OutputFileName = fileName
	return d
}

// SetFigureName sets a filename of the rendered picture
func (d *Draw) SetFigureName(fn string) *Draw {
	d.FigureName = fn
	return d
}

// SetXaxisName sets a filename of the rendered picture
func (d *Draw) SetXaxisName(xname string) *Draw {
	d.XaxisName = xname
	return d
}

// SetYaxisName sets a filename of the rendered picture
func (d *Draw) SetYaxisName(yname string) *Draw {
	d.YaxisName = yname
	return d
}

// SetXmin sets a minimum bound for an X-axis
func (d *Draw) SetXmin(xmin int64) *Draw {
	*d.xmin = xmin
	return d
}

// SetXmax sets a maximum bound for an X-axis
func (d *Draw) SetXmax(xmax int64) *Draw {
	*d.xmax = xmax
	return d
}

// SetYmin sets a minimum bound for an Y-axis
func (d *Draw) SetYmin(ymin int64) *Draw {
	*d.ymin = ymin
	return d
}

// SetYmax sets a maximum bound for an Y-axis
func (d *Draw) SetYmax(ymax int64) *Draw {
	*d.ymax = ymax
	return d
}

// SetXLength sets actual length of an X-axis in inches, cm or mm
func (d *Draw) SetXLength(xl vg.Length) *Draw {
	d.XLength = xl
	return d
}

// SetYLength sets actual length of an Y-axis in inches, cm or mm
func (d *Draw) SetYLength(yl vg.Length) *Draw {
	d.YLength = yl
	return d
}

// Coordinates is a structure that carries all information about the nodes and their coordinates in
// the systemmodel.SystemModel structure
type Coordinates struct {
	Points map[string]*Coordinate // this map contains as a key Name of the node and it's coordinates..
	Labels plotter.XYLabels       // this is to hold a name of the instance and plot it on the graph (per node)
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
	//ds.Labels = make(plotter.XYLabels, 0)
	for i := 1; i <= len(sm.Layers); i++ {
		layer := sm.Layers[i]
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
			// adding labels to figure
			ds.Labels.XYs = append(ds.Labels.XYs, data[0])
			ds.Labels.Labels = append(ds.Labels.Labels, v.Name)
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

// item1 type mimics type plotter.item - this is done to adjust function AddScatterSquare
type item1 struct {
	name  string
	value plot.Thumbnailer
}

// AddScattersSquare adds a Scatter plotters to a plot.
// The variadic arguments must be either strings
// or plotter.XYers.  Each plotter.XYer is added to
// the plot using the next color, and square glyph shape
// via the Color and Shape functions. If a
// plotter.XYer is immediately preceded by
// a string then a legend entry is added to the plot
// using the string as the name.
//
// If an error occurs then none of the plotters are added
// to the plot, and the error is returned.
func AddScattersSquare(plt *plot.Plot, vs ...interface{}) error {
	var ps []plot.Plotter
	var items []item1
	var i int
	for _, v := range vs {
		switch t := v.(type) {
		case map[string]*Coordinate:
			for k, val := range t {
				s, err := plotter.NewScatter(val.Coordinates)
				if err != nil {
					return err
				}
				if strings.Contains(k, "MAIS") {
					s.Color = plotutil.Color(2)
				} else if strings.HasPrefix(k, "VI") {
					s.Color = plotutil.Color(1)
				} else {
					s.Color = plotutil.Color(0)
				}
				s.Shape = plotutil.Shape(6) // filled with color box
				s.Radius = 0.25 * vg.Centimeter
				i++
				ps = append(ps, s)
			}
			// adding legend to the figure
			// for MAIS
			mais, err := plotter.NewScatter(t["MAIS"].Coordinates)
			if err != nil {
				return err
			}
			mais.Color = plotutil.Color(2) // setting a color to correspond to MAIS colour
			mais.Shape = plotutil.Shape(6) // filled with color box
			mais.Radius = 0.2 * vg.Centimeter
			items = append(items, item1{name: "MAIS", value: mais})
			// for VI
			vi, err := plotter.NewScatter(t["MAIS"].Coordinates)
			if err != nil {
				return err
			}
			vi.Color = plotutil.Color(1) // setting a color to correspond to VI colour
			vi.Shape = plotutil.Shape(6) // filled with color box
			vi.Radius = 0.2 * vg.Centimeter
			items = append(items, item1{name: "VI", value: vi})
			// for App
			app, err := plotter.NewScatter(t["MAIS"].Coordinates)
			if err != nil {
				return err
			}
			app.Color = plotutil.Color(0) // setting a color to correspond to App colour
			app.Shape = plotutil.Shape(6) // filled with color box
			app.Radius = 0.2 * vg.Centimeter
			items = append(items, item1{name: "App", value: app})

		default:
			panic(fmt.Sprintf("plotutil: AddScattersSquare only map[string]*Coordinate type, got %T", t))
		}
	}
	plt.Add(ps...)
	for _, v := range items {
		plt.Legend.Add(v.name, v.value)
	}
	return nil
}
