package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)

	var (
		lineCount      = 0
		countOfOneBits []int
	)
	for ; sc.Scan(); lineCount++ {
		binNum := sc.Text()
		if len(binNum) > len(countOfOneBits) {
			// Add leading zeros.
			countOfOneBits = append(make([]int, len(binNum)-len(countOfOneBits), len(binNum)), countOfOneBits...)
		}
		for i := range binNum {
			switch b := binNum[len(binNum)-1-i]; b {
			case '0':
			case '1':
				countOfOneBits[len(countOfOneBits)-1-i]++
			default:
				log.Panicf("unknown bit value %q on line %d", b, lineCount)
			}
		}
	}
	if err = sc.Err(); err != nil {
		log.Panicf("failed to scan: %s", err)
	}

	var gammaRate int
	for i := range countOfOneBits {
		oneCount := countOfOneBits[len(countOfOneBits)-1-i]
		zeroCount := lineCount - oneCount
		if oneCount > zeroCount {
			gammaRate ^= 1 << i
		}
	}
	epsilonRate := gammaRate
	for i := range countOfOneBits {
		epsilonRate ^= 1 << i
	}

	fmt.Println(gammaRate * epsilonRate)
}
