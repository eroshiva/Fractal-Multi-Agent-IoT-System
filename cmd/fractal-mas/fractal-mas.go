// main package structures experiment which measures time complexity of the algorithms
package main

import (
	"fmt"
	"time"
)

func main() {
	start := time.Now()
	fmt.Println("Hello, World!")
	duration := time.Since(start)
	fmt.Printf("It took %d us to print previous message\n", duration.Microseconds())
}
