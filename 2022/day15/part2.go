package main

import (
	"bufio"
	"errors"
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

const (
	distressBeaconMinCoordinate = 0
	distressBeaconMaxCoordinate = 4000000
	tuningFrequencyXMultiplier  = 4000000
)

type sensor struct {
	location point
	beacon   point
}

type point struct {
	x, y int
}

type sensorRegion struct {
	sensor sensor
	region scanRegion
}

type scanRegion struct {
	xRange, yRange scanRange
}

type scanRange struct {
	from, to int // inclusive.
}

func solveBeaconExclusionZone() (int, error) {
	sensors, err := parseSensors()
	if err != nil {
		return 0, fmt.Errorf("parse sensors: %v", err)
	}

	// Add covered by sensors regions. They have shape of square, rotated by 45 degrees.
	var sensorsRegions []sensorRegion
	for _, s := range sensors {
		manhattanDistance := abs(s.beacon.x-s.location.x) + abs(s.beacon.y-s.location.y)
		r := scanRegion{
			xRange: scanRange{s.location.x - manhattanDistance, s.location.x + manhattanDistance},
			yRange: scanRange{s.location.y - manhattanDistance, s.location.y + manhattanDistance},
		}
		if r.xRange.to < distressBeaconMinCoordinate || r.xRange.from > distressBeaconMaxCoordinate ||
			r.yRange.to < distressBeaconMinCoordinate || r.yRange.from > distressBeaconMaxCoordinate {
			continue
		}

		sensorsRegions = append(sensorsRegions, sensorRegion{
			sensor: s,
			region: r,
		})
	}

	if len(sensorsRegions) == 0 {
		return 0, fmt.Errorf("sensor regions not found")
	}

	// Sort from top to bottom, from left to right.
	sort.Slice(sensorsRegions, func(i, j int) bool {
		l, r := sensorsRegions[i].region, sensorsRegions[j].region
		return l.yRange.from < r.yRange.from ||
			l.yRange.from == r.yRange.from &&
				(l.yRange.to < r.yRange.to || l.yRange.to == r.yRange.to &&
					(l.xRange.from < r.xRange.from || l.xRange.from == r.xRange.from && l.xRange.to < r.xRange.to))
	})

	for calculatedRowY := distressBeaconMinCoordinate; calculatedRowY <= distressBeaconMaxCoordinate; calculatedRowY++ {
		var coveredXRanges []scanRange
		for _, sensorRegion := range sensorsRegions {
			region := sensorRegion.region
			if region.yRange.to < calculatedRowY {
				continue
			}
			if region.yRange.from > calculatedRowY {
				break
			}
			// Region covers this row.

			deltaY := abs(sensorRegion.sensor.location.y - calculatedRowY)
			coveredXRanges = append(coveredXRanges,
				scanRange{
					max(region.xRange.from+deltaY, distressBeaconMinCoordinate),
					min(region.xRange.to-deltaY, distressBeaconMaxCoordinate),
				})
		}

		mergedRanges := mergeRanges(coveredXRanges)

		if len(mergedRanges) == 1 &&
			mergedRanges[0].from == distressBeaconMinCoordinate && mergedRanges[0].to == distressBeaconMaxCoordinate {
			continue
		}

		distressBeaconX := 0
		for _, r := range mergedRanges {
			if r.from > distressBeaconX {
				return distressBeaconX*tuningFrequencyXMultiplier + calculatedRowY, nil
			}

			distressBeaconX = r.to + 1
		}
		return 0, fmt.Errorf("found gaps at row y %d, but no distress beacon", calculatedRowY)
	}

	return 0, errors.New("distress beacon not found")
}

func mergeRanges(rs []scanRange) []scanRange {
	sort.Slice(rs, func(i, j int) bool {
		return rs[i].from < rs[j].from ||
			rs[i].from == rs[j].from && rs[i].to < rs[j].to
	})
	var merged []scanRange
	prev := rs[0]
	for _, r := range rs[1:] {
		if r.from > prev.to {
			merged = append(merged, prev)
			prev = r
			continue
		}

		if r.to > prev.to {
			prev.to = r.to
		}
	}
	merged = append(merged, prev)

	return merged
}

func max(l, r int) int {
	if l >= r {
		return l
	}
	return r
}

func min(l, r int) int {
	if l <= r {
		return l
	}
	return r
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
