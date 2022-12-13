package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	answer, err := solveTreetopTreeHouse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to solve treetop tree house: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(answer)
}

type rowCol struct {
	row, col int
}

func solveTreetopTreeHouse() (int, error) {
	grid, err := readGrid()
	if err != nil {
		return 0, fmt.Errorf("read grid: %v", err)
	}

	maxScenicScore := 0
	for i, row := range grid {
		for j := range row {
			scenicScore := calculateScenicScore(grid, i, j)
			if scenicScore > maxScenicScore {
				maxScenicScore = scenicScore
			}
		}
	}

	return maxScenicScore, nil
}

func readGrid() ([][]int, error) {
	f, err := os.Open("input.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)

	var grid [][]int
	for i := 1; sc.Scan(); i++ {
		var row []int
		if i != 1 {
			row = make([]int, 0, len(grid[0]))
		}
		for j, ch := range sc.Bytes() {
			if ch < '0' || ch > '9' {
				return nil, fmt.Errorf("wrong character %c on line:pos %d:%d, expected [0-9]", ch, i, j+1)
			}
			row = append(row, int(ch-'0'))
		}

		if i != 1 && len(row) != len(grid[0]) {
			return nil, fmt.Errorf("expected grid, row on line %d has different length", i)
		}

		grid = append(grid, row)
	}
	if err = sc.Err(); err != nil {
		return nil, fmt.Errorf("scan: %v", err)
	}

	return grid, nil
}

func calculateScenicScore(grid [][]int, rowIdx, colIdx int) int {
	if rowIdx == 0 || rowIdx == len(grid)-1 || colIdx == 0 || colIdx == len(grid[0])-1 {
		return 0
	}

	povTreeValue := grid[rowIdx][colIdx]
	povCount := [4]int{}
	// Up.
	currPovCount := 0
	for i := rowIdx - 1; i >= 0; i-- {
		currPovCount++
		v := grid[i][colIdx]
		if v >= povTreeValue {
			break
		}
	}
	povCount[0] = currPovCount
	// To the right.
	currPovCount = 0
	for _, v := range grid[rowIdx][colIdx+1:] {
		currPovCount++
		if v >= povTreeValue {
			break
		}
	}
	povCount[1] = currPovCount
	// Down.
	currPovCount = 0
	for i := rowIdx + 1; i < len(grid); i++ {
		currPovCount++
		v := grid[i][colIdx]
		if v >= povTreeValue {
			break
		}
	}
	povCount[2] = currPovCount
	// To the left.
	currPovCount = 0
	for i := colIdx - 1; i >= 0; i-- {
		currPovCount++
		v := grid[rowIdx][i]
		if v >= povTreeValue {
			break
		}
	}
	povCount[3] = currPovCount

	score := 1
	for _, v := range povCount {
		score *= v
	}
	return score
}
