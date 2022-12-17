package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
)

func main() {
	answer, err := solveDistressSignal()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to solve distress signal: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(answer)
}

func solveDistressSignal() (int, error) {
	lists, err := parseLists()
	if err != nil {
		return 0, fmt.Errorf("parse lists: %v", err)
	}
	dividerPacket2, dividerPacket6 := []any{[]any{2.0}}, []any{[]any{6.0}}
	lists = append(lists, dividerPacket2, dividerPacket6)

	sort.Slice(lists, func(i, j int) bool {
		return isListPairInTheRightOrder(lists[i], lists[j]) == comparisonResultRightOrder
	})

	dividerPacket2Idx, ok := sort.Find(len(lists), func(i int) int {
		return int(isListPairInTheRightOrder(dividerPacket2, lists[i]))
	})
	if !ok {
		return 0, errors.New("divider packet 2 not found")
	}
	dividerPacket2Idx++
	dividerPacket6Idx, ok := sort.Find(len(lists), func(i int) int {
		return int(isListPairInTheRightOrder(dividerPacket6, lists[i]))
	})
	if !ok {
		return 0, errors.New("divider packet 6 not found")
	}
	dividerPacket6Idx++

	decoderKey := dividerPacket2Idx * dividerPacket6Idx
	return decoderKey, nil
}

type comparisonResult int

const (
	comparisonResultRightOrder = iota - 1
	comparisonResultContinue
	comparisonResultWrongOrder
)

func isPairInTheRightOrder(l, r any) comparisonResult {
	switch typedL := l.(type) {
	case float64:
		switch typedR := r.(type) {
		case float64:
			return isFloatPairInTheRightOrder(typedL, typedR)
		case []any:
			return isListPairInTheRightOrder([]any{l}, typedR)
		default:
			panic(fmt.Sprintf("unknown type %T", r))
		}
	case []any:
		switch typedR := r.(type) {
		case float64:
			return isListPairInTheRightOrder(typedL, []any{r})
		case []any:
			return isListPairInTheRightOrder(typedL, typedR)
		default:
			panic(fmt.Sprintf("unknown type %T", l))
		}
	default:
		panic(fmt.Sprintf("unknown type %T", l))
	}
}

func isListPairInTheRightOrder(l, r []any) comparisonResult {
	for i, lv := range l {
		if i == len(r) {
			return comparisonResultWrongOrder
		}

		switch res := isPairInTheRightOrder(lv, r[i]); res {
		case comparisonResultRightOrder, comparisonResultWrongOrder:
			return res
		}
	}
	if len(l) < len(r) {
		return comparisonResultRightOrder
	}
	return comparisonResultContinue
}

func isFloatPairInTheRightOrder(l, r float64) comparisonResult {
	switch {
	case l < r:
		return comparisonResultRightOrder
	case l == r:
		return comparisonResultContinue
	default:
		return comparisonResultWrongOrder
	}
}

func parseLists() ([][]any, error) {
	f, err := os.Open("input.txt")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)

	var lists [][]any
	for line := 1; sc.Scan(); line++ {
		b := sc.Bytes()
		if line%3 == 0 {
			if len(b) != 0 {
				return nil, fmt.Errorf("expected empty line after pair on line %d", line)
			}
			continue
		}

		var list []any
		if err = json.Unmarshal(b, &list); err != nil {
			return nil, fmt.Errorf("unmarshal line %d", line)
		}
		lists = append(lists, list)
	}
	if err = sc.Err(); err != nil {
		return nil, fmt.Errorf("scan: %v", err)
	}

	for i, l := range lists {
		if err = validateList(l); err != nil {
			return nil, fmt.Errorf("validate list %d: %v", i, err)
		}
	}

	return lists, nil
}

func validateList(l []any) error {
	for i, el := range l {
		switch el := el.(type) {
		case float64:
		case []any:
			if err := validateList(el); err != nil {
				return fmt.Errorf("list at index %d", i)
			}
		default:
			return fmt.Errorf("list element at index %d has wrong type %T", i, el)
		}
	}
	return nil
}
