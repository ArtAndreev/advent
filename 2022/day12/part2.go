package main

import (
	"bufio"
	"container/heap"
	"errors"
	"fmt"
	"os"
)

func main() {
	answer, err := solveHillClimbingAlgorithm()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to solve hill climbing algorithm: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(answer)
}

type point struct {
	X, Y int
}

func solveHillClimbingAlgorithm() (int, error) {
	rawHeightmap, err := parseHeightmap()
	if err != nil {
		return 0, fmt.Errorf("parse heightmap: %v", err)
	}

	var startPoints []point
	endPoint := point{-1, -1}
	for i, row := range rawHeightmap {
		for j, el := range row {
			switch el {
			case 'S':
				row[j] = 'a'
				fallthrough
			case 'a':
				startPoints = append(startPoints, point{j, i})
			case 'E':
				endPoint = point{j, i}
				row[j] = 'z'
			}
		}
	}
	if len(startPoints) == 0 {
		return 0, errors.New("start points not found")
	}
	if endPoint.X == -1 {
		return 0, errors.New("end not found")
	}

	heightmap := make([][]int, 0, len(rawHeightmap))
	for _, row := range rawHeightmap {
		newRow := make([]int, 0, len(row))
		for _, el := range row {
			newRow = append(newRow, int(el-'a'))
		}
		heightmap = append(heightmap, newRow)
	}

	// TODO: we can cache paths somehow.
	minStepCount := len(heightmap) * len(heightmap[0])
	for _, startPoint := range startPoints {
		stepCount := runAStar(heightmap, startPoint, endPoint)
		if stepCount == -1 {
			continue
		}
		if stepCount < minStepCount {
			minStepCount = stepCount
		}
	}

	return minStepCount, nil
}

func parseHeightmap() ([][]byte, error) {
	f, err := os.Open("input.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)

	var grid [][]byte
	for line := 1; sc.Scan(); line++ {
		var row []byte
		if line != 1 {
			row = make([]byte, 0, len(grid[0]))
		}
		for i, ch := range sc.Bytes() {
			if ch < 'a' && ch != 'E' && ch != 'S' || ch > 'z' {
				return nil, fmt.Errorf("wrong character %c on line:pos %d:%d, expected [a-zES]", ch, line, i+1)
			}
			row = append(row, ch)
		}

		if line != 1 && len(row) != len(grid[0]) {
			return nil, fmt.Errorf("expected grid, row on line %d has different length", line)
		}

		grid = append(grid, row)
	}
	if err = sc.Err(); err != nil {
		return nil, fmt.Errorf("scan: %v", err)
	}

	return grid, nil
}

type (
	priorityQueue struct {
		heap *priorityHeap
	}

	priorityHeap []prioritizedPoint

	prioritizedPoint struct {
		point    pointWithPathLen
		priority int
	}
)

func (h priorityHeap) Len() int               { return len(h) }
func (h priorityHeap) Less(i int, j int) bool { return h[i].priority < h[j].priority }
func (h priorityHeap) Swap(i int, j int)      { h[i], h[j] = h[j], h[i] }

func (h *priorityHeap) Push(x any) {
	*h = append(*h, x.(prioritizedPoint))
}

func (h *priorityHeap) Pop() any {
	old := *h
	x := old[len(old)-1]
	*h = old[:len(old)-1]
	return x
}

func newPriorityQueue(vs ...prioritizedPoint) *priorityQueue {
	h := (*priorityHeap)(&vs)
	heap.Init(h)

	return &priorityQueue{
		heap: h,
	}
}

func (q *priorityQueue) Push(v prioritizedPoint) {
	heap.Push(q.heap, v)
}

func (q *priorityQueue) Pop() (pointWithPathLen, bool) {
	if q.heap.Len() == 0 {
		return pointWithPathLen{}, false
	}

	return heap.Pop(q.heap).(prioritizedPoint).point, true
}

type pointWithPathLen struct {
	point      point
	stepsCount int
}

func runAStar(heightmap [][]int, startPoint, endPoint point) int {
	pq := newPriorityQueue(newPrioritizedPoint(0, startPoint, endPoint))
	visited := make(map[point]bool)
	for v, ok := pq.Pop(); ok; v, ok = pq.Pop() {
		if visited[v.point] {
			continue
		}
		visited[v.point] = true

		if v.point == endPoint {
			return v.stepsCount
		}

		nextPoints := [4]point{}
		nextPointsLen := 0
		if topY := v.point.Y - 1; topY >= 0 {
			nextPoints[nextPointsLen] = point{X: v.point.X, Y: topY}
			nextPointsLen++
		}
		if rightX := v.point.X + 1; rightX < len(heightmap[0]) {
			nextPoints[nextPointsLen] = point{X: rightX, Y: v.point.Y}
			nextPointsLen++
		}
		if bottomY := v.point.Y + 1; bottomY < len(heightmap) {
			nextPoints[nextPointsLen] = point{X: v.point.X, Y: bottomY}
			nextPointsLen++
		}
		if leftX := v.point.X - 1; leftX >= 0 {
			nextPoints[nextPointsLen] = point{X: leftX, Y: v.point.Y}
			nextPointsLen++
		}

		currentPointValue := heightmap[v.point.Y][v.point.X]
		for _, nextPoint := range nextPoints[:nextPointsLen] {
			if !visited[nextPoint] {
				nextPointValue := heightmap[nextPoint.Y][nextPoint.X]
				if nextPointValue-currentPointValue < 2 {
					pq.Push(newPrioritizedPoint(v.stepsCount+1, nextPoint, endPoint))
				}
			}
		}
	}

	return -1
}

func newPrioritizedPoint(stepsCount int, newPoint, endPoint point) prioritizedPoint {
	return prioritizedPoint{
		point:    pointWithPathLen{point: newPoint, stepsCount: stepsCount},
		priority: f(newPoint, endPoint, stepsCount),
	}
}

func f(currentPoint, endPoint point, stepsCount int) int {
	return g(stepsCount) + h(currentPoint, endPoint)
}

func g(stepsCount int) int {
	return stepsCount
}

func h(currentPoint, endPoint point) int {
	return countManhattanDistance(currentPoint, endPoint)
}

func countManhattanDistance(l, r point) int {
	return abs(l.X-r.X) + abs(l.Y-r.Y)
}

func abs(l int) int {
	if l >= 0 {
		return l
	}
	return -l
}
