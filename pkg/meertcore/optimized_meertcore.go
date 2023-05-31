// Package meertcore implements ME-ERT-CORE reliability model. This file in particular implements optimized version of ME-ERT-CORE.
package meertcore

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// ComputeReliabilityOptimized function implements a generic optimized version of a ME-ERT-CORE reliability computation
// Total Reliability of System Model equals to weighted sum of each application reliability and weighted sm of all
// reliabilities at the last layer of System Model
func (me *MeErtCore) ComputeReliabilityOptimized() (float64, error) {
	var reliability float64
	for k, v := range me.SystemModel.Applications {
		if v.State && !strings.HasPrefix(k, "VI") {
			rlblty, err := v.GetReliability()
			if err != nil {
				return 0, fmt.Errorf("application %s: %w", k, err)
			}
			cc, err := v.GetChainCoefficient()
			if err != nil {
				return 0, fmt.Errorf("application %s: %w", k, err)
			}
			reliability += rlblty * cc
		} else if strings.HasPrefix(k, "VI") {
			// gather reliability of all VIs, which do not deploy any further instance
			var viRel float64
			for d := len(me.SystemModel.Layers); d > 0; d-- {
				for _, val := range me.SystemModel.Layers[d].Instances {
					if val.IsVI() && len(val.Relations) == 0 {
						priority, err := val.GetPriority()
						if err != nil {
							return 0, fmt.Errorf("application %s: %w", val.Name, err)
						}
						rlblty, err := val.GetReliability()
						if err != nil {
							return 0, fmt.Errorf("application %s: %w", val.Name, err)
						}
						cc, err := val.GetChainCoefficient()
						if err != nil {
							return 0, fmt.Errorf("application %s: %w", val.Name, err)
						}
						viRel += rlblty * priority * cc
					}
				}
			}
			me.SystemModel.Applications[k].SetReliability(viRel)
			reliability += viRel
		}
	}
	me.Reliability = reliability
	return reliability, nil
}

// ComputeReliabilityOptimizedSimple function implements simplified and optimized version of a ME-ERT-CORE reliability
// computation. Simplification consists in assuming that Chain Coefficient (CC) is 1 for all VIs.
// Total Reliability of System Model equals to weighted sum of each application reliability and weighted sm of all
// reliabilities at the last layer of System Model.
//
// This function does not produce the same output as ComputeReliabilityPerDefinition(). This is because the System Model
// should be set a bit differently - priority values for VIs are distributed differently. This is the main cause of
// different results.
func (me *MeErtCore) ComputeReliabilityOptimizedSimple() (float64, error) {
	var reliability float64
	for k, v := range me.SystemModel.Applications {
		if v.State && !strings.HasPrefix(k, "VI") {
			rlblty, err := v.GetReliability()
			if err != nil {
				return 0, fmt.Errorf("application %s: %w", k, err)
			}
			priority, err := v.GetPriority()
			if err != nil {
				return 0, fmt.Errorf("application %s: %w", k, err)
			}
			reliability += rlblty * priority
			me.SystemModel.Applications[k].SetReliability(rlblty * priority)
		} else if strings.HasPrefix(k, "VI") {
			// gather reliability of all VIs, which do not deploy any further instance
			var viRel float64
			viPriority, err := me.SystemModel.Applications["VI"].GetPriority()
			if err != nil {
				return 0, err
			}
			for d := len(me.SystemModel.Layers); d > 0; d-- {
				for _, val := range me.SystemModel.Layers[d].Instances {
					if val.IsVI() && len(val.Relations) == 0 {
						priority, err := val.GetPriority()
						if err != nil {
							return 0, fmt.Errorf("application %s: %w", val.Name, err)
						}
						rlblty, err := val.GetReliability()
						if err != nil {
							return 0, fmt.Errorf("application %s: %w", val.Name, err)
						}
						viRel += rlblty * priority * viPriority
					}
				}
			}
			reliability += viRel
		}
	}
	me.Reliability = reliability
	return reliability, nil
}

// ComputeMeErtCoreCoefficient computes ME-ERT-CORE coefficient, which describes better obtained reliability values.
// It requires on input to get a reliability value and a System Model depth
func ComputeMeErtCoreCoefficient(rel float64, depth int) (float64, error) {
	scaledCoef := math.Pow10(depth-1) * rel
	return getDecimal(scaledCoef)
}

// getDecimal function gets decimal part of the number (numbers after point)
func getDecimal(val float64) (float64, error) {
	decimalPart := 0.0
	s := fmt.Sprintf("%v", val)
	split := strings.Split(s, ".")
	decimalStr := "0." + split[1]
	decimalPart, err := strconv.ParseFloat(decimalStr, 64)
	if err != nil {
		return 0, fmt.Errorf("something went wrong in float64 parsing from string: %w", err)
	}
	return decimalPart, nil
}
