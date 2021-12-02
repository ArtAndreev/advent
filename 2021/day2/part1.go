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
		horizontal = 0
		depth      = 0
	)
	for i := 0; sc.Scan(); i++ {
		var (
			name  string
			value int
		)
		if _, err = fmt.Sscanf(sc.Text(), "%s %d", &name, &value); err != nil {
			log.Panicf("failed to parse depth value on line %d: %s", i, err)
		}

		switch name {
		case "forward":
			horizontal += value
		case "down":
			depth += value
		case "up":
			depth -= value
		default:
			log.Panicf("unknown command name %q on line %d", name, i)
		}
	}
	if err = sc.Err(); err != nil {
		log.Panicf("failed to scan: %s", err)
	}

	fmt.Println(horizontal * depth)
}
