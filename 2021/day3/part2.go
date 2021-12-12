package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)

	var (
		values         []int
		maxBinValueLen = 0
	)
	for i := 0; sc.Scan(); i++ {
		t := sc.Text()
		if len(t) > maxBinValueLen {
			maxBinValueLen = len(t)
		}
		v, err := strconv.ParseInt(t, 2, 64)
		if err != nil {
			log.Panicf("failed to parse value on line %d: %s", i, err)
		}
		values = append(values, int(v))
	}
	if err = sc.Err(); err != nil {
		log.Panicf("failed to scan: %s", err)
	}

	if len(values) == 0 {
		log.Panicln("file is empty")
	}

	oxygenGR := findRating(values, maxBinValueLen, func(zerosCount, onesCount int) bool {
		return onesCount >= zerosCount
	})
	fmt.Printf("found oxygen GR %d (%b)\n", oxygenGR, oxygenGR)

	co2SR := findRating(values, maxBinValueLen, func(zerosCount, onesCount int) bool {
		return onesCount < zerosCount
	})
	fmt.Printf("found co2 GR %d (%b)\n", co2SR, co2SR)

	fmt.Println(oxygenGR * co2SR)
}

func findRating(values []int, maxBinValueLen int, keepOnes func(zerosCount, onesCount int) bool) int {
	left, right := 0, len(values)
	for pos := maxBinValueLen - 1; left-right != 1 && pos >= 0; pos-- {
		part := values[left:right]

		sort.Slice(part, func(i, j int) bool {
			return (part[i]>>pos)&1 < (part[j]>>pos)&1
		})

		idx := sort.Search(len(part), func(i int) bool { return (part[i]>>pos)&1 == 1 })
		zerosCount := idx
		onesCount := len(part) - zerosCount
		if keepOnes(zerosCount, onesCount) {
			left += idx
		} else {
			right = left + idx
		}
	}

	return values[left]
}
