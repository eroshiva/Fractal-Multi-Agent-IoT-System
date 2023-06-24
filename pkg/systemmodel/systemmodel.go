// Package systemmodel provides a basic structure to carry System Model information
// and contains basic functions for its initialization
package systemmodel

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
)

// SystemModel structure represents structure of the system model
type SystemModel struct {
	Depth  int            // represents depth of the system model
	Layers map[int]*Layer // represents list of layers in the system model
	// a key for the Application can be whatever string you want (e.g., name of the application)
	Applications map[string]*Application // represents a list of applications, which were deployed at this layer - this is to track all deployed applications over all layers
	VIcount      *uint64                 // this pointer is used to be passed to the various functions and change its value once VI is being deployed. This is done to distinguish various VIs on the same layer
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
	Rules       int               // number of instances that application can deploy
	Probability float32           // probability of the application deployment
	State       bool              // true for deployed, false for not deployed
	Aspect      map[string]string // holds aspects of the Application, like Reliability or its Priority (= weight)
}

// InitializeSystemModel initializes SystemModel with provided values
func (sm *SystemModel) InitializeSystemModel(numApps int, depth int) *SystemModel {
	sm.Depth = depth
	sm.Layers = make(map[int]*Layer, depth)
	sm.Applications = make(map[string]*Application, numApps)
	viCount := uint64(0)
	sm.VIcount = &viCount
	return sm
}

// CreateInstance creates an empty Instance with given name and instance type
func (i *Instance) CreateInstance(name string, tp InstanceType) *Instance {
	i.Name = name
	i.Type = tp
	i.Relations = make([]*Instance, 0)
	i.Aspect = make(map[string]string, 0)
	return i
}

// CreateInstanceRnd creates an Instance with given name, instance type and random parameters (priority and chain coefficient)
func (i *Instance) CreateInstanceRnd(name string, tp InstanceType, cc float64) *Instance {
	i.Name = name
	i.Type = tp
	i.Relations = make([]*Instance, 0)
	// setting some aspects
	i.Aspect = make(map[string]string, 0)
	priority := rand.Float64()
	i.SetPriority(priority)
	i.SetChainCoefficient(cc * priority)
	return i
}

// AddRelation adds an instance to the instance list (i.e., Relations)
func (i *Instance) AddRelation(relation *Instance) *Instance {
	i.Relations = append(i.Relations, relation)
	return i
}

// AddAspect adds an Aspect to the Instance list (i.e., reliability, etc.)
func (i *Instance) AddAspect(aspectType string, aspectValue string) *Instance {
	i.Aspect[aspectType] = aspectValue
	return i
}

// InitializeLayer initializes Layer
func (l *Layer) InitializeLayer() *Layer {
	l.Instances = make([]*Instance, 0)
	l.VIwasDeployed = false
	return l
}

// AddLayer adds layer to the system model at a given level
func (sm *SystemModel) AddLayer(layer *Layer, level int) *SystemModel {
	sm.Layers[level] = layer
	return sm
}

// AddInstanceToLayer adds a given Instance to the Layer and checks if it is of type VI
// to indicate that VI was deployed at this Layer
func (l *Layer) AddInstanceToLayer(instance *Instance) *Layer {
	l.Instances = append(l.Instances, instance)
	// checking if an instance is of type VI
	if instance.Type == VI {
		l.VIwasDeployed = true
	}
	return l
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
func (sm *SystemModel) CreateApplication(numInstances int, probability float32, name string) *SystemModel {
	application := &Application{
		Rules:       numInstances,
		Probability: probability,
		State:       false,
		Aspect:      make(map[string]string, 0),
	}
	sm.Applications[name] = application
	return sm
}

// CreateRandomApplications creates a set of applications with random parameters given the pre-defined names
func (sm *SystemModel) CreateRandomApplications(names []string, minNumInstances int, maxNumInstances int) *SystemModel {
	sm.Applications = make(map[string]*Application, len(names))
	probabilitySum := float32(1)
	for _, name := range names {
		// generate random probability of the application deployment within a range
		probability := rand.Float32() * probabilitySum // consider using normal distribution with rand.NormFloat64()
		// taking care of the case when the only one instance resides within Application
		var numInstances = 1
		if maxNumInstances-minNumInstances > 0 {
			numInstances = rand.Intn(maxNumInstances-minNumInstances) + minNumInstances
		}
		sm.CreateApplication(numInstances, probability, name)
		probabilitySum -= probability
	}
	return sm
}

// Deploy function sets the state of the Application to true
func (a *Application) Deploy() *Application {
	a.State = true
	return a
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
func (i *Instance) DeployApplications(apps map[string]*Application, currentLevel int, viCount *uint64) (bool, map[string]*Application) {
	updatedApps := make(map[string]*Application, len(apps))
	viWasDeployed := false

	for appName, app := range apps {
		// if the application was not yet deployed, or it is a VI (which can be deployed multiple times)
		if !app.State || strings.HasPrefix(appName, "VI") {
			updatedApp := &Application{
				Rules:       app.Rules,
				Probability: app.Probability,
				State:       false,
			}
			// specific to VI - if it was already deployed, then keep its flag as deployed..
			if strings.HasPrefix(appName, "VI") && app.State {
				updatedApp.State = true
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
				for j := 1; j <= app.Rules; j++ {
					appInstance := &Instance{}
					name := appName + "-" + strconv.Itoa(j) + "-" + strconv.FormatInt(int64(currentLevel), 10)
					if strings.HasPrefix(appName, "VI") {
						name = "VI#" + strconv.FormatUint(*viCount, 10) + "-" + strconv.Itoa(j) + "-" + strconv.FormatInt(int64(currentLevel), 10)
					}
					appInstance.CreateInstance(name, tp)
					i.AddRelation(appInstance)
				}
				updatedApp.State = true
			}
			updatedApps[appName] = updatedApp
		} else {
			// if application was already deployed, just save it
			updatedApps[appName] = app
		}
	}

	return viWasDeployed, updatedApps
}

// CreateLayer creates a layer of the SystemModel and updates the Applications list to reflect the current deployment state
func (l *Layer) CreateLayer(apps map[string]*Application, currentLevel int, viCount *uint64) (map[string]*Application, *Layer) {
	nextLayer := &Layer{}
	nextLayer.InitializeLayer()
	updApps := apps

	for _, instance := range l.Instances {
		// checking if the instance is of type VI or the root instance, MAIS
		if strings.HasPrefix(instance.Name, "VI") || strings.Contains(instance.Name, "MAIS") {
			viWasDeployed, updatedApps := instance.DeployApplications(updApps, currentLevel, viCount)
			if viWasDeployed {
				l.VIwasDeployed = true
			}
			// overwrite the apps with regard to what was deployed
			updApps = updatedApps
			// adding instances to layer
			for _, inst := range instance.Relations {
				nextLayer.AddInstanceToLayer(inst)
			}
		}
	}

	return updApps, nextLayer
}

// GenerateSystemModel generates system model with regard to provided input data
func (sm *SystemModel) GenerateSystemModel() *SystemModel {
	sm.InitializeRootLayer()
	for i := 2; i <= sm.Depth; i++ {
		if sm.Layers[i-1].VIwasDeployed {
			apps, nextLayer := sm.Layers[i-1].CreateLayer(sm.Applications, i, sm.VIcount)
			sm.Applications = apps // updating Applications map
			// if something was deployed, then add Layer to the SystemModel, otherwise stop
			if len(nextLayer.Instances) > 0 {
				sm.AddLayer(nextLayer, i)
			} else {
				break
			}
		} else {
			break
		}
	}
	return sm
}

// InitializeRootLayer initializes 1st level Layer of SystemModel as a single instance with given name "MAIS",
// which behaves as a VI.
func (sm *SystemModel) InitializeRootLayer() *SystemModel {
	rootInstance := &Instance{}
	rootInstance.CreateInstance("MAIS", CreateInstanceTypeVI())
	rootLayer := &Layer{}
	rootLayer.InitializeLayer()
	rootLayer.VIwasDeployed = true
	*sm.VIcount++
	rootLayer.AddInstanceToLayer(rootInstance)
	sm.AddLayer(rootLayer, 1)
	return sm
}

// PrettyPrintApplications prints Application related information
func (sm *SystemModel) PrettyPrintApplications() *SystemModel {
	for k, v := range sm.Applications {
		log.Printf("%s has probability %v and deploys %v instances. Number of aspects is %d."+
			" Deployed status: %v\n", k, v.Probability, v.Rules, len(v.Aspect), v.State)
		for key, value := range v.Aspect {
			log.Printf("Aspect %s, value %s\n", key, value)
		}
	}
	return sm
}

// PrettyPrintLayers prints Layers related information
func (sm *SystemModel) PrettyPrintLayers() *SystemModel {
	for k := 1; k <= len(sm.Layers); k++ {
		v := sm.Layers[k]
		fmt.Printf("----> Layer %v, VI deployed %v, Instances deployed %v, detailed info about deployed instances:\n", k, v.VIwasDeployed, len(v.Instances))
		v.PrettyPrintLayer()
	}
	return sm
}

// PrettyPrintLayer prints Layer related information
func (l *Layer) PrettyPrintLayer() {
	for _, v := range l.Instances {
		if len(v.Relations) > 0 {
			fmt.Printf("--> Instance %s of type %v has %v following relations:\n", v.Name, v.Type, len(v.Relations))
			for _, val := range v.Relations {
				fmt.Printf("Related is Instance %s of type %v\n", val.Name, val.Type)
			}
			if len(v.Aspect) > 0 {
				fmt.Printf("-> Instance %s of type %v has %d aspects, they are:\n", v.Name, v.Type, len(v.Aspect))
				for k, v := range v.Aspect {
					fmt.Printf("Aspect %s, value %s\n", k, v)
				}
			} else {
				fmt.Printf("-> Instance %s of type %v has 0 aspects\n", v.Name, v.Type)
			}
		} else {
			fmt.Printf("--> Instance %s of type %v has no relations and %d aspects\n", v.Name, v.Type, len(v.Aspect))
			for k, v := range v.Aspect {
				fmt.Printf("Aspect %s, value %s\n", k, v)
			}
		}
	}
}

// GetTotalNumberOfInstances gets total number of instances (i.e., nodes) in the SystemModel structure
func (sm *SystemModel) GetTotalNumberOfInstances() int64 {
	var total int64
	for k := 1; k <= len(sm.Layers); k++ {
		v := sm.Layers[k]
		total += int64(len(v.Instances))
	}
	return total
}

// GetTheGreatestNumberOfInstancesPerLayer gets the greatest number of instances per layer
func (sm *SystemModel) GetTheGreatestNumberOfInstancesPerLayer() int64 {
	var max int64
	for k := 1; k <= len(sm.Layers); k++ {
		if max < int64(len(sm.Layers[k].Instances)) {
			max = int64(len(sm.Layers[k].Instances))
		}
	}
	return max
}

// GenerateAppNames generates Application names as a string composition of "App#" and its ordering number
func GenerateAppNames(maxNumInstances int) []string {
	res := make([]string, maxNumInstances+1)
	res[0] = "VI"
	for i := 1; i <= maxNumInstances; i++ {
		res[i] = "App#" + strconv.FormatInt(int64(i), 10)
	}

	return res
}

// GetSystemModelParameters iterates over a provided data variable (assuming that it is a benchmarked data) and
// determines maximum Depth, maximum Applications number and maximum number of Instances (all defined as input data for SystemModel)
func GetSystemModelParameters(data map[int]map[int]map[int]float64) (int, int, int, error) {

	var depth, apps, instances int

	for k, v := range data {
		for k1, v1 := range v {
			for k2 := range v1 {
				if instances < k2 {
					instances = k2
				}
			}
			if apps < k1 {
				apps = k1
			}
		}
		if depth < k {
			depth = k
		}
	}

	if depth == 0 || apps == 0 || instances == 0 {
		return -1, -1, -1, fmt.Errorf("something went wrong during determination os SystemModel parameters - "+
			"probably empty data were passed: %v\n", data)
	}

	return depth, apps, instances, nil
}

// GetAspect function returns an Aspect of an Instance with a given key
func (i *Instance) GetAspect(key string) (string, error) {
	var aspect string
	aspect, ok := i.Aspect[key]
	if !ok {
		return "", fmt.Errorf("can't find '%v' aspect for instance %v", key, i.Name)
	}
	return aspect, nil
}

// SetAspect function sets an Aspect for an Instance with a given key and a given value
func (i *Instance) SetAspect(key, value string) *Instance {
	i.Aspect[key] = value
	return i
}

// GetAppName returns an application name, which can be used as a key to get the Application out of
// the applications list. It assumes that all Application instances in the name have dashes
func (i *Instance) GetAppName() (string, error) {
	if i.IsVI() {
		// if it is a VI, then returning the exact VI instance name
		return i.Name, nil
	}
	// if it is an Application, doing some parsing..
	idx := strings.Index(i.Name, "-") // getting the index of a first dash
	subStr := i.Name[idx+1:]
	idx2 := strings.Index(subStr, "-") // getting the index of a second dash
	if idx != -1 && idx2 != -1 {
		// obtaining actual application number
		appNumber := i.Name[idx+1 : idx+idx2+1]
		return "App#" + appNumber, nil
	}
	return "", fmt.Errorf("can't parse first dash in the instance name %s", i.Name)
}

// GetInstanceNumber returns an instance number for Application or VI.
// It assumes that all Application (and VI) instances in the name have dashes
func (i *Instance) GetInstanceNumber() (int64, error) {
	idx := strings.LastIndex(i.Name, "-") // getting the last index of a dash
	if idx != -1 {
		indexInt, err := strconv.ParseInt(i.Name[idx+1:], 10, 64)
		if err != nil {
			return -1, fmt.Errorf("conversion of an instance index to integer failed: %w", err)
		}
		return indexInt, nil
	}
	return -1, fmt.Errorf("something went wrong when attempted to parse a %s instance number", i.Name)
}

// IsVI function returns true, if the instance is of type VI, otherwise false
func (i *Instance) IsVI() bool {
	return i.Type == VI
}

// IsApp function returns true, if the instance is of type Application, otherwise false
func (i *Instance) IsApp() bool {
	return i.Type == App
}
