package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func main() {
	answer, err := supplyStacks()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to supply stacks: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", answer)
}

type stack struct {
	arr []rune
}

func (s *stack) Push(v rune) {
	s.arr = append(s.arr, v)
}

func (s *stack) Pop() (v rune, ok bool) {
	if len(s.arr) == 0 {
		return 0, false
	}

	v = s.arr[len(s.arr)-1]
	s.arr = s.arr[:len(s.arr)-1]
	return v, true
}

type movement struct {
	count    int
	from, to int
}

func supplyStacks() (string, error) {
	f, err := os.Open("input.txt")
	if err != nil {
		return "", err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)

	var i int
	var levels []string
	for i = 1; sc.Scan(); i++ {
		t := sc.Text()
		if t == "" {
			break
		}

		levels = append(levels, t)
	}
	if len(levels) < 2 {
		return "", errors.New("level count (including indices) must be greater than 2")
	}

	// Ignore indices line, because it's unused.
	levels = levels[:len(levels)-1]

	stacks, err := parseStacksByLevels(levels)
	if err != nil {
		return "", fmt.Errorf("parse stacks by levels: %v", err)
	}

	for ; sc.Scan(); i++ {
		move, err := parseMovement(sc.Bytes())
		if err != nil {
			return "", fmt.Errorf("parse movement on line %d: %v", i, err)
		}
		if move.from == 0 {
			return "", fmt.Errorf("movement direction from must start from 1, line %d", i)
		}
		move.from--
		if move.from > len(stacks)-1 {
			return "", fmt.Errorf("movement direction from out of range, line %d", i)
		}
		if move.to == 0 {
			return "", fmt.Errorf("movement direction to must start from 1, line %d", i)
		}
		move.to--
		if move.to > len(stacks)-1 {
			return "", fmt.Errorf("movement direction to out of range, line %d", i)
		}

		if move.from == move.to {
			continue
		}

		for j := 0; j < move.count; j++ {
			c, ok := stacks[move.from].Pop()
			if !ok {
				return "", fmt.Errorf("movement on line %d failed, stack in position %d is empty", i, move.from+1)
			}
			stacks[move.to].Push(c)
		}
	}
	if err = sc.Err(); err != nil {
		return "", fmt.Errorf("scan: %v", err)
	}

	var topCrates []rune
	for i, st := range stacks {
		c, ok := st.Pop()
		if !ok {
			return "", fmt.Errorf("stack %d is empty after all movements", i+1)
		}
		topCrates = append(topCrates, c)
	}

	return string(topCrates), nil
}

func parseStacksByLevels(levels []string) ([]stack, error) {
	var stacks []stack
	for i := len(levels) - 1; i >= 0; i-- {
		level := levels[i]
		if level[0] == ' ' {
			return nil, errors.New("levels must not have leading whitespaces")
		}

		stackIdx := -1
		startedCrate := false
		symbolAppended := false
		for j, ch := range level {
			if j%4 == 0 {
				stackIdx++
			}
			switch ch {
			case '[':
				if startedCrate {
					return nil, fmt.Errorf("malformed crate on line %d", i+1)
				}
				startedCrate = true
			case ']':
				if !startedCrate {
					return nil, fmt.Errorf("malformed crate on line %d", i+1)
				}
				startedCrate = false
				symbolAppended = false
			case ' ':
				if startedCrate || symbolAppended {
					return nil, fmt.Errorf("unknown symbol on crate in stack %d on line %d", stackIdx+1, i+1)
				}
			default:
				if !startedCrate {
					return nil, fmt.Errorf("malformed crate on line %d", i+1)
				}
				if symbolAppended || ch < 'A' || ch > 'Z' {
					return nil, fmt.Errorf("unknown symbol on crate in stack %d on line %d", stackIdx+1, i+1)
				}
				if len(stacks) <= stackIdx {
					stacks = append(stacks, stack{})
				}
				stacks[stackIdx].Push(ch)
				symbolAppended = true
			}
		}
	}

	return stacks, nil
}

var moveRe = regexp.MustCompile(`^move (\d+) from (\d+) to (\d+)$`)

func parseMovement(b []byte) (movement, error) {
	matches := moveRe.FindSubmatch(b)
	if len(matches) == 0 {
		return movement{}, errors.New("malformed movement")
	}

	count, err := strconv.Atoi(string(matches[1])) // count is non-negative.
	if err != nil {
		return movement{}, fmt.Errorf("parse count: %v", err)
	}
	from, err := strconv.Atoi(string(matches[2])) // from is non-negative.
	if err != nil {
		return movement{}, fmt.Errorf("parse from: %v", err)
	}
	to, err := strconv.Atoi(string(matches[3])) // to is non-negative.
	if err != nil {
		return movement{}, fmt.Errorf("parse to: %v", err)
	}

	return movement{
		count: count,
		from:  from,
		to:    to,
	}, nil
}
