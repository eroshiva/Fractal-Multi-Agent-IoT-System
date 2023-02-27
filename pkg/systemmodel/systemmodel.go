// Package systemmodel provides a basic structure to carry System Model information
// and contains basic functions for its initialization
package systemmodel

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

// SystemModel structure represents structure of the system model
type SystemModel struct {
	Depth  int32            // represents depth of the system model
	Layers map[int32]*Layer // represents list of layers in the system model
	// a key for the Application can be whatever string you want (e.g., name of the application)
	Applications map[string]*Application // represents a list of applications, which were deployed at this layer - this is to track all deployed applications over all layers
	VIcount      *int64                  // this pointer is used to be passed to the various functions and change its value once VI is being deployed. This is done to distinguish various VIs on the same layer
}

// Layer structure represents the layer of the system model (e.g., Layer[3] corresponds to the 3-rd level of the SystemModel)
type Layer struct {
	Instances     []*Instance // represents deployed instances on this layer
	VIwasDeployed bool        // this is to indicate whether VI was deployed at this Layer
}

// Instance structure represents an instance
type Instance struct {
	Name      string            // carries a name of the instance (e.g., VI#2-1, App#3-2-4, etc...)
	Type      InstanceType      // specifies a type of the instance
	Relations []*Instance       // carries relations to the other instances
	Aspect    map[string]string // carries aspects of the Instance (e.g., Reliability, etc.)
}

// InstanceType defines a type of the instance. It is either VI, or Application (App)
type InstanceType uint

const (
	VI  InstanceType = 0 // VI indicates a type of Virtual Instance
	App InstanceType = 1 // App indicates a type of Application
)

// Application structure defines initial set of rules for applications, probabilities of the application deployment and
// carries a map of the deployed apps
type Application struct {
	Rules       int32   // number of instances that application can deploy
	Probability float32 // probability of the application deployment
	State       bool    // true for deployed, false for not deployed
}

// InitializeSystemModel initializes SystemModel with provided values
func (sm *SystemModel) InitializeSystemModel(numApps int32, depth int32) {
	sm.Depth = depth
	sm.Layers = make(map[int32]*Layer, depth)
	sm.Applications = make(map[string]*Application, numApps)
	viCount := int64(0)
	sm.VIcount = &viCount
}

// CreateInstance creates an instance with given name and instance type
func (i *Instance) CreateInstance(name string, tp InstanceType) {
	i.Name = name
	i.Type = tp
	i.Relations = make([]*Instance, 0)
	i.Aspect = make(map[string]string, 0)
}

// AddRelation adds an instance to the instance list (i.e., Relations)
func (i *Instance) AddRelation(relation *Instance) {
	i.Relations = append(i.Relations, relation)
}

// AddAspect adds an Aspect to the Instance list (i.e., reliability, etc.)
func (i *Instance) AddAspect(aspectType string, aspectValue string) {
	i.Aspect[aspectType] = aspectValue
}

// InitializeLayer initializes Layer
func (l *Layer) InitializeLayer() {
	l.Instances = make([]*Instance, 0)
	l.VIwasDeployed = false
}

// AddLayer adds layer to the system model at a given level
func (sm *SystemModel) AddLayer(layer *Layer, level int32) {
	sm.Layers[level] = layer
}

// AddInstanceToLayer adds a given Instance to the Layer and checks if it is of type VI
// to indicate that VI was deployed at this Layer
func (l *Layer) AddInstanceToLayer(instance *Instance) {
	l.Instances = append(l.Instances, instance)
	// checking if an instance is of type VI
	if instance.Type == VI {
		l.VIwasDeployed = true
	}
}

// CreateInstanceTypeVI creates an enumerator for VI type
func CreateInstanceTypeVI() InstanceType {
	return VI
}

// CreateInstanceTypeApp creates an enumerator for Application type
func CreateInstanceTypeApp() InstanceType {
	return App
}

// CreateApplication creates an application with given parameters and adds it to the list
func (sm *SystemModel) CreateApplication(numInstances int32, probability float32, name string) {
	application := &Application{
		Rules:       numInstances,
		Probability: probability,
		State:       false,
	}
	sm.Applications[name] = application
}

// CreateRandomApplications creates a set of applications with random parameters given the pre-defined names
func (sm *SystemModel) CreateRandomApplications(names []string, maxNumInstances int32) {
	sm.Applications = make(map[string]*Application, len(names))
	probabilitySum := float32(1)
	for _, name := range names {
		// generate random probability of the application deployment within a range
		probability := rand.Float32() * probabilitySum     // consider using normal distribution with rang.NormFloat64()
		numInstances := rand.Int31n(maxNumInstances-1) + 1 // it should deploy at least 1 instance
		sm.CreateApplication(numInstances, probability, name)
		probabilitySum -= probability
	}
}

// DeployApplication checks if the Application is deployed. It generates random probability and compares with
// the probability of the Application deployment. If it is smaller than the Application deployment probability,
// then the Application is deployed. In the other case, Application is not deployed.
func (a *Application) DeployApplication() bool {
	probability := rand.Float32()
	return probability < a.Probability
}

// DeployApplications iterates over a map of Applications and checks, whether application is deployed or not.
// It returns updated list of Applications, which denotes the updated state of applications
func (i *Instance) DeployApplications(apps map[string]*Application, currentLevel int32, viCount *int64) (bool, map[string]*Application) {
	updatedApps := make(map[string]*Application, 0)
	viWasDeployed := false

	for appName, app := range apps {
		// if the application was not yet deployed or it is a VI (which can be deployed multiple times)
		if !app.State || strings.HasPrefix(appName, "VI") {
			updatedApp := &Application{
				Rules:       app.Rules,
				Probability: app.Probability,
				State:       false,
			}
			deployed := app.DeployApplication()
			if deployed {
				tp := CreateInstanceTypeApp()
				// if the Application/VI was deployed, creating a new instances
				if strings.HasPrefix(appName, "VI") {
					viWasDeployed = true
					tp = CreateInstanceTypeVI()
					*viCount++
				}
				for j := 1; int32(j) <= app.Rules; j++ {
					appInstance := &Instance{}
					name := appName + "-" + strconv.Itoa(j) + "-" + strconv.FormatInt(int64(currentLevel), 10)
					if strings.HasPrefix(appName, "VI") {
						name = "VI#" + strconv.FormatInt(*viCount, 10) + "-" + strconv.Itoa(j) + "-" + strconv.FormatInt(int64(currentLevel), 10)
					}
					appInstance.CreateInstance(name, tp)
					i.AddRelation(appInstance)
				}
				updatedApp.State = true
			}
			updatedApps[appName] = updatedApp
		}
	}

	return viWasDeployed, updatedApps
}

// CreateLayer creates a layer of the SystemModel and updates the Applications list to reflect the current deployment state
func (l *Layer) CreateLayer(apps map[string]*Application, currentLevel int32, viCount *int64) (map[string]*Application, *Layer) {

	nextLayer := &Layer{}
	nextLayer.InitializeLayer()

	for _, instance := range l.Instances {
		// checking if the instance is of type VI or the root instance, MAIS
		if strings.HasPrefix(instance.Name, "VI") || strings.Contains(instance.Name, "MAIS") {
			viWasDeployed, updatedApps := instance.DeployApplications(apps, currentLevel, viCount)
			if viWasDeployed {
				l.VIwasDeployed = true
			}
			// overwrite the apps with regard to what was deployed
			apps = updatedApps
			// adding instances to layer
			for _, inst := range instance.Relations {
				nextLayer.AddInstanceToLayer(inst)
			}
		}
	}

	return apps, nextLayer
}

// GenerateSystemModel generates system model with regard to provided input data
func (sm *SystemModel) GenerateSystemModel() {
	sm.InitializeRootLayer()
	for i := 2; int32(i) <= sm.Depth; i++ {
		if sm.Layers[int32(i)-1].VIwasDeployed {
			var nextLayer *Layer
			sm.Applications, nextLayer = sm.Layers[int32(i)-1].CreateLayer(sm.Applications, int32(i), sm.VIcount)
			// if something was deployed, then add Layer to the SystemModel, otherwise stop
			if len(nextLayer.Instances) > 0 {
				sm.AddLayer(nextLayer, int32(i))
			} else {
				break
			}
		} else {
			break
		}
	}
}

// InitializeRootLayer initializes 1st level Layer of SystemModel as a single instance with given name "MAIS",
// which behaves as a VI.
func (sm *SystemModel) InitializeRootLayer() {
	rootInstance := &Instance{}
	rootInstance.CreateInstance("MAIS", CreateInstanceTypeVI())
	rootLayer := &Layer{}
	rootLayer.InitializeLayer()
	rootLayer.VIwasDeployed = true
	*sm.VIcount++
	rootLayer.AddInstanceToLayer(rootInstance)
	sm.AddLayer(rootLayer, 1)
}

// PrettyPrintApplications prints Application related information
func (sm *SystemModel) PrettyPrintApplications() {
	for k, v := range sm.Applications {
		fmt.Printf("%s has probability %v and deploys %v instances. Deployed status: %v\n", k, v.Probability, v.Rules, v.State)
	}
}

// PrettyPrintLayers prints Layers related information
func (sm *SystemModel) PrettyPrintLayers() {
	for k := 1; k <= len(sm.Layers); k++ {
		v := sm.Layers[int32(k)]
		fmt.Printf("----> Layer %v, VI deployed %v, Instances deployed %v, detailed info about deployed instances:\n", k, v.VIwasDeployed, len(v.Instances))
		v.PrettyPrintLayer()
	}
}

// PrettyPrintLayer prints Layer related information
func (l *Layer) PrettyPrintLayer() {
	for _, v := range l.Instances {
		if len(v.Relations) > 0 {
			fmt.Printf("--> Instance %s of type %v has %v following relations:\n", v.Name, v.Type, len(v.Relations))
			for _, val := range v.Relations {
				fmt.Printf("Related is Instance %s of type %v\n", val.Name, val.Type)
			}
		} else {
			fmt.Printf("--> Instance %s of type %v has no relations\n", v.Name, v.Type)
		}
	}
}

// GetTotalNumberOfInstances gets total number of instances (i.e., nodes) in the SystemModel structure
func (sm *SystemModel) GetTotalNumberOfInstances() int64 {
	var total int64
	for k := 1; k <= len(sm.Layers); k++ {
		v := sm.Layers[int32(k)]
		total += int64(len(v.Instances))
	}
	return total
}
