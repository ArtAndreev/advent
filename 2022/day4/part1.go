package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	answer, err := cleanupCamp()
	if err != nil {
		fmt.Fprintf(os.Stdout, "failed to cleanup camp: %v", err)
		os.Exit(1)
	}
	fmt.Println(answer)
}

type section struct {
	begin, end int
}

func cleanupCamp() (int, error) {
	f, err := os.Open("input.txt")
	if err != nil {
		return 0, err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)

	pairsCount := 0
	for i := 1; sc.Scan(); i++ {
		pairSections := strings.Split(sc.Text(), ",")
		if len(pairSections) != 2 {
			return 0, fmt.Errorf("not pair sections on line %d", i)
		}

		var sections []section
		for j, s := range pairSections {
			chunks := strings.SplitN(s, "-", 2)
			if len(chunks) < 2 {
				return 0, fmt.Errorf("wrong section %d value on line %d: must be `begin-end` format", j, i)
			}
			begin, err := strconv.Atoi(chunks[0])
			if err != nil {
				return 0, fmt.Errorf("wrong begin value in section %d on line %d: %v", j, i, err)
			}
			end, err := strconv.Atoi(chunks[1])
			if err != nil {
				return 0, fmt.Errorf("wrong end value in section %d on line %d: %v", j, i, err)
			}
			if end < begin {
				return 0, fmt.Errorf("wrong section %d on line %d: begin must be less or equal to end", j, i)
			}

			sections = append(sections, section{
				begin: begin,
				end:   end,
			})
		}

		sort.Slice(sections, func(i, j int) bool {
			if sections[i].begin < sections[j].begin {
				return true
			}
			if sections[i].begin > sections[j].begin {
				return false
			}
			return sections[i].end > sections[j].end
		})

		if sections[0].end >= sections[1].end {
			pairsCount++
		}
	}
	if err = sc.Err(); err != nil {
		return 0, fmt.Errorf("scan: %v", err)
	}

	return pairsCount, nil
}
