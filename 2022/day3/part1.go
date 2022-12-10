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
	for i := 1; sc.Scan(); i++ {
		rucksack := sc.Text()
		itemsInFirstCompartment := make(map[rune]bool)
		if len(rucksack)%2 != 0 {
			return 0, fmt.Errorf("rucksack has odd item count on line %d", i)
		}
		for _, i := range rucksack[:len(rucksack)/2] {
			itemsInFirstCompartment[i] = true
		}
		var sharedItem rune
		for _, i := range rucksack[len(rucksack)/2:] {
			if itemsInFirstCompartment[i] {
				sharedItem = i
				break
			}
		}

		if sharedItem == 0 {
			return 0, fmt.Errorf("rucksack doesn't have shared item on line %d", i)
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

	return prioritySum, nil
}
