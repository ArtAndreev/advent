package main

import (
	"fmt"
	"os"
)

const markerLen = 14

func main() {
	answer, err := solveTuningTrouble(markerLen)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to solve tuning trouble: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(answer)
}
