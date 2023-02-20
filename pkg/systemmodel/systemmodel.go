// Package systemmodel provides a basic structure to carry System Model information
// and contains basic functions for its initialization
package systemmodel

import (
	"math/rand"
	"strings"
)

// SystemModel structure represents structure of the system model
type SystemModel struct {
	Depth  int32            // represents depth of the system model
	Layers map[int32]*Layer // represents list of layers in the system model
	// a key for the Application can be whatever string you want (e.g., name of the application)
	Applications map[string]*Application // represents a list of applications, which were deployed at this layer - this is to track all deployed applications over all layers
	// ToDo - how to store relations of the nodes?? That breaks everything...
}

// Layer structure represents layer of the system model (e.g., Layer[3] corresponds to the 3-rd level of the SystemModel)
type Layer struct {
	DeployedApps  map[string]*Application // represents a list of applications, which were deployed at this layer
	VIwasDeployed bool                    // this is to indicate whether VI was deployed at this Layer
}

// Application structure defines initial set of rules for applications, probabilities of the application deployment and
// carries a map of the deployed apps
type Application struct {
	Rules       int     // number of instances that application can deploy
	Probability float32 // probability of the application deployment
	State       bool    // true for deployed, false for not deployed
}

// InitializeSystemModel initializes SystemModel with provided values
func (sm *SystemModel) InitializeSystemModel(numApps int32, deps int32) {
	sm.Depth = deps
	sm.Layers = make(map[int32]*Layer, numApps)
	sm.Applications = make(map[string]*Application, numApps)
}

// CreateLayer creates a new layer of the fractal system model
func (l *Layer) CreateLayer(apps map[string]*Application) map[string]*Application {

	resApps := make(map[string]*Application, len(apps))

	for k, v := range apps {
		// copying an initial value
		resApps[k] = v
		// firstly, checking if the app was deployed
		// VI has a special rule - it can be deployed more than once
		if v.State && !strings.Contains(k, "VI") { // if application was deployed, then skip this iteration
			continue
		}

		// generate random probability
		probability := rand.Float32() // consider using normal distribution with rang.NormFloat64()
		// if probability of the application deployment is lower than the generated probability, then the application should be deployed
		if probability < v.Probability {
			app := &Application{
				Rules:       v.Rules,
				Probability: v.Probability,
				State:       true,
			}
			l.DeployedApps[k] = app
			resApps[k] = app // updating a value
			if strings.Contains(k, "VI") {
				l.VIwasDeployed = true
			}
		}
	}

	return resApps
}

// AddLayer adds layer to the system model at a given level
func (sm *SystemModel) AddLayer(layer *Layer, level int32) {
	sm.Layers[level] = layer
}

// CreateApplication creates an application with given parameters and adds it to the list
func (sm *SystemModel) CreateApplication(numInstances int, probability float32, name string) {
	application := &Application{
		Rules:       numInstances,
		Probability: probability,
		State:       false,
	}
	sm.Applications[name] = application
}

// CreateRandomApplications creates a set of applications with random parameters given the pre-defined names
func (sm *SystemModel) CreateRandomApplications(names []string, maxNumInstances int) {
	sm.Applications = make(map[string]*Application, len(names))
	probabilitySum := float32(1)
	for _, name := range names {
		// generate random probability of the application deployment within a range
		probability := rand.Float32() * probabilitySum   // consider using normal distribution with rang.NormFloat64()
		numInstances := rand.Intn(maxNumInstances-1) + 1 // it should deploy at least 1 instance
		sm.CreateApplication(numInstances, probability, name)
		probabilitySum -= probability
	}
}

// GenerateRandomSystemModel generates a random system model with regard to the parameters provided at input (i.e., depth of a fractal,
// applications inside of the MAIS, maximum number of instances per application)
func (sm *SystemModel) GenerateRandomSystemModel(depth int32, appNames []string, maxNumInstances int) {
	sm.InitializeSystemModel(int32(len(appNames)), depth)
	sm.CreateRandomApplications(appNames, maxNumInstances)
	for i := 1; int32(i) <= depth; i++ {
		if i > 1 && !sm.Layers[int32(i-1)].VIwasDeployed {
			break
		}
		// ToDo - make sure that the next iteration is being entered only when the VI was deployed, otherwise it doesn't make sense..
		//  Also, think about a way how to keep relations between the instances and the nodes (VIs)..
		layer := &Layer{
			DeployedApps: make(map[string]*Application, 0),
		}
		sm.Applications = layer.CreateLayer(sm.Applications)
		sm.AddLayer(layer, int32(i))
	}
}
