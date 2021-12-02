package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

const (
	slidingWindowSize = 3
)

type slidingWindow struct {
	// q doesn't support non-full sliding window.
	q       [slidingWindowSize]int
	tailIdx int
	sum     int
}

func (w *slidingWindow) Add(v int) {
	w.sum -= w.q[w.tailIdx]
	w.q[w.tailIdx] = v
	w.tailIdx = (w.tailIdx + 1) % slidingWindowSize
	w.sum += v
}

func (w *slidingWindow) Sum() int {
	return w.sum
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)

	var window slidingWindow
	for i := 0; i < slidingWindowSize; i++ {
		if sc.Scan() {
			v, err := strconv.Atoi(sc.Text())
			if err != nil {
				log.Panicf("failed to parse first line depth value: %s", err)
			}
			window.Add(v)
		}
		if err = sc.Err(); err != nil {
			log.Panicf("failed to read first line: %s", err)
		}
	}

	var (
		prevSum = window.Sum()
		count   = 0
	)
	for i := slidingWindowSize + 1; sc.Scan(); i++ {
		v, err := strconv.Atoi(sc.Text())
		if err != nil {
			log.Panicf("failed to parse depth value on line %d: %s", i, err)
		}

		window.Add(v)
		currSum := window.Sum()
		if currSum > prevSum {
			count++
		}
		prevSum = currSum
	}
	if err = sc.Err(); err != nil {
		log.Panicf("failed to scan: %s", err)
	}

	fmt.Println(count)
}
