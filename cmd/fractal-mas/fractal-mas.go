// main package structures experiment which measures time complexity of the algorithms
package main

import (
	"fmt"
	"gitlab.fel.cvut.cz/eroshiva/fractal-multi-agent-system/pkg/systemmodel"
	"time"
)

func main() {
	start := time.Now()
	fmt.Println("Hello, World!")
	duration := time.Since(start)
	fmt.Printf("It took %d us to print previous message\n", duration.Microseconds())
	sm := systemmodel.SystemModel{}
	sm.InitializeSystemModel(20, 5)
}
