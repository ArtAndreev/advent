package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)

	drawNumbers, err := readDrawNumbers(sc)
	if err != nil {
		log.Panicf("failed to read draw numbers: %s", err)
	}
	assertNoDuplicateNumbers(drawNumbers)

	boards, rowNum, colNum, err := readBoards(sc)
	if err != nil {
		log.Panicf("failed to read boards: %s", err)
	}
	assertNoDuplicateCellNumbers(boards)

	unmarkedSum, lastWinDrawNumberIdx, ok := playSquidGame(boards, rowNum, colNum, drawNumbers)
	if !ok {
		log.Panicln("no winner at the game")
	}
	lastWinDrawNumber := drawNumbers[lastWinDrawNumberIdx]

	fmt.Println(unmarkedSum * lastWinDrawNumber)
}
