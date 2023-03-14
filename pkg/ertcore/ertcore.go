// Package ertcore implements an ERT-CORE reliability model
package ertcore

import (
	"fmt"
)

// ErtCore defines an ERT-CORE instance reliability
type ErtCore struct {
	ReliabilityPerParameter ReliabilityPerParameter // map key is a parameter name
	Priorities              map[string]float64      // map key is a parameter name
	Reliability             float64                 // holds computed reliability
	EntityName              string                  // holds an entity name
}

// ReliabilityPerParameter defines reliability value per parameter
type ReliabilityPerParameter struct {
	InputMetrics            map[string]*InputMetric // map key is a parameter name
	Priorities              map[string]Priorities   // map key is a parameter name
	ReliabilityPerParameter map[string]float64      // holds computed reliability per parameter vector, key is a parameter name
}

// Priorities contains a vector of First-level priorities, which hold priority values (value) per component (key)
type Priorities map[string]float64

// InputMetric defines input metric per parameter
type InputMetric struct {
	InputData                 map[string]float64 // map key is a component name
	SLA                       map[string]float64 // map key is a component name
	InputMetrics              map[string]float64 // holds computed Input Metrics vector values
	ValuesRangeGreaterThanSLA bool               // indicates if value range of the input data is greater than SLA
}

// Initialize initializes an InstanceReliability structure
func (r *ErtCore) Initialize() *ErtCore {
	r.Priorities = make(map[string]float64, 0)
	r.Reliability = 0.0
	return r
}

// SetName sets a name of an ERT-CORE instance
func (r *ErtCore) SetName(name string) *ErtCore {
	r.EntityName = name
	return r
}

// ComputeReliability computes an instance Reliability value
func (r *ErtCore) ComputeReliability() (float64, error) {
	var reliability float64

	for k, v := range r.Priorities {
		r.ReliabilityPerParameter.InputMetrics[k].ComputeInputMetric()
		rpp, err := r.ReliabilityPerParameter.ComputeReliabilityPerGivenParameter(k)
		if err != nil {
			return reliability, err
		}
		reliability += rpp * v
	}
	r.Reliability = reliability
	return reliability, nil
}

// SetReliabilityPerParameter sets a reliability per parameter
func (r *ErtCore) SetReliabilityPerParameter(rpp ReliabilityPerParameter) *ErtCore {
	r.ReliabilityPerParameter = rpp
	return r
}

// SetPriorities sets a priority vector as a map, where key is a parameter name and a value is a priority (= weight)
// of reliability per parameter
func (r *ErtCore) SetPriorities(priorities map[string]float64) *ErtCore {
	r.Priorities = priorities
	return r
}

// Initialize initializes a ReliabilityPerParameter structure
func (rp *ReliabilityPerParameter) Initialize() *ReliabilityPerParameter {
	rp.InputMetrics = make(map[string]*InputMetric, 0)
	rp.Priorities = make(map[string]Priorities, 0)
	rp.ReliabilityPerParameter = make(map[string]float64, 0)
	return rp
}

// ComputeReliabilityPerParameter computes a reliability per parameter value
func (rp *ReliabilityPerParameter) ComputeReliabilityPerParameter() *ReliabilityPerParameter {
	for k, v := range rp.InputMetrics {
		reliability := 0.0
		for c, pr := range rp.Priorities[k] {
			reliability += v.InputMetrics[c] * pr
		}
		rp.ReliabilityPerParameter[k] = reliability
	}
	return rp
}

// ComputeReliabilityPerGivenParameter computes a reliability for provided parameter value
func (rp *ReliabilityPerParameter) ComputeReliabilityPerGivenParameter(param string) (float64, error) {
	reliability := 0.0
	im, ok := rp.InputMetrics[param]
	if !ok {
		return 0.0, fmt.Errorf("map entry %v does not exist in ap %v", param, rp.InputMetrics)
	}
	for c, pr := range rp.Priorities[param] {
		reliability += im.InputMetrics[c] * pr
	}
	rp.ReliabilityPerParameter[param] = reliability
	return reliability, nil
}

// SetInputMetrics sets input metrics as a map, where key is a parameter name and a value is an input metric vector
func (rp *ReliabilityPerParameter) SetInputMetrics(inmetrics map[string]*InputMetric) *ReliabilityPerParameter {
	rp.InputMetrics = inmetrics
	return rp
}

// SetPriorities sets a priorities vector as a map, where key is a parameter name and a value is a priority (per component) vector
func (rp *ReliabilityPerParameter) SetPriorities(priorities map[string]Priorities) *ReliabilityPerParameter {
	rp.Priorities = priorities
	return rp
}

// Initialize initializes an InputMetric structure
func (im *InputMetric) Initialize() *InputMetric {
	im.InputData = make(map[string]float64, 0)
	im.SLA = make(map[string]float64, 0)
	im.InputMetrics = make(map[string]float64, 0)
	im.ValuesRangeGreaterThanSLA = false // by default we are assuming that input data values are in range (0, 1)
	return im
}

// SetValuesAreGreaterThanSLA sets a flag that indicates that the input data values are typically in range (0, +inf).
func (im *InputMetric) SetValuesAreGreaterThanSLA() *InputMetric {
	im.ValuesRangeGreaterThanSLA = true
	return im
}

// ComputeInputMetric computes input metric vector for the given input data and SLA
func (im *InputMetric) ComputeInputMetric() {

	for k, v := range im.InputData {
		if im.ValuesRangeGreaterThanSLA {
			im.InputMetrics[k] = im.SLA[k] / v
		} else {
			im.InputMetrics[k] = v / im.SLA[k]
		}
	}
}

// SetInputData sets an input data vector as a map, where key is the component name, value is a SLA value
func (im *InputMetric) SetInputData(indata map[string]float64) *InputMetric {
	im.InputData = indata
	return im
}

// SetSLA sets SLA vector as a map, where key is the component name, value is a SLA value
func (im *InputMetric) SetSLA(sla map[string]float64) *InputMetric {
	im.SLA = sla
	return im
}

// CreateInputDataVector creates a map entry to pass as an input to SetInputData function of InputMetric structure
func CreateInputDataVector(components []string, indata []float64) map[string]float64 {
	res := make(map[string]float64, 0)
	for i, val := range components {
		res[val] = indata[i]
	}
	return res
}

// CreateSLAVector creates a map entry to pass as an input parameter to SetSLA function of InputMetrics structure
func CreateSLAVector(components []string, sla []float64) map[string]float64 {
	res := make(map[string]float64, 0)
	for i, val := range components {
		res[val] = sla[i]
	}
	return res
}

// CreateInputMetricsVectorForReliabilityPerParameter creates a map entry to pass as an input to SetInputMetrics function of ReliabilityPerParameter structure
func CreateInputMetricsVectorForReliabilityPerParameter(parameters []string, v ...InputMetric) map[string]*InputMetric {
	res := make(map[string]*InputMetric, 0)
	for i, val := range parameters {
		res[val] = &v[i]
	}
	return res
}

// CreateReliabilityPerParameterPriorities creates a map entry to pass as an input to SetPriorities function of ReliabilityPerParameter structure
func CreateReliabilityPerParameterPriorities(parameters []string, priorities ...Priorities) map[string]Priorities {
	res := make(map[string]Priorities, 0)
	for i, val := range parameters {
		res[val] = priorities[i]
	}
	return res
}

// CreatePrioritiesVectorForInstanceReliability creates a map entry to pass as an input to SetPriorities function of InstanceReliability structure
func CreatePrioritiesVectorForInstanceReliability(parameters []string, priorities []float64) map[string]float64 {
	res := make(map[string]float64, 0)
	for i, val := range parameters {
		res[val] = priorities[i]
	}
	return res
}
