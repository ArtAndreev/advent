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

	counted := make(map[rowCol]bool)
	// From left to right and from right to left.
	for i, row := range grid {
		// Rightmost tree index.
		maxHeight, maxHeightTreeIdx := -1, -1
		// From left to right.
		for j, v := range row {
			if v >= maxHeight {
				if v > maxHeight {
					maxHeight = v
					counted[rowCol{i, j}] = true
				}
				maxHeightTreeIdx = j
			}
		}
		// From right to left.
		maxHeight = -1
		for j := len(row) - 1; j >= maxHeightTreeIdx; j-- {
			v := row[j]
			if v > maxHeight {
				maxHeight = v
				counted[rowCol{i, j}] = true
			}
		}
	}
	// From top to bottom and from bottom to top.
	for colIdx := range grid[0] {
		// Bottommost tree index.
		maxHeight, maxHeightTreeIdx := -1, -1
		// From top to bottom.
		for i, row := range grid {
			v := row[colIdx]
			if v >= maxHeight {
				if v > maxHeight {
					maxHeight = v
					counted[rowCol{i, colIdx}] = true
				}
				maxHeightTreeIdx = i
			}
		}
		// From bottom to top.
		maxHeight = -1
		for i := len(grid) - 1; i >= maxHeightTreeIdx; i-- {
			v := grid[i][colIdx]
			if v > maxHeight {
				maxHeight = v
				counted[rowCol{i, colIdx}] = true
			}
		}
	}

	return len(counted), nil
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
