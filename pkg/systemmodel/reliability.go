// Package systemmodel implements means of Fractal MAS system model. This file in particular holds all functions
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

// GatherApplicationInstanceReliabilities gathers reliability for all entities of an application.
// This function does not work for VI - it returns empty map and nil error
func (sm *SystemModel) GatherApplicationInstanceReliabilities(appName string) (map[string]float64, error) {
	app, ok := sm.Applications[appName]
	if !ok {
		sm.PrettyPrintApplications()
		return nil, fmt.Errorf("application %s was not initialized in a SystemModel", appName)
	}
	if app.State && !strings.HasPrefix(appName, "VI") {
		res := make(map[string]float64, app.Rules)
		var appReliability float64
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
				}
			}
		}
		sm.Applications[appName].SetReliability(appReliability)
		return res, nil
	} else if strings.HasPrefix(appName, "VI") { // do not return anything for VI
		return nil, nil
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
				return nil, fmt.Errorf("application %s: %v", k, err)
			}
			rel, err := v.GetReliability()
			if err != nil {
				return nil, fmt.Errorf("application %s: %v", k, err)
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
			"Empty output for a key %v", i.Name, reliabilityKey)
	}
	reliability, err := strconv.ParseFloat(relStr, 64)
	if err != nil {
		return 0, fmt.Errorf("can't parse string %v to float64", reliability)
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
		return 0, fmt.Errorf("can't parse string %v to float64", priority)
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
		return 0, fmt.Errorf("looks like Priority aspect was not defined for the application.. "+
			"Empty output for a key %v", priorityKey)
	}
	priority, err := strconv.ParseFloat(priorStr, 64)
	if err != nil {
		return 0, fmt.Errorf("can't parse string %v to float64", priority)
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
		return 0, fmt.Errorf("can't parse string %v to float64", reliability)
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
	for k, v := range sm.Applications {
		if v.State && !strings.Contains(k, "VI") {
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
			instCount := *sm.VIcount
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
					if strings.HasPrefix(inst.Name, "VI") && inst.Type == VI {
						rnd := rand.Float64() * priorSum
						inst.SetPriority(rnd)
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
	return nil
}

// SetInstanceReliabilitiesRandom sets random reliabilities for instances of each application
func (sm *SystemModel) SetInstanceReliabilitiesRandom() error {
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
			for _, inst := range sm.Layers[len(sm.Layers)].Instances {
				if strings.Contains(inst.Name, "VI") {
					rnd := rand.Float64()
					inst.SetReliability(rnd)
				}
			}
		}
	}
	return nil
}
