package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	answer, err := playRockPaperScissors()
	if err != nil {
		log.Fatalf("failed to play rock paper scissors: %v", err)
	}
	fmt.Println(answer)
}

const (
	loseStrategy = "X"
	drawStrategy = "Y"
	winStrategy  = "Z"

	rockShape     = "A"
	paperShape    = "B"
	scissorsShape = "C"
)

var (
	// strategies are mapping strategy -> op shape -> my shape.
	strategies = map[string]map[string]string{
		loseStrategy: {
			rockShape:     scissorsShape,
			paperShape:    rockShape,
			scissorsShape: paperShape,
		},
		drawStrategy: {
			rockShape:     rockShape,
			paperShape:    paperShape,
			scissorsShape: scissorsShape,
		},
		winStrategy: {
			rockShape:     paperShape,
			paperShape:    scissorsShape,
			scissorsShape: rockShape,
		},
	}
	shapeScores = map[string]int{
		rockShape:     1,
		paperShape:    2,
		scissorsShape: 3,
	}
	outcomeScores = map[string]int{
		loseStrategy: 0,
		drawStrategy: 3,
		winStrategy:  6,
	}
)

func playRockPaperScissors() (int, error) {
	f, err := os.Open("input.txt")
	if err != nil {
		return 0, fmt.Errorf("open file: %s", err)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)

	totalScore := 0
	for i := 1; sc.Scan(); i++ {
		shapes := strings.Fields(sc.Text())
		if len(shapes) != 2 {
			return 0, fmt.Errorf("wrong shape count on line %d", i)
		}

		opShape, myShape := shapes[0], shapes[1]
		if myShape < "X" || myShape > "Z" {
			return 0, fmt.Errorf("wrong my shape on line %d", i)
		}
		if opShape < "A" || opShape > "C" {
			return 0, fmt.Errorf("wrong op shape on line %d", i)
		}

		totalScore += outcomeScores[myShape]
		totalScore += shapeScores[strategies[myShape][opShape]]
	}
	if err = sc.Err(); err != nil {
		return 0, fmt.Errorf("scan: %s", err)
	}

	return totalScore, nil
}
