// Package meertcore implements ME-ERT-CORE reliability model. This file in particular implements optimized version of ME-ERT-CORE.
package meertcore

import (
	"fmt"
	"strings"
)

// ComputeReliabilityOptimized function implements optimized version of a ME-ERT-CORE reliability computation
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
