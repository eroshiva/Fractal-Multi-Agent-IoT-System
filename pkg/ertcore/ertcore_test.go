package ertcore

import (
	"gotest.tools/assert"
	"testing"
)

func TestComputeReliability1(t *testing.T) {
	// setting input data for Workload
	w := InputMetric{}
	w.Initialize().SetInputData(map[string]float64{
		"CPU":     0.46,
		"RAM":     0.7,
		"NI":      0.27,
		"Storage": 0.46,
	}).SetSLA(map[string]float64{
		"CPU":     1.0,
		"RAM":     1.0,
		"NI":      1.0,
		"Storage": 1.0,
	}).ComputeInputMetric()
	workloadPrior := Priorities{
		"CPU":     0.6,
		"RAM":     0.01,
		"NI":      0.25,
		"Storage": 0.09,
	}

	// setting ARPT input data
	arpt := InputMetric{}
	arpt.Initialize().SetInputData(map[string]float64{
		"CPU":     0.3,
		"RAM":     0.009,
		"NI":      0.07,
		"Storage": 0.44,
	}).SetSLA(map[string]float64{
		"CPU":     0.5,
		"RAM":     0.015,
		"NI":      0.15,
		"Storage": 0.6,
	}).ComputeInputMetric()
	arptPrior := Priorities{
		"CPU":     0.55,
		"RAM":     0.01,
		"NI":      0.34,
		"Storage": 0.1,
	}

	// setting Availability input data
	av := InputMetric{}
	av.Initialize().SetInputData(map[string]float64{
		"CPU":     0.9,
		"RAM":     0.85,
		"NI":      0.7,
		"Storage": 0.84,
	}).SetSLA(map[string]float64{
		"CPU":     0.95,
		"RAM":     0.8,
		"NI":      0.5,
		"Storage": 0.84,
	}).SetValuesAreGreaterThanSLA().ComputeInputMetric()
	avPrior := Priorities{
		"CPU":     0.45,
		"RAM":     0.09,
		"NI":      0.1,
		"Storage": 0.36,
	}

	// setting a Reliability per Parameter vector
	rpp := ReliabilityPerParameter{}
	rpp.Initialize().SetPriorities(map[string]Priorities{
		"Workload":     workloadPrior,
		"ARPT":         arptPrior,
		"Availability": avPrior,
	}).SetInputMetrics(map[string]*InputMetric{
		"Workload":     &w,
		"ARPT":         &arpt,
		"Availability": &av,
	}).ComputeReliabilityPerParameter()

	r := InstanceReliability{}
	r.Initialize().SetReliabilityPerParameter(rpp).SetPriorities(map[string]float64{
		"Workload":     0.4,
		"ARPT":         0.55,
		"Availability": 0.05,
	})
	err := r.ComputeInstanceReliability()
	assert.NilError(t, err)
	assert.Assert(t, r.Reliability > 0.0)
	t.Logf("Computed reliability is %v\n", r.Reliability)
}

// This is an example on how to save some typing..
func TestComputeReliability2(t *testing.T) {
	// defining components
	c := []string{"CPU", "RAM", "NI", "Storage"}
	// defining parameters
	params := []string{"Workload", "ARPT", "Availability"}

	// setting input data for Workload
	w := InputMetric{}
	w.Initialize().SetInputData(CreateInputDataVector(c, []float64{0.46, 0.7, 0.27, 0.46})).
		SetSLA(CreateSLAVector(c, []float64{1.0, 1.0, 1.0, 1.0})).ComputeInputMetric()
	workloadPrior := Priorities{
		"CPU":     0.6,
		"RAM":     0.01,
		"NI":      0.25,
		"Storage": 0.09,
	}

	// setting ARPT input data
	arpt := InputMetric{}
	arpt.Initialize().SetInputData(CreateInputDataVector(c, []float64{0.3, 0.009, 0.07, 0.44})).
		SetSLA(CreateSLAVector(c, []float64{0.5, 0.015, 0.15, 0.6})).ComputeInputMetric()
	arptPrior := Priorities{
		"CPU":     0.55,
		"RAM":     0.01,
		"NI":      0.34,
		"Storage": 0.1,
	}

	// setting Availability input data
	av := InputMetric{}
	av.Initialize().SetInputData(CreateInputDataVector(c, []float64{0.9, 0.85, 0.7, 0.84})).
		SetSLA(CreateSLAVector(c, []float64{0.95, 0.8, 0.5, 0.84})).SetValuesAreGreaterThanSLA().ComputeInputMetric()
	avPrior := Priorities{
		"CPU":     0.45,
		"RAM":     0.09,
		"NI":      0.1,
		"Storage": 0.36,
	}

	// setting a Reliability per Parameter vector
	rpp := ReliabilityPerParameter{}
	rpp.Initialize().SetPriorities(CreateReliabilityPerParameterPriorities(params, workloadPrior, arptPrior, avPrior)).
		SetInputMetrics(CreateInputMetricsVectorForReliabilityPerParameter(params, w, arpt, av)).
		ComputeReliabilityPerParameter()

	r := InstanceReliability{}
	r.Initialize().SetReliabilityPerParameter(rpp).
		SetPriorities(CreatePrioritiesVectorForInstanceReliability(params, []float64{0.4, 0.55, 0.05}))
	err := r.ComputeInstanceReliability()
	assert.NilError(t, err)
	assert.Assert(t, r.Reliability > 0.0)
	t.Logf("Computed reliability is %v\n", r.Reliability)
}
