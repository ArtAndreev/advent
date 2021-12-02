package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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
		prev  int
		count = 0
	)
	if sc.Scan() {
		if prev, err = strconv.Atoi(sc.Text()); err != nil {
			log.Panicf("failed to parse first line depth value: %s", err)
		}
	}
	if err = sc.Err(); err != nil {
		log.Panicf("failed to read first line: %s", err)
	}

	for i := 1; sc.Scan(); i++ {
		curr, err := strconv.Atoi(sc.Text())
		if err != nil {
			log.Panicf("failed to parse depth value on line %d: %s", i, err)
		}
		if curr > prev {
			count++
		}
		prev = curr
	}
	if err = sc.Err(); err != nil {
		log.Panicf("failed to scan: %s", err)
	}

	fmt.Println(count)
}
