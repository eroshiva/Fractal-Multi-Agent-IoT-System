// Package meertcore implements ME-ERT-CORE reliability model. This package implements canonical version of ME-ERT-CORE.
package meertcore

import (
	"fmt"
	"gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/pkg/systemmodel"
)

// MeErtCore structure represents an ME-ERT-CORE instance reliability
type MeErtCore struct {
	SystemModel *systemmodel.SystemModel // holds System Model definition
	Reliability float64                  // contains Reliability of the System Model
}

// ComputeReliabilityPerDefinition computes reliability of Fractal MAIS (i.e., System Model), per canonical definition
func (me *MeErtCore) ComputeReliabilityPerDefinition() (float64, error) {
	// no need to iterate over the last layer - reliabilities of all instances should be present for our disposal
	for d := len(me.SystemModel.Layers) - 1; d > 0; d-- {
		layer, ok := me.SystemModel.Layers[d]
		if !ok {
			me.SystemModel.PrettyPrintApplications().PrettyPrintLayers()
			return 0.0, fmt.Errorf("couldn't extract layer %d out of the System Model", d)
		}
		for _, inst := range layer.Instances {
			var instRel float64 // there would be resulting reliability of an instance
			// if we are on the last layer,
			// there is nothing to do - reliabilities should all be set.
			// if we are not on the last layer, we are
			// iterating over the instance relations and computing reliability of an instance
			if len(inst.Relations) != 0 {
				for _, rel := range inst.Relations {
					reliability, err := rel.GetReliability()
					if err != nil {
						return 0, err
					}
					priority, err := rel.GetPriority()
					if err != nil {
						return 0, err
					}

					if rel.IsApp() {
						// extracting coefficient of an Application (i.e., priority)
						appName, err := rel.GetAppName()
						if err != nil {
							return 0, err
						}
						app, okie := me.SystemModel.Applications[appName]
						if !okie {
							// try one more time
							appNameOld, err := rel.GetAppName()
							if err != nil {
								return 0, err
							}
							app, okie2 := me.SystemModel.Applications[appNameOld]
							if !okie2 {
								return 0, fmt.Errorf("couldn't extract application with a key %s", appName)
							}
							appInstancePriority, err := app.GetPriority()
							if err != nil {
								return 0, err
							}

							instRel += reliability * priority * appInstancePriority
						} else {
							appInstancePriority, err := app.GetPriority()
							if err != nil {
								return 0, err
							}

							instRel += reliability * priority * appInstancePriority
						}
					} else { // Treating the VI case
						app, okie := me.SystemModel.Applications["VI"]
						if !okie {
							return 0, fmt.Errorf("couldn't extract VI from an application dictionary")
						}
						viPriority, err := app.GetPriority()
						if err != nil {
							return 0, err
						}

						instRel += reliability * priority * viPriority
					}
				}
				// setting computed reliability to the instance
				inst.SetReliability(instRel)
			}
		}
	}

	// getting total reliability of the System Model - at the layer 1 there is only one instance, i.e., MAIS!
	totalReliability, err := me.SystemModel.Layers[1].Instances[0].GetReliability()
	if err != nil {
		me.SystemModel.PrettyPrintApplications()
		me.SystemModel.PrettyPrintLayers()
		return 0, err
	}
	me.Reliability = totalReliability

	return totalReliability, nil
}
