package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	answer, err := countCalories()
	if err != nil {
		log.Fatalf("failed to count calories: %v", err)
	}
	fmt.Println(answer)
}

func countCalories() (int, error) {
	f, err := os.Open("input.txt")
	if err != nil {
		return 0, fmt.Errorf("open file: %s", err)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)

	var (
		mostTotalCalories = 0
		currTotalCalories = 0
	)
	for i := 1; sc.Scan(); i++ {
		t := sc.Text()
		if t == "" {
			if currTotalCalories > mostTotalCalories {
				mostTotalCalories = currTotalCalories
			}
			currTotalCalories = 0
			continue
		}

		calories, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("parse text on line %d: %v", i, err)
		}
		currTotalCalories += int(calories)
	}
	if err = sc.Err(); err != nil {
		return 0, fmt.Errorf("scan: %s", err)
	}

	return mostTotalCalories, nil
}
