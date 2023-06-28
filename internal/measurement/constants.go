// Package measurement provides a measurement logic and all helper functions. This file in particular
// specifies constants/input data used in the measurement
package measurement

// deviation defines a deviation of a Normal distribution
const (
	app1Name    = "App#1"
	app2Name    = "App#2"
	app3Name    = "App#3"
	app4Name    = "App#4"
	appFailName = "App#1" // all other apps won't fail in large-scale measurement
)

// this variable is adjusted in each measurement scenario to match a VIaaS (one per system)
var viName = "VI"

var deviation = 0.025

// app1inst1 defines reliabilities values for Instance #1 of the Application #1
var app1inst1 = []inputData{
	{
		from:  1,
		to:    50,
		value: 0.5,
	},
	{
		from:  51,
		to:    75,
		value: 0.35,
	},
	{
		from:  76,
		to:    300,
		value: 0.5,
	},
}

// app1inst2 defines reliabilities values for Instance #2 of the Application #1
var app1inst2 = []inputData{
	{
		from:  1,
		to:    25,
		value: 0.37,
	},
	{
		from:  26,
		to:    40,
		value: 0.27,
	},
	{
		from:  41,
		to:    300,
		value: 0.37,
	},
}

// app1inst3 defines reliabilities values for Instance #3 of the Application #1
var app1inst3 = []inputData{
	{
		from:  1,
		to:    100,
		value: 0.71,
	},
	{
		from:  101,
		to:    125,
		value: 0.59,
	},
	{
		from:  126,
		to:    300,
		value: 0.71,
	},
}

// app2inst1 defines reliabilities values for Instance #1 of the Application #2
var app2inst1 = []inputData{
	{
		from:  1,
		to:    130,
		value: 0.46,
	},
	{
		from:  131,
		to:    150,
		value: 0.24,
	},
	{
		from:  151,
		to:    300,
		value: 0.46,
	},
}

// app2inst2 defines reliabilities values for Instance #2 of the Application #2
var app2inst2 = []inputData{
	{
		from:  1,
		to:    160,
		value: 0.69,
	},
	{
		from:  161,
		to:    175,
		value: 0.0,
	},
	{
		from:  176,
		to:    300,
		value: 0.9,
	},
}

// app3inst1 defines reliabilities values for Instance #1 of the Application #3
var app3inst1 = []inputData{
	{
		from:  1,
		to:    190,
		value: 0.54,
	},
	{
		from:  191,
		to:    215,
		value: 0.38,
	},
	{
		from:  216,
		to:    300,
		value: 0.54,
	},
}

// app3inst2 defines reliabilities values for Instance #2 of the Application #3
var app3inst2 = []inputData{
	{
		from:  1,
		to:    230,
		value: 0.47,
	},
	{
		from:  231,
		to:    245,
		value: 0.33,
	},
	{
		from:  246,
		to:    300,
		value: 0.47,
	},
}

// app4inst1 defines reliabilities values for Instance #1 of the Application #4
var app4inst1 = []inputData{
	{
		from:  1,
		to:    250,
		value: 0.8,
	},
	{
		from:  251,
		to:    261,
		value: 0.0,
	},
	{
		from:  262,
		to:    300,
		value: 0.8,
	},
}

// viaas defines reliabilities values for VIaaS
var viaas = []inputData{
	{
		from:  1,
		to:    270,
		value: 0.94,
	},
	{
		from:  271,
		to:    285,
		value: 0.21,
	},
	{
		from:  286,
		to:    300,
		value: 0.94,
	},
}

// appInst1 defines reliabilities values for Instance #1 of the Application (which does NOT fails)
var appInst1 = []inputData{
	{
		from:  1,
		to:    300,
		value: 0.46,
	},
}

// appInst2 defines reliabilities values for Instance #2 of the Application (which does NOT fail)
var appInst2 = []inputData{
	{
		from:  1,
		to:    300,
		value: 0.69,
	},
}

// appFailInst1 defines reliabilities values for Instance #1 of the Application (which FAILS)
var appFailInst1 = []inputData{
	{
		from:  1,
		to:    99,
		value: 0.46,
	},
	{
		from:  100,
		to:    130,
		value: 0.3,
	},
	{
		from:  131,
		to:    300,
		value: 0.46,
	},
}

// appFailInst2 defines reliabilities values for Instance #2 of the Application (which FAILS)
var appFailInst2 = []inputData{
	{
		from:  1,
		to:    160,
		value: 0.69,
	},
	{
		from:  161,
		to:    185,
		value: 0.59,
	},
	{
		from:  186,
		to:    300,
		value: 0.69,
	},
}
