// Package systemmodel implements means of Fractal MAIS system model. This file in particular holds all functions
// related to the reliability estimation and supports meertcore package with necessary functionality.
package systemmodel

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

const reliabilityKey = "Reliability"
const priorityKey = "Priority"
const chainCoefKey = "ChainCoefficient" // represents a chain coefficient, which is ultimately a multiplication of all priorities on top of the instance

// GatherApplicationInstanceReliabilities gathers reliability for all entities of an application.
// This function does not work for VI
func (sm *SystemModel) GatherApplicationInstanceReliabilities(appName string) (map[string]float64, error) {
	app, ok := sm.Applications[appName]
	if !ok {
		sm.PrettyPrintApplications()
		return nil, fmt.Errorf("application %s was not initialized in a SystemModel", appName)
	}
	if app.State && !strings.HasPrefix(appName, "VI") {
		res := make(map[string]float64, app.Rules)
		var appReliability float64
		count := app.Rules
		for i := len(sm.Layers); i > 0; i-- {
			layer, ok := sm.Layers[i]
			if !ok {
				return nil, fmt.Errorf("map entry for SystemModel.Layers with key %v does not exist", i)
			}
			for _, v := range layer.Instances {
				if strings.HasPrefix(v.Name, appName+"-") {
					reliability, err := v.GetReliability()
					if err != nil {
						return nil, err
					}
					res[v.Name] = reliability
					priority, err := v.GetPriority()
					if err != nil {
						return nil, err
					}
					appReliability += reliability * priority
					count--
				}
				if count == 0 {
					break
				}
			}
		}
		sm.Applications[appName].SetReliability(appReliability)
		return res, nil
	}

	return nil, fmt.Errorf("application %s was not deployed", appName)
}

// GatherAllApplicationsReliabilities function gathers reliability of each application and returns it in a map
func (sm *SystemModel) GatherAllApplicationsReliabilities() (map[string]float64, error) {
	reliabilities := make(map[string]float64, 0)

	for k, v := range sm.Applications {
		if !strings.HasPrefix(k, "VI") {
			_, err := sm.GatherApplicationInstanceReliabilities(k)
			if err != nil {
				return nil, fmt.Errorf("application %s: %w", k, err)
			}
			rel, err := v.GetReliability()
			if err != nil {
				return nil, fmt.Errorf("application %s: %w", k, err)
			}
			reliabilities[k] = rel
		}
	}

	return reliabilities, nil
}

// SetReliability sets Reliability Aspect for an Instance
func (i *Instance) SetReliability(reliability float64) *Instance {
	if i.Aspect == nil {
		i.Aspect = make(map[string]string, 0)
	}
	i.Aspect[reliabilityKey] = strconv.FormatFloat(reliability, 'f', -1, 64)
	return i
}

// GetReliability returns Reliability Aspect of an Instance
func (i *Instance) GetReliability() (float64, error) {
	relStr, ok := i.Aspect[reliabilityKey]
	if !ok {
		return 0, fmt.Errorf("looks like Reliability aspect for instance %v was not defined.. "+
			"Empty output for a key %v\n%v", i.Name, reliabilityKey, i)
	}
	reliability, err := strconv.ParseFloat(relStr, 64)
	if err != nil {
		return 0, fmt.Errorf("can't parse string %s to float64", relStr)
	}
	return reliability, nil
}

// SetPriority sets Priority Aspect for an Instance
func (i *Instance) SetPriority(priority float64) *Instance {
	if i.Aspect == nil {
		i.Aspect = make(map[string]string, 0)
	}
	i.Aspect[priorityKey] = strconv.FormatFloat(priority, 'f', -1, 64)
	return i
}

// GetPriority returns Priority Aspect of an Instance
func (i *Instance) GetPriority() (float64, error) {
	priorStr, ok := i.Aspect[priorityKey]
	if !ok {
		return 0, fmt.Errorf("looks like Priority aspect for instance %v was not defined.. "+
			"Empty output for a key %v", i.Name, priorityKey)
	}
	priority, err := strconv.ParseFloat(priorStr, 64)
	if err != nil {
		return 0, fmt.Errorf("can't parse string %s to float64", priorStr)
	}
	return priority, nil
}

// SetPriority sets Priority Aspect for an Application
func (a *Application) SetPriority(priority float64) *Application {
	if a.Aspect == nil {
		a.Aspect = make(map[string]string, 0)
	}
	a.Aspect[priorityKey] = strconv.FormatFloat(priority, 'f', -1, 64)
	return a
}

// GetPriority returns Priority Aspect of an Application
func (a *Application) GetPriority() (float64, error) {
	priorStr, ok := a.Aspect[priorityKey]
	if !ok {
		return 0, fmt.Errorf("looks like Priority aspect was not defined for the Application.. "+
			"Empty output for a key %v", priorityKey)
	}
	priority, err := strconv.ParseFloat(priorStr, 64)
	if err != nil {
		return 0, fmt.Errorf("can't parse string %s to float64", priorStr)
	}
	return priority, nil
}

// SetReliability sets Reliability Aspect for an Application
func (a *Application) SetReliability(reliability float64) *Application {
	if a.Aspect == nil {
		a.Aspect = make(map[string]string, 0)
	}
	a.Aspect[reliabilityKey] = strconv.FormatFloat(reliability, 'f', -1, 64)
	return a
}

// GetReliability returns Reliability Aspect of an Application
func (a *Application) GetReliability() (float64, error) {
	relStr, ok := a.Aspect[reliabilityKey]
	if !ok {
		return 0, fmt.Errorf("looks like Reliability aspect was not defined for an Application.. "+
			"Empty output for a key %v", reliabilityKey)
	}
	reliability, err := strconv.ParseFloat(relStr, 64)
	if err != nil {
		return 0, fmt.Errorf("can't parse string %s to float64", relStr)
	}
	return reliability, nil
}

// GetInstance returns Instance with a given name from the SystemModel
func (sm *SystemModel) GetInstance(instName string) (*Instance, error) {

	for i := len(sm.Layers); i > 0; i-- {
		layer, ok := sm.Layers[i]
		if !ok {
			return nil, fmt.Errorf("no layer at level %d exists", i)
		}
		for _, v := range layer.Instances {
			// exact matching the name of an instance
			if v.Name == instName {
				return v, nil
			}
		}
	}

	return nil, fmt.Errorf("couldn't find instance with name %s", instName)
}

// SetApplicationPrioritiesRandom sets random priorities for each application
func (sm *SystemModel) SetApplicationPrioritiesRandom() *SystemModel {
	probSum := 1.0
	for _, v := range sm.Applications {
		rnd := rand.Float64() * probSum
		v.SetPriority(rnd)
		probSum -= rnd
	}
	return sm
}

// SetInstancePrioritiesRandom sets random priorities for instances of each application
func (sm *SystemModel) SetInstancePrioritiesRandom() error {
	if len(sm.Layers) != 1 {
		for k, v := range sm.Applications {
			if v.State && !strings.HasPrefix(k, "VI") {
				instCount := v.Rules
				priorSum := 1.0
				// Instances of the application usually sit at the same level,
				// thus it is convenient to avoid redundant iterations and break the loop here
				for i := len(sm.Layers); i > 0 && instCount > 0; i-- {
					layer, ok := sm.Layers[i]
					if !ok {
						return fmt.Errorf("no layer at level %d exists", i)
					}
					for _, inst := range layer.Instances {
						if strings.HasPrefix(inst.Name, k+"-") && inst.Type == App {
							rnd := rand.Float64() * priorSum
							inst.SetPriority(rnd)
							priorSum -= rnd
							instCount--
						}
					}
					if instCount == 0 {
						break
					}
				}
				if instCount != 0 {
					return fmt.Errorf("PANIC!!! Not all instances were found for Application %s!"+
						" %d instances were NOT found", k, instCount)
				}
			} else if v.State { // handling the VI case..
				instCount := int64(*sm.VIcount-1) * int64(v.Rules) // get total amount of VI instances (-1 is to exclude root instance, MAIS)
				priorSum := 1.0
				for i := 1; i <= len(sm.Layers) && instCount > 0; i++ {
					layer, ok := sm.Layers[i]
					if !ok {
						return fmt.Errorf("no layer at level %d exists", i)
					}
					// if there are no VIs, then there is nothing to do
					if !layer.VIwasDeployed {
						break
					}
					for _, inst := range layer.Instances {
						if strings.HasPrefix(inst.Name, "VI") && inst.IsVI() {
							rnd := rand.Float64() * priorSum
							inst.SetPriority(rnd)
							//log.Printf("Setting priority %v to instance %s\n", rnd, inst.Name)
							priorSum -= rnd
							instCount--
						}
					}
				}
				if instCount != 0 {
					return fmt.Errorf("PANIC!!! Not all instances were found for %s!"+
						" %d instances were NOT found", k, instCount)
				}
			}
		}
	} else { // treating the case of a single instance System Model
		if len(sm.Layers) == 1 {
			// drilling to root layer
			if len(sm.Layers[1].Instances) == 1 {
				// double-check
				if strings.HasPrefix(sm.Layers[1].Instances[0].Name, "MAIS") {
					prty := rand.Float64()
					sm.Layers[1].Instances[0].SetPriority(prty)
				}
			}
		}
	}
	//sm.PrettyPrintLayers()

	return nil
}

// SetInstanceReliabilitiesRandom sets random reliabilities for instances of each application
func (sm *SystemModel) SetInstanceReliabilitiesRandom() error {
	if len(sm.Applications) != 0 {
		// checking if at least 1 application was deployed..
		var deployed = 0
		for _, v := range sm.Applications {
			if v.State {
				deployed++
			}
		}
		if deployed != 0 {
			for k, v := range sm.Applications {
				if v.State && !strings.Contains(k, "VI") {
					instCount := v.Rules
					// Instances of the application usually sit at the same level,
					// thus it is convenient to avoid redundant iterations and break the loop here
					for i := len(sm.Layers); i > 0 && instCount > 0; i-- {
						layer, ok := sm.Layers[i]
						if !ok {
							return fmt.Errorf("no layer at level %d exists", i)
						}
						for _, inst := range layer.Instances {
							if strings.HasPrefix(inst.Name, k+"-") && inst.Type == App {
								rnd := rand.Float64()
								inst.SetReliability(rnd)
								instCount--
							}
						}
						if instCount == 0 {
							break
						}
					}
					if instCount != 0 {
						return fmt.Errorf("PANIC!!! Not all instances were found for Application %s!"+
							" %d instances were NOT found", k, instCount)
					}
				} else if strings.Contains(k, "VI") {
					// setting reliabilities only for the VIs which do not deploy any further instances
					for d := len(sm.Layers); d > 0; d-- {
						layer, ok := sm.Layers[d]
						if !ok {
							sm.PrettyPrintApplications().PrettyPrintLayers()
							return fmt.Errorf("couldn't extract layer %d out of system model", d)
						}
						for _, inst := range layer.Instances {
							if strings.Contains(inst.Name, "VI") && len(inst.Relations) == 0 {
								rnd := rand.Float64()
								inst.SetReliability(rnd)
							}
						}
					}
				}
			}
		} else { // if no applications were deployed, setting
			if len(sm.Layers) == 1 {
				// drilling to root layer
				if len(sm.Layers[1].Instances) == 1 {
					// double-check
					if strings.HasPrefix(sm.Layers[1].Instances[0].Name, "MAIS") {
						rlblty := rand.Float64()
						sm.Layers[1].Instances[0].SetReliability(rlblty)
					}
				}
			}
		}
	} else { // treating the case of a single instance System Model
		if len(sm.Layers) == 1 {
			// drilling to root layer
			if len(sm.Layers[1].Instances) == 1 {
				// double-check
				if strings.HasPrefix(sm.Layers[1].Instances[0].Name, "MAIS") {
					rlblty := rand.Float64()
					sm.Layers[1].Instances[0].SetReliability(rlblty)
				}
			}
		}
	}
	return nil
}

// SetChainCoefficient sets Chain Coefficient for an Instance
func (i *Instance) SetChainCoefficient(coef float64) *Instance {
	if i.Aspect == nil {
		i.Aspect = make(map[string]string, 0)
	}
	i.Aspect[chainCoefKey] = strconv.FormatFloat(coef, 'f', -1, 64)
	return i
}

// GetChainCoefficient returns chain coefficient for an Instance
func (i *Instance) GetChainCoefficient() (float64, error) {
	ccStr, ok := i.Aspect[chainCoefKey]
	if !ok {
		return 0, fmt.Errorf("looks like Chain Coefficient aspect was not defined for an Instance.. "+
			"Empty output for a key %v", chainCoefKey)
	}
	chainCoef, err := strconv.ParseFloat(ccStr, 64)
	if err != nil {
		return 0, fmt.Errorf("can't parse string %v to float64", ccStr)
	}
	return chainCoef, nil
}

// SetChainCoefficientForInstance gathers all coefficient in a SystemModel tree on top of the instance
func (sm *SystemModel) SetChainCoefficientForInstance(instName string) error {

	// this is to hold chain coefficient value
	var cc float64 = 1

	targetInst := instName // this is to hold an instance name, which has in its relation a required instance
	// compute chain coefficient
	for i := len(sm.Layers); i > 1; i-- {
		layer, ok := sm.Layers[i]
		if !ok {
			return fmt.Errorf("no layer at level %d exists", i)
		}
		for _, v := range layer.Instances {
			for _, rel := range v.Relations {
				// exact matching the name of an instance
				if rel.Name == targetInst {
					instPriority, err := v.GetPriority()
					if err != nil {
						return err
					}
					cc *= instPriority
					targetInst = v.Name
				}
			}
		}
	}

	// set chain coefficient
	inst, err := sm.GetInstance(instName)
	if err != nil {
		return fmt.Errorf("instance %s was not found in SystemModel: %w", instName, err)
	}
	inst.SetChainCoefficient(cc)

	return nil
}

// SetChainCoefficients sets chain coefficient for all instances, which have no relations
func (sm *SystemModel) SetChainCoefficients() error {

	// setting chain coefficients for instances with no relations
	for i := len(sm.Layers); i > 1; i-- {
		layer, ok := sm.Layers[i]
		if !ok {
			return fmt.Errorf("no layer at level %d exists", i)
		}
		for _, v := range layer.Instances {
			if len(v.Relations) == 0 {
				err := sm.SetChainCoefficientForInstance(v.Name)
				if err != nil {
					return err
				}
			}
		}
	}

	// setting chain coefficients to corresponding Applications
	for k, v := range sm.Applications {
		if !strings.HasPrefix(k, "VI") {
			// treating an application case
			for i := len(sm.Layers); i > 0; i-- {
				flag := false
				layer, ok := sm.Layers[i]
				if !ok {
					return fmt.Errorf("no layer at level %d exists", i)
				}
				for _, val := range layer.Instances {
					// exact matching the name of an instance
					if strings.HasPrefix(val.Name, k) {
						// obtaining a chain coefficient for an instance (they should be the same for all instance within the same Application)
						cc, err := val.GetChainCoefficient()
						if err != nil {
							return err
						}
						// setting chain coefficient to the Application
						v.SetChainCoefficient(cc)
						flag = true
						break
					}
				}
				// if instance is found, then breaking an algorithm
				if flag {
					break
				}
			}
		}
	}

	return nil
}

// SetChainCoefficient sets Chain Coefficient for an Application
func (a *Application) SetChainCoefficient(coef float64) *Application {
	if a.Aspect == nil {
		a.Aspect = make(map[string]string, 0)
	}
	a.Aspect[chainCoefKey] = strconv.FormatFloat(coef, 'f', -1, 64)
	return a
}

// GetChainCoefficient returns chain coefficient for an Application
func (a *Application) GetChainCoefficient() (float64, error) {
	ccStr, ok := a.Aspect[chainCoefKey]
	if !ok {
		return 0, fmt.Errorf("looks like Chain Coefficient aspect was not defined for an Application.. "+
			"Empty output for a key %v", chainCoefKey)
	}
	chainCoef, err := strconv.ParseFloat(ccStr, 64)
	if err != nil {
		return 0, fmt.Errorf("can't parse string %v to float64", ccStr)
	}
	return chainCoef, nil
}
