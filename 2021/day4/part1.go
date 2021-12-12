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

	winBoardIdx, winDrawNumberIdx, ok := playGame(boards, rowNum, colNum, drawNumbers)
	if !ok {
		log.Panicln("no winner at the game")
	}
	winDrawNumber := drawNumbers[winDrawNumberIdx]
	fmt.Printf("board %d won with number %d\n", winBoardIdx, winDrawNumber)

	unmarkedSum := 0
	for _, c := range boards[winBoardIdx] {
		if !c.marked {
			unmarkedSum += c.number
		}
	}

	fmt.Println(unmarkedSum * winDrawNumber)
}
