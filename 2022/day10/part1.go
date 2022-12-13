package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	answer, err := solveCathodeRayTube()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to solve cathode-ray tube %v\n", err)
		os.Exit(1)
	}
	fmt.Println(answer)
}

func solveCathodeRayTube() (int, error) {
	f, err := os.Open("input.txt")
	if err != nil {
		return 0, err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)

	strengthSum := 0
	cycle := 1
	register := 1
	observedCycles := []int{20, 60, 100, 140, 180, 220}
	for i := 1; sc.Scan(); i++ {
		t := sc.Text()
		if t == "noop" {
			cycle++
			continue
		}

		chunks := strings.SplitN(t, " ", 2)
		if len(chunks) < 2 || chunks[0] != "addx" {
			return 0, fmt.Errorf("command on line %d is wrong, expected noop or addx with arg", i)
		}
		value, err := strconv.Atoi(chunks[1])
		if err != nil {
			return 0, fmt.Errorf("parse addx arg on line %d: %v", i, err)
		}

		cycle += 2
		for len(observedCycles) != 0 && cycle > observedCycles[0] {
			strengthSum += observedCycles[0] * register
			observedCycles = observedCycles[1:]
		}

		register += value
	}
	if err = sc.Err(); err != nil {
		return 0, fmt.Errorf("scan: %v", err)
	}

	for len(observedCycles) != 0 {
		strengthSum += observedCycles[0] * register
		observedCycles = observedCycles[1:]
	}

	return strengthSum, nil
}
