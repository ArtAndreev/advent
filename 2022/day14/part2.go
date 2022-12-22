package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
)

func main() {
	answer, err := solveRegolithReservoir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to solve regolith reservoir: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(answer)
}

type point struct {
	x, y int
}

type material int

const (
	materialAir material = iota
	materialRock
	materialSand
)

var startPoint = point{500, 0}

func solveRegolithReservoir() (int, error) {
	field, err := parseField()
	if err != nil {
		return 0, fmt.Errorf("parse field: %v", err)
	}
	// Add floor.
	drawRockPath(&field, point{0, len(field) + 1}, point{len(field[0]) - 1, len(field) + 1})

	if len(field[0]) <= startPoint.x {
		return 0, fmt.Errorf("wrong start x %d for field with max x %d", startPoint.x, len(field[0]))
	}
	if len(field) <= startPoint.y {
		return 0, fmt.Errorf("wrong start y %d for field with max y %d", startPoint.y, len(field))
	}

	sandCount := 0
	for getFieldValue(field, startPoint) != materialSand {
		// Spawn sand and drop it.
		sandPoint := startPoint
		for {
			downY := sandPoint.y + 1
			downPoint := point{sandPoint.x, downY}
			if getFieldValue(field, downPoint) == materialAir {
				sandPoint = downPoint
				continue
			}

			leftX := sandPoint.x - 1
			if leftX < 0 {
				expandFieldToTheLeft(&field)
				sandPoint.x++
				leftX++
				startPoint.x++
			}
			downLeftPoint := point{leftX, downY}
			if getFieldValue(field, downLeftPoint) == materialAir {
				sandPoint = downLeftPoint
				continue
			}

			rightX := sandPoint.x + 1
			if rightX == len(field[0]) {
				expandFieldToTheRight(&field)
			}
			downRightPoint := point{rightX, downY}
			if getFieldValue(field, downRightPoint) == materialAir {
				sandPoint = downRightPoint
				continue
			}

			setFieldValue(field, sandPoint, materialSand)
			sandCount++
			break
		}
	}

	return sandCount, nil
}

func expandFieldToTheLeft(field *[][]material) {
	for i, row := range *field {
		newRow := make([]material, len(row)+1)
		copy(newRow[1:], row)
		(*field)[i] = newRow
	}
	(*field)[len(*field)-1][0] = materialRock
}

func expandFieldToTheRight(field *[][]material) {
	for i, row := range *field {
		newRow := make([]material, len(row)+1)
		copy(newRow, row)
		(*field)[i] = newRow
	}
	lastRow := (*field)[len(*field)-1]
	lastRow[len(lastRow)-1] = materialRock
}

func parseField() ([][]material, error) {
	f, err := os.Open("input.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)

	var field [][]material
	for line := 1; sc.Scan(); line++ {
		chunks := bytes.Split(sc.Bytes(), []byte(" -> "))
		if len(chunks) == 0 {
			return nil, fmt.Errorf("wrong draw on line %d", line)
		}

		var prevPath *point
		for i, ch := range chunks {
			rawPath := bytes.Split(ch, []byte(","))
			if len(rawPath) != 2 {
				return nil, fmt.Errorf("wrong path %d on line %d", i, line)
			}
			x, err := strconv.Atoi(string(rawPath[0]))
			if err != nil {
				return nil, fmt.Errorf("parse first (x) coordinate from path %d on line %d", i, line)
			}
			if x < 0 {
				return nil, fmt.Errorf("coordinate (x) from path %d must be non-negative on line %d", i, line)
			}
			y, err := strconv.Atoi(string(rawPath[1]))
			if err != nil {
				return nil, fmt.Errorf("parse first (x) coordinate from path %d on line %d", i, line)
			}
			if y < 0 {
				return nil, fmt.Errorf("coordinate (y) from path %d must be non-negative on line %d", i, line)
			}
			currPath := point{x, y}

			if prevPath != nil {
				if currPath.x != prevPath.x && currPath.y != prevPath.y {
					return nil, fmt.Errorf("x,y coordinates from path %d on line %d draw not ortogonal path", i, line)
				}
				drawRockPath(&field, *prevPath, currPath)
			}

			prevPath = &currPath
		}
	}
	if err = sc.Err(); err != nil {
		return nil, fmt.Errorf("scan: %v", err)
	}

	return field, nil
}

func drawRockPath(field *[][]material, from, to point) {
	if len(*field) == 0 {
		*field = make([][]material, 1)
	}

	// First reallocate x, to decrease allocation number for y.
	maxX := max(from.x, to.x)
	if len((*field)[0]) <= maxX {
		for i, row := range *field {
			newRow := make([]material, maxX+1)
			copy(newRow, row)
			(*field)[i] = newRow
		}
	}
	maxY := max(from.y, to.y)
	if len(*field) <= maxY {
		newField := make([][]material, maxY+1)
		copy(newField, *field)
		for i := len(*field); i < len(newField); i++ {
			newField[i] = make([]material, len((*field)[0]))
		}
		*field = newField
	}

	if from.x == to.x {
		x := from.x
		minY, maxY := from.y, to.y
		if minY > maxY {
			minY, maxY = maxY, minY
		}
		for y := minY; y <= maxY; y++ {
			setFieldValue(*field, point{x, y}, materialRock)
		}
	} else if from.y == to.y {
		y := from.y
		minX, maxX := from.x, to.x
		if minX > maxX {
			minX, maxX = maxX, minX
		}
		for x := minX; x <= maxX; x++ {
			setFieldValue(*field, point{x, y}, materialRock)
		}
	}
}

func getFieldValue(field [][]material, pnt point) material {
	return field[pnt.y][pnt.x]
}

func setFieldValue(field [][]material, pnt point, value material) {
	field[pnt.y][pnt.x] = value
}

func max(l, r int) int {
	if l >= r {
		return l
	}
	return r
}
