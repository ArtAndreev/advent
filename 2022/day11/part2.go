package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	answer, err := solveMonkeyInTheMiddle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to solve monkey in the middle: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(answer)
}

type monkey struct {
	items                 []int
	operation             func(v int) int
	testDivisible         int
	trueThrow, falseThrow int
}

const (
	roundCount = 10000
	mostCount  = 2
)

func solveMonkeyInTheMiddle() (int, error) {
	monkeys, err := parseMonkeys()
	if err != nil {
		return 0, fmt.Errorf("parse monkeys: %v", err)
	}

	divisibleByTotal := 1
	for _, m := range monkeys {
		divisibleByTotal *= m.testDivisible
	}

	inspectedCount := make([]int, len(monkeys))
	for round := 0; round < roundCount; round++ {
		for i, m := range monkeys {
			inspectedCount[i] += len(m.items)

			for _, it := range m.items {
				// Inspect.
				it = m.operation(it) % divisibleByTotal
				// Test
				nextMonkey := monkeys[m.trueThrow]
				if it%m.testDivisible != 0 {
					nextMonkey = monkeys[m.falseThrow]
				}
				nextMonkey.items = append(nextMonkey.items, it)
			}
			m.items = m.items[:0]
		}
	}

	most := make([]int, mostCount)
	for _, v := range inspectedCount {
		for i, mostOne := range most {
			if v > mostOne {
				most = append(most[:i+1], most[i:len(most)-1]...)
				most[i] = v
				break
			}
		}
	}

	multiplication := 1
	for _, v := range most {
		multiplication *= v
	}
	return multiplication, nil
}

const (
	parseStateMonkey = iota
	parseStateItems
	parseStateOperation
	parseStateTestDivisible
	parseStateTrueThrow
	parseStateFalseThrow
)

const (
	monkeyPrefix        = "Monkey "
	monkeySuffix        = ":"
	itemsPrefix         = "  Starting items: "
	operationPrefix     = "  Operation: new = old "
	testDivisiblePrefix = "  Test: divisible by "
	trueThrowPrefix     = "    If true: throw to monkey "
	falseThrowPrefix    = "    If false: throw to monkey "

	operationSignPlus   = "+"
	operationSignStar   = "*"
	operationOldOperand = "old"
)

func getPrefixForParseState(s int) (string, error) {
	switch s {
	case parseStateMonkey:
		return monkeyPrefix, nil
	case parseStateItems:
		return itemsPrefix, nil
	case parseStateOperation:
		return operationPrefix, nil
	case parseStateTestDivisible:
		return testDivisiblePrefix, nil
	case parseStateTrueThrow:
		return trueThrowPrefix, nil
	case parseStateFalseThrow:
		return falseThrowPrefix, nil
	default:
		return "", fmt.Errorf("unknown parse state %d", s)
	}
}

func parseMonkeys() ([]*monkey, error) {
	f, err := os.Open("input.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)

	var (
		monkeys []*monkey

		monkeyNum  int
		parseState = parseStateMonkey
	)
	for line := 1; sc.Scan(); line++ {
		t := sc.Text()
		if parseState > parseStateFalseThrow {
			if t != "" {
				return nil, fmt.Errorf("unexpected input on line %d", line)
			}
			parseState = parseStateMonkey
			continue
		}

		prefix, err := getPrefixForParseState(parseState)
		if err != nil {
			return nil, fmt.Errorf("internal error on line %d: %v", line, err)
		}
		trimmed := strings.TrimPrefix(t, prefix)
		if len(trimmed) == len(t) {
			return nil, fmt.Errorf("state %d: no prefix %q on line %d", parseState, prefix, line)
		}
		switch parseState {
		case parseStateMonkey:
			rawMonkeyNum := strings.TrimSuffix(trimmed, monkeySuffix)
			if len(rawMonkeyNum) == len(trimmed) {
				return nil, fmt.Errorf("state %d: no suffix %q on line %d", parseState, monkeySuffix, line)
			}
			monkeyNum, err = strconv.Atoi(rawMonkeyNum)
			if err != nil {
				return nil, fmt.Errorf("state %d: wrong monkey num on line %d: %v", parseState, line, err)
			}
			if len(monkeys)-1 < monkeyNum {
				delta := monkeyNum - len(monkeys) + 1
				monkeys = append(monkeys, make([]*monkey, delta)...)
				monkeys[monkeyNum] = new(monkey)
			}
		case parseStateItems:
			rawItems := strings.Split(trimmed, ", ")
			monkey := monkeys[monkeyNum]
			for i, ri := range rawItems {
				item, err := strconv.Atoi(ri)
				if err != nil {
					return nil, fmt.Errorf("state %d: parse monkey item %d on line %d: %v", parseState, i, line, err)
				}
				monkey.items = append(monkey.items, item)
			}
		case parseStateOperation:
			rawOpAndValue := strings.SplitN(trimmed, " ", 2)
			if len(rawOpAndValue) < 2 {
				return nil, fmt.Errorf("state %d: parse op and value on line %d: wrong format", parseState, line)
			}
			op, rawValue := rawOpAndValue[0], rawOpAndValue[1]
			monkey := monkeys[monkeyNum]
			if rawValue == operationOldOperand {
				switch op {
				case operationSignPlus:
					monkey.operation = func(v int) int { return v + v }
				case operationSignStar:
					monkey.operation = func(v int) int { return v * v }
				default:
					return nil, fmt.Errorf("state %d: unknown operation sign %q on line %d", parseState, op, line)
				}
			} else {
				value, err := strconv.Atoi(rawValue)
				if err != nil {
					return nil, fmt.Errorf("state %d: parse op value on line %d: %v", parseState, line, err)
				}
				switch op {
				case operationSignPlus:
					monkey.operation = func(v int) int { return v + value }
				case operationSignStar:
					monkey.operation = func(v int) int { return v * value }
				default:
					return nil, fmt.Errorf("state %d: unknown operation sign %q on line %d", parseState, op, line)
				}
			}
		case parseStateTestDivisible:
			value, err := strconv.Atoi(trimmed)
			if err != nil {
				return nil, fmt.Errorf("state %d: parse divisible value on line %d: %v", parseState, line, err)
			}
			monkeys[monkeyNum].testDivisible = value
		case parseStateTrueThrow:
			value, err := strconv.Atoi(trimmed)
			if err != nil {
				return nil, fmt.Errorf("state %d: parse true throw value on line %d: %v", parseState, line, err)
			}
			monkeys[monkeyNum].trueThrow = value
		case parseStateFalseThrow:
			value, err := strconv.Atoi(trimmed)
			if err != nil {
				return nil, fmt.Errorf("state %d: parse false throw value on line %d: %v", parseState, line, err)
			}
			monkeys[monkeyNum].falseThrow = value
		}

		parseState++
	}
	if err = sc.Err(); err != nil {
		return nil, fmt.Errorf("scan: %v", err)
	}

	for i, m := range monkeys {
		if m == nil {
			return nil, fmt.Errorf("monkey %d is not initialized", i)
		}
	}

	return monkeys, nil
}
