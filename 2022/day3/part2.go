package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	answer, err := reorganizeRucksack()
	if err != nil {
		log.Fatalf("failed to reorganize rucksack: %v", err)
	}
	fmt.Println(answer)
}

func reorganizeRucksack() (int, error) {
	f, err := os.Open("input.txt")
	if err != nil {
		return 0, fmt.Errorf("open file: %s", err)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)

	prioritySum := 0
	sharedByGroupItems := make(map[rune]bool)
	var i int
	for i = 1; sc.Scan(); i++ {
		rucksack := sc.Text()
		items := make(map[rune]bool)
		for _, i := range rucksack {
			items[i] = true
		}

		if i%3 != 1 {
			sharedByGroupItems = findSetIntersection(sharedByGroupItems, items)
			if len(sharedByGroupItems) == 0 {
				return 0, fmt.Errorf("rucksack on line %d has no shared items with previous", i)
			}
		} else {
			sharedByGroupItems = items
		}

		if i%3 != 0 {
			continue
		}

		if len(sharedByGroupItems) > 1 {
			return 0, fmt.Errorf("group on lines %d-%d has multiple shared items", i-2, i)
		}
		var sharedItem rune
		for i := range sharedByGroupItems {
			sharedItem = i
		}

		switch {
		case sharedItem >= 'A' && sharedItem <= 'Z':
			prioritySum += 27 + int(sharedItem-'A')
		case sharedItem >= 'a' && sharedItem <= 'z':
			prioritySum += 1 + int(sharedItem-'a')
		default:
			return 0, fmt.Errorf("unknown shared item %c on line %d", sharedItem, i)
		}
	}
	if err = sc.Err(); err != nil {
		return 0, fmt.Errorf("scan: %s", err)
	}
	if i%3 != 1 {
		return 0, fmt.Errorf("rucksack on line %d isn't last in group of three", i)
	}

	return prioritySum, nil
}

func findSetIntersection(l, r map[rune]bool) map[rune]bool {
	res := make(map[rune]bool)
	for v := range r {
		if l[v] {
			res[v] = true
		}
	}

	return res
}
