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
	head, tail := point{0, 0}, point{0, 0}
	for i := 1; sc.Scan(); i++ {
		mot, err := parseMotion(sc.Text())
		if err != nil {
			return 0, fmt.Errorf("parse motion on line %d: %v", i, err)
		}

		for i := 0; i < mot.count; i++ {
			prevHeadPoint := head
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
			diffX := abs(head.x - tail.x)
			diffY := abs(head.y - tail.y)
			if diffX > 1 || diffY > 1 {
				tail = prevHeadPoint
				visited[tail] = true
			}
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
