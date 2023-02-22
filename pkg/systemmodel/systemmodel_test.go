package systemmodel

import (
	"gotest.tools/assert"
	"testing"
)

func TestCreateSystemModel(t *testing.T) {

	// defining input parameters
	var a int32 = 10
	var l int32 = 4

	systemModel := &SystemModel{}
	systemModel.InitializeSystemModel(a, l)
	t.Logf("System model is\n%v", systemModel)
	assert.Equal(t, systemModel.Depth, l)

	// defining a maximum number of instances per application
	var maxNumInstances int32 = 15 // 15 instances per app
	// defining list of application names
	names := []string{"VI", "App#1", "App#2", "App#3", "App#4", "App#5", "App#6", "App#7",
		"App#8", "App#9"}
	systemModel.CreateRandomApplications(names, maxNumInstances)
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
	assert.Assert(t, sum <= float32(1))
}

func TestGenerateRandomSystemModel(t *testing.T) {
	systemModel := &SystemModel{}
	// defining input parameters
	var l int32 = 4
	// defining a maximum number of instances per application
	var maxNumInstances int32 = 15 // 15 instances per app
	// defining list of application names
	names := []string{"VI", "App#1", "App#2", "App#3", "App#4", "App#5", "App#6", "App#7",
		"App#8", "App#9"}
	systemModel.InitializeSystemModel(maxNumInstances, l)
	systemModel.CreateRandomApplications(names, maxNumInstances)
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

	t.Logf("----------------------- Printing the output result of MAIS ----------------------\n")
	systemModel.PrettyPrintApplications()
	systemModel.PrettyPrintLayers()
}

func BenchmarkGenerateSystemModel(b *testing.B) {
	systemModel := &SystemModel{}
	// defining input parameters
	var l int32 = 4
	// defining a maximum number of instances per application
	var maxNumInstances int32 = 15 // 15 instances per app
	// defining list of application names
	names := []string{"VI", "App#1", "App#2", "App#3", "App#4", "App#5", "App#6", "App#7",
		"App#8", "App#9"}
	systemModel.InitializeSystemModel(maxNumInstances, l)
	systemModel.CreateRandomApplications(names, maxNumInstances)
	systemModel.GenerateSystemModel()
	//b.Logf("System model is\n%v", systemModel)
}
