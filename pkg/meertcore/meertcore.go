// Package meertcore implements ME-ERT-CORE reliability model. This package implements canonical version of ME-ERT-CORE.
package meertcore

import "gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/pkg/systemmodel"

// MeErtCore structure represents an ME-ERT-CORE instance reliability
type MeErtCore struct {
	SystemModel *systemmodel.SystemModel // holds System Model definition
	Reliability float64                  // contains Reliability of the System Model
}
