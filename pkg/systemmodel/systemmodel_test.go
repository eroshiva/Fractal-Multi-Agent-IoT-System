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
	var maxNumInstances = 15 // 15 instances per app
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
	var maxNumInstances = 15 // 15 instances per app
	// defining list of application names
	names := []string{"VI", "App#1", "App#2", "App#3", "App#4", "App#5", "App#6", "App#7",
		"App#8", "App#9"}
	systemModel.GenerateRandomSystemModel(l, names, maxNumInstances)
	t.Logf("System model is\n%v", systemModel)

	// check that the probabilities are of total 1
	sum := float32(0)
	for k, v := range systemModel.Applications {
		sum += v.Probability
		t.Logf("%s has probability %v and deploys %v applications\n", k, v.Probability, v.Rules)
	}
	t.Logf("Total probabilities are %v\n", sum)
	assert.Assert(t, sum <= float32(1))

	for k, v := range systemModel.Layers {
		t.Logf("Layer %v has %v deplyed applications\n%v", k, len(v.DeployedApps), v)
	}
}

func BenchmarkGenerateRandomSystemModel(b *testing.B) {
	systemModel := &SystemModel{}

	// defining input parameters
	var l int32 = 4
	// defining a maximum number of instances per application
	var maxNumInstances = 15 // 15 instances per app
	// defining list of application names
	names := []string{"VI", "App#1", "App#2", "App#3", "App#4", "App#5", "App#6", "App#7",
		"App#8", "App#9"}
	systemModel.GenerateRandomSystemModel(l, names, maxNumInstances)
	//b.Logf("System model is\n%v", systemModel)
}
