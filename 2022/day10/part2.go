package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	answer, err := solveCathodeRayTube()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to solve cathode-ray tube %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%s", answer)
}

const (
	litPixel  = '#'
	darkPixel = '.'
)

func solveCathodeRayTube() ([]byte, error) {
	f, err := os.Open("input.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)

	screen := []byte{litPixel}
	drawPixel := func(spriteIdx int) {
		idx := len(screen) % 40
		if idx == spriteIdx || idx == spriteIdx-1 || idx == spriteIdx+1 {
			screen = append(screen, litPixel)
		} else {
			screen = append(screen, darkPixel)
		}
	}

	cycle := 1
	register := 1
	for i := 1; sc.Scan(); i++ {
		t := sc.Text()
		if t == "noop" {
			cycle++
			continue
		}

		chunks := strings.SplitN(t, " ", 2)
		if len(chunks) < 2 || chunks[0] != "addx" {
			return nil, fmt.Errorf("command on line %d is wrong, expected noop or addx with arg", i)
		}
		value, err := strconv.Atoi(chunks[1])
		if err != nil {
			return nil, fmt.Errorf("parse addx arg on line %d: %v", i, err)
		}

		cycle += 2
		for len(screen) < cycle-1 {
			drawPixel(register)
		}

		register += value
	}
	if err = sc.Err(); err != nil {
		return nil, fmt.Errorf("scan: %v", err)
	}

	for len(screen) < cycle-1 {
		drawPixel(register)
	}
	for ; cycle <= 240; cycle++ {
		drawPixel(register)
	}
	splitIdx := 40
	for splitIdx <= len(screen) {
		screen = append(screen[:splitIdx], append([]byte{'\n'}, screen[splitIdx:]...)...)
		splitIdx += 41
	}

	return screen, nil
}
