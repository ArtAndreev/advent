package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
)

func main() {
	answer, err := solveBeaconExclusionZone()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to solve beacon exclusion zone: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(answer)
}

const calculatedRowY = 2000000

type sensor struct {
	location point
	beacon   point
}

type point struct {
	x, y int
}

type scanRange struct {
	from, to int // inclusive.
}

func solveBeaconExclusionZone() (int, error) {
	sensors, err := parseSensors()
	if err != nil {
		return 0, fmt.Errorf("parse sensors: %v", err)
	}

	var scannedRanges []scanRange
	beaconsXsInCalculatedRow := make(map[int]bool)
	for _, s := range sensors {
		manhattanDistance := abs(s.beacon.x-s.location.x) + abs(s.beacon.y-s.location.y)

		deltaY := abs(s.location.y - calculatedRowY)
		manhattanDistance -= deltaY
		if manhattanDistance < 0 {
			continue
		}

		scannedRanges = append(scannedRanges, scanRange{s.location.x - manhattanDistance, s.location.x + manhattanDistance})

		if s.beacon.y == calculatedRowY {
			beaconsXsInCalculatedRow[s.beacon.x] = true
		}
	}

	if len(scannedRanges) == 0 {
		return 0, nil
	}

	sort.Slice(scannedRanges, func(i, j int) bool {
		return scannedRanges[i].from < scannedRanges[j].from ||
			scannedRanges[i].from == scannedRanges[j].from && scannedRanges[i].to < scannedRanges[j].to
	})
	var mergedRanges []scanRange
	prev := scannedRanges[0]
	for _, sr := range scannedRanges[1:] {
		if sr.from > prev.to {
			mergedRanges = append(mergedRanges, prev)
			prev = sr
			continue
		}

		if sr.to > prev.to {
			prev.to = sr.to
		}
	}
	mergedRanges = append(mergedRanges, prev)

	bannedPositions := 0
	for _, mr := range mergedRanges {
		bannedPositions += mr.to - mr.from + 1
	}
	return bannedPositions - len(beaconsXsInCalculatedRow), nil
}

func abs(i int) int {
	if i >= 0 {
		return i
	}
	return -i
}

var sensorRe = regexp.MustCompile(`^Sensor at x=(-?\d+), y=(-?\d+): closest beacon is at x=(-?\d+), y=(-?\d+)$`)

func parseSensors() ([]sensor, error) {
	f, err := os.Open("input.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)

	var sensors []sensor
	for line := 1; sc.Scan(); line++ {
		matches := sensorRe.FindSubmatch(sc.Bytes())
		if len(matches) == 0 {
			return nil, fmt.Errorf("wrong line %d format, must be %s", line, sensorRe.String())
		}

		nums := make([]int, 0, 4)
		for i, m := range matches[1:] {
			num, err := strconv.Atoi(string(m))
			if err != nil {
				return nil, fmt.Errorf("wrong %s at line %d: %v", getNumName(i), line, err)
			}
			nums = append(nums, num)
		}

		sensors = append(sensors, sensor{
			location: point{
				x: nums[0],
				y: nums[1],
			},
			beacon: point{
				x: nums[2],
				y: nums[3],
			},
		})
	}
	if err = sc.Err(); err != nil {
		return nil, fmt.Errorf("scan: %v", err)
	}

	return sensors, nil
}

func getNumName(i int) string {
	switch i {
	case 0:
		return "sensor x coordinate"
	case 1:
		return "sensor y coordinate"
	case 2:
		return "beacon x coordinate"
	case 3:
		return "beacon y coordinate"
	default:
		panic("unknown num")
	}
}
