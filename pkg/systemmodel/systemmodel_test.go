package systemmodel

import (
	"gotest.tools/assert"
	"testing"
)

func TestCreateSystemModel(t *testing.T) {

	// defining input parameters
	var a = 10
	var l = 4

	systemModel := &SystemModel{}
	systemModel.InitializeSystemModel(a, l)
	t.Logf("System model is\n%v", systemModel)
	assert.Equal(t, systemModel.Depth, l)

	// defining a maximum number of instances per application
	var maxNumInstances = 15 // maximum 15 instances per app
	var minNumInstances = 1  // minimum 1 instances per app
	// defining list of application names
	names := GenerateAppNames(a)
	t.Logf("Generated Application names are:\n%v\n", names)
	systemModel.CreateRandomApplications(names, minNumInstances, maxNumInstances)
	t.Logf("System model is\n%v", systemModel)
	assert.Equal(t, len(names), len(systemModel.Applications))
	t.Logf("VI is\n%v", systemModel.Applications["VI"])
	// check that the probabilities are of total 1
	sum := float32(0)
	for k, v := range systemModel.Applications {
		sum += v.Probability
		t.Logf("%s has probability %v and deploys %v applications\n", k, v.Probability, v.Rules)
	}
	t.Logf("Total probabilities are %v\n", sum)
	assert.Assert(t, sum <= float32(1.0001))
}

func TestGenerateRandomSystemModel(t *testing.T) {
	systemModel := &SystemModel{}
	// defining input parameters
	var a = 10
	var l = 4
	// defining a maximum number of instances per application
	var maxNumInstances = 15 // maximum 15 instances per app
	var minNumInstances = 1  // minimum 1 instances per app
	// defining list of application names
	names := GenerateAppNames(a)
	systemModel.InitializeSystemModel(maxNumInstances, l)
	systemModel.CreateRandomApplications(names, minNumInstances, maxNumInstances)
	systemModel.GenerateSystemModel()
	t.Logf("System model is\n%v", systemModel)

	// check that the probabilities are of total 1
	sum := float32(0)
	for k, v := range systemModel.Applications {
		sum += v.Probability
		t.Logf("%s has probability %v and deploys %v applications\n", k, v.Probability, v.Rules)
	}
	t.Logf("Total probabilities are %v\n", sum)
	assert.Assert(t, sum <= float32(1.0001)) // leaving .0001 as a possible overhead due to float32 operations..
}

func BenchmarkGenerateSystemModel(b *testing.B) {
	systemModel := &SystemModel{}
	// defining input parameters
	var a = 10
	var l = 4
	// defining a minimum and maximum number of instances per application
	var maxNumInstances = 15 // maximum 15 instances per app
	var minNumInstances = 1  // minimum 1 instances per app
	// defining list of application names
	names := GenerateAppNames(a)
	systemModel.InitializeSystemModel(maxNumInstances, l)
	systemModel.CreateRandomApplications(names, minNumInstances, maxNumInstances)
	systemModel.GenerateSystemModel()
	//b.Logf("System model is\n%v", systemModel)
}
