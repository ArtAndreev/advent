package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	answer, err := solveRopeBridge()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to solve rope bridge %v\n", err)
		os.Exit(1)
	}
	fmt.Println(answer)
}

const knotCount = 10

type direction string

const (
	directionUp    direction = "U"
	directionDown  direction = "D"
	directionLeft  direction = "L"
	directionRight direction = "R"
)

type point struct {
	x, y int
}

type motion struct {
	direction direction
	count     int
}

func solveRopeBridge() (int, error) {
	f, err := os.Open("input.txt")
	if err != nil {
		return 0, err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)

	visited := map[point]bool{
		{0, 0}: true,
	}
	knots := make([]point, knotCount)
	for i := 1; sc.Scan(); i++ {
		mot, err := parseMotion(sc.Text())
		if err != nil {
			return 0, fmt.Errorf("parse motion on line %d: %v", i, err)
		}

	MOVE:
		for i := 0; i < mot.count; i++ {
			head := &knots[0]
			switch mot.direction {
			case directionUp:
				head.y++
			case directionDown:
				head.y--
			case directionLeft:
				head.x--
			case directionRight:
				head.x++
			}
			prevKnot := *head
			for i, k := range knots[1:] {
				diffX := prevKnot.x - k.x
				diffY := prevKnot.y - k.y
				if diffX >= -1 && diffX <= 1 && diffY >= -1 && diffY <= 1 {
					continue MOVE
				}
				if diffX != 0 {
					knots[i+1].x += diffX / abs(diffX) // move by 1 or -1.
				}
				if diffY != 0 {
					knots[i+1].y += diffY / abs(diffY) // move by 1 or -1.
				}
				prevKnot = knots[i+1]
			}
			visited[knots[len(knots)-1]] = true
		}
	}
	if err = sc.Err(); err != nil {
		return 0, fmt.Errorf("scan: %v", err)
	}

	return len(visited), nil
}

func abs(v int) int {
	if v >= 0 {
		return v
	}
	return -v
}

func parseMotion(s string) (motion, error) {
	chunks := strings.SplitN(s, " ", 2)
	if len(chunks) < 2 {
		return motion{}, errors.New("expected motion in format '[UDLR] (\\d)+'")
	}

	dir := direction(chunks[0])
	switch dir {
	case directionUp, directionDown, directionLeft, directionRight:
	default:
		return motion{}, fmt.Errorf("unknown direction %s", dir)
	}

	count, err := strconv.Atoi(chunks[1])
	if err != nil {
		return motion{}, fmt.Errorf("parse count: %v", err)
	}

	return motion{
		direction: dir,
		count:     count,
	}, nil
}
