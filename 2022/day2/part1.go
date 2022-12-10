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
		switch myShape {
		case "X":
			myShape = "A"
			totalScore += 1
		case "Y":
			myShape = "B"
			totalScore += 2
		case "Z":
			myShape = "C"
			totalScore += 3
		default:
			return 0, fmt.Errorf("wrong my shape on line %d", i)
		}

		if opShape == myShape {
			totalScore += 3
			continue
		}

		switch opShape {
		case "A":
			if myShape == "B" {
				totalScore += 6
			}
		case "B":
			if myShape == "C" {
				totalScore += 6
			}
		case "C":
			if myShape == "A" {
				totalScore += 6
			}
		default:
			return 0, fmt.Errorf("wrong op shape on line %d", i)
		}
	}
	if err = sc.Err(); err != nil {
		return 0, fmt.Errorf("scan: %s", err)
	}

	return totalScore, nil
}
